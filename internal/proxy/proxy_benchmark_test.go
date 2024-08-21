package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"reverse_proxy/internal/testutil"
)

func BenchmarkStreamingDecompression(b *testing.B) {
	pairs := [][2]int{{10, 100}, {100, 100}, {1000, 100}, {10000, 100}, {100000, 100}}

	for _, pair := range pairs {
		b.Run(fmt.Sprintf("Size-%d", pair[0]), func(b *testing.B) {
			testServer := testutil.StartTestServer(pair[0], pair[1])
			defer testServer.Close()

			targetURL, _ := url.Parse(testServer.URL)
			proxy := NewStreamingReverseProxy(targetURL)
			proxyServer := httptest.NewServer(proxy)
			defer proxyServer.Close()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				client := &http.Client{}
				req, err := http.NewRequest("GET", proxyServer.URL, nil)
				if err != nil {
					b.Fatal(err)
				}

				resp, err := client.Do(req)
				if err != nil {
					b.Fatal(err)
				}

				err = consumeBody(resp.Body)
				if err != nil {
					b.Fatal(err)
				}

				err = resp.Body.Close()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func consumeBody(body io.ReadCloser) error {
	buffer := make([]byte, 4096)
	for {
		_, err := body.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
