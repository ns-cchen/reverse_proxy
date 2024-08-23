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

		response.Body = &gzipReaderCloser{
			gzipReader: reader,
			body:       response.Body,
		}
		response.Header.Del("Content-Encoding")
		response.Header.Del("Content-Length")
		response.Header.Set("Transfer-Encoding", "chunked")
		response.ContentLength = -1
		response.Uncompressed = true
	}
	return nil
}

type gzipReaderCloser struct {
	gzipReader *gzip.Reader
	body       io.Closer
}

func (grc *gzipReaderCloser) Read(p []byte) (n int, err error) {
	return grc.gzipReader.Read(p)
}

func (grc *gzipReaderCloser) Close() error {
	if err := grc.gzipReader.Close(); err != nil {
		return err
	}

	if err := grc.body.Close(); err != nil {
		return err
	}

	return nil
}
