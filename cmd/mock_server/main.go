package main

import (
	"flag"
	"fmt"
	"log"

	"reverse_proxy/internal/testutil"
)

func main() {
	port := flag.Int("port", 8080, "Port to run the test server on")
	size := flag.Int("size", 999999, "Size parameter for StartTestServer")
	timeUpperBound := flag.Int("time", 100, "Time upper bound in milliseconds for StartTestServer")
	flag.Parse()

	mockServer := testutil.StartMockServer(*size, *timeUpperBound)

	addr := fmt.Sprintf(":%d", *port)
	mockServer.Addr = addr
	log.Fatal(mockServer.ListenAndServe())
}
