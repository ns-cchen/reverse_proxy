package proxy

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type StreamingReverseProxy struct {
	*httputil.ReverseProxy
}

func NewStreamingReverseProxy(target *url.URL) *StreamingReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = ModifyResponse
	return &StreamingReverseProxy{
		ReverseProxy: proxy,
	}
}

func ModifyResponse(response *http.Response) error {
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return err
		}

		response.Body = io.NopCloser(reader)
		response.Header.Del("Content-Encoding")
		response.Header.Del("Content-Length")
		response.Header.Set("Transfer-Encoding", "chunked")
		response.ContentLength = -1
		response.Uncompressed = true
	}
	return nil
}
