package proxy

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reverse_proxy/internal/testutil"
)

func TestStreamingDecompression(t *testing.T) {
	testServer := testutil.StartTestServer(100, 100)
	defer testServer.Close()

	targetURL, _ := url.Parse(testServer.URL)
	proxy := NewStreamingReverseProxy(targetURL)
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	req, _ := http.NewRequest("GET", proxyServer.URL, nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	require.NotNil(t, resp)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Empty(t, resp.Header.Get("Content-Encoding"))

	totalBytes, _ := testutil.VerifyStreamingDecompression(resp.Body)
	assert.NotZero(t, totalBytes)
}

func TestProxyModifyResponse(t *testing.T) {
	proxy := &StreamingReverseProxy{}

	tests := []struct {
		name            string
		contentEncoding string
		body            string
		expectedBody    string
		expectedRemoved bool
	}{
		{"Gzip encoding", "gzip", "test body", "test body", true},
		{"No encoding", "", "test body", "test body", false},
		{"Different encoding", "deflate", "test body", "test body", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.ReadCloser
			var contentLength int

			if tt.contentEncoding == "gzip" {
				var buf bytes.Buffer
				gzipWriter := gzip.NewWriter(&buf)
				_, err := gzipWriter.Write([]byte(tt.body))
				require.NoError(t, err)
				err = gzipWriter.Close()
				require.NoError(t, err)
				body = io.NopCloser(&buf)
				contentLength = buf.Len()
			} else {
				body = io.NopCloser(strings.NewReader(tt.body))
				contentLength = len(tt.body)
			}

			res := &http.Response{
				Header:        make(http.Header),
				Body:          body,
				ContentLength: int64(contentLength),
			}
			if tt.contentEncoding != "" {
				res.Header.Set("Content-Encoding", tt.contentEncoding)
			}
			res.Header.Set("Content-Length", strconv.Itoa(contentLength))

			err := proxy.ModifyResponse(res)
			require.NoError(t, err)

			if tt.expectedRemoved {
				assert.Empty(t, res.Header.Get("Content-Encoding"))
				assert.Empty(t, res.Header.Get("Content-Length"))
				assert.Equal(t, int64(len(tt.expectedBody)), res.ContentLength)

				// Check the actual body content
				bodyContent, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(bodyContent))
			} else {
				assert.Equal(t, tt.contentEncoding, res.Header.Get("Content-Encoding"))
				assert.Equal(t, strconv.Itoa(contentLength), res.Header.Get("Content-Length"))
				assert.Equal(t, int64(contentLength), res.ContentLength)

				// Check that the body is unchanged
				bodyContent, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.body, string(bodyContent))
			}
		})
	}
}
