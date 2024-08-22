package main

import (
	"log"
	"net/http"
	"net/url"

	"reverse_proxy/internal/proxy"
)

func main() {
	targetUrl, _ := url.Parse("http://localhost:8080")
	reverseProxy := proxy.NewStreamingReverseProxy(targetUrl)
	log.Fatal(http.ListenAndServe(":8000", reverseProxy))
}
