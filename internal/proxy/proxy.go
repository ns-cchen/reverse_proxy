package proxy

import (
	"bytes"
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
	return &StreamingReverseProxy{
		ReverseProxy: httputil.NewSingleHostReverseProxy(target),
	}
}

func (p *StreamingReverseProxy) ModifyResponse(res *http.Response) error {
	if res.Header.Get("Content-Encoding") == "gzip" {
		decompressedBody, err := fullyReadAndDecompress(res.Body)
		if err != nil {
			return err
		}
		_ = res.Body.Close()

		res.Body = io.NopCloser(bytes.NewReader(decompressedBody))
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Length")
		res.ContentLength = int64(len(decompressedBody))
		res.Uncompressed = true
	}
	return nil
}

func fullyReadAndDecompress(body io.ReadCloser) ([]byte, error) {
	reader, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}
	defer func(gzipReader *gzip.Reader) {
		_ = gzipReader.Close()
	}(reader)

	return io.ReadAll(reader)
}

type gzipReader struct {
	body io.ReadCloser
	zr   *gzip.Reader
}

func (g *gzipReader) Read(p []byte) (n int, err error) {
	if g.zr == nil {
		g.zr, err = gzip.NewReader(g.body)
		if err != nil {
			return 0, err
		}
	}
	return g.zr.Read(p)
}

func (g *gzipReader) Close() error {
	if g.zr != nil {
		_ = g.zr.Close()
	}
	return g.body.Close()
}
