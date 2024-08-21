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

	targetURL, err := url.Parse(testServer.URL)
	assert.NoError(t, err)
	proxy := NewStreamingReverseProxy(targetURL)
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	request, err := http.NewRequest("GET", proxyServer.URL, nil)
	assert.NoError(t, err)
	client := &http.Client{}
	response, err := client.Do(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, []string{"chunked"}, response.TransferEncoding)
	assert.Empty(t, response.Header.Get("Content-Encoding"))
	assert.Empty(t, response.Header.Get("Content-Length"))

	totalBytes, _ := testutil.ReadBody(response.Body)
	assert.NotZero(t, totalBytes)
}

func TestProxyModifyResponse(t *testing.T) {
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

			response := &http.Response{
				Header:        make(http.Header),
				Body:          body,
				ContentLength: int64(contentLength),
			}
			if tt.contentEncoding != "" {
				response.Header.Set("Content-Encoding", tt.contentEncoding)
			}
			response.Header.Set("Content-Length", strconv.Itoa(contentLength))

			err := ModifyResponse(response)
			require.NoError(t, err)

			if tt.expectedRemoved {
				assert.Empty(t, response.Header.Get("Content-Encoding"))
				assert.Empty(t, response.Header.Get("Content-Length"))
				assert.Equal(t, "chunked", response.Header.Get("Transfer-Encoding"))
				assert.Equal(t, int64(-1), response.ContentLength)
				bodyContent, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(bodyContent))
			} else {
				assert.Equal(t, tt.contentEncoding, response.Header.Get("Content-Encoding"))
				assert.Equal(t, strconv.Itoa(contentLength), response.Header.Get("Content-Length"))
				assert.Equal(t, int64(contentLength), response.ContentLength)
				bodyContent, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.body, string(bodyContent))
			}
		})
	}
}
