package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"reverse_proxy/internal/proxy"
)

func main() {
	targetURL := flag.String("target", "http://localhost:8080", "Target server URL")
	listenAddr := flag.String("listen", ":8000", "Listen address")
	flag.Parse()

	target, err := url.Parse(*targetURL)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	reverseProxy := proxy.NewStreamingReverseProxy(target)
	log.Fatal(http.ListenAndServe(*listenAddr, reverseProxy))
}
