package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"reverse_proxy/internal/proxy"
	"reverse_proxy/internal/testutil"
)

func main() {
	listenAddr := flag.String("listen", ":8000", "Listen address")
	flag.Parse()

	testServer := testutil.StartTestServer(999999, 100)
	targetURL, err := url.Parse(testServer.URL)
	if err != nil {
		log.Fatal(err)
	}
	reverseProxy := proxy.NewStreamingReverseProxy(targetURL)
	log.Fatal(http.ListenAndServe(*listenAddr, reverseProxy))
}
