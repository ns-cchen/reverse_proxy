package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"reverse_proxy/internal/testutil"
)

func BenchmarkStreamingDecompression(b *testing.B) {
	pairs := [][2]int{{100, 100}, {1000, 100}, {10000, 100}}

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
				req, _ := http.NewRequest("GET", proxyServer.URL, nil)
				resp, _ := client.Do(req)

				_, _ = testutil.VerifyStreamingDecompression(resp.Body)
				_ = resp.Body.Close()
			}
		})
	}
}
