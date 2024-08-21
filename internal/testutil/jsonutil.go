package testutil

import (
	"fmt"
	"io"
	"time"
)

func VerifyStreamingDecompression(r io.Reader) (int64, error) {
	buffer := make([]byte, 4096)
	totalBytes := int64(0)
	startTime := time.Now()
	lastReportTime := startTime
	lastReportBytes := int64(0)

	for {
		n, err := r.Read(buffer)
		totalBytes += int64(n)

		// Report progress every second
		if time.Since(lastReportTime) >= time.Second {
			bytesPerSecond := (totalBytes - lastReportBytes) / int64(time.Since(lastReportTime).Seconds())
			fmt.Printf("Received %d bytes. Current speed: %d bytes/second\n", totalBytes, bytesPerSecond)
			lastReportTime = time.Now()
			lastReportBytes = totalBytes
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return totalBytes, fmt.Errorf("error reading stream: %v", err)
		}
	}

	duration := time.Since(startTime)
	avgSpeed := float64(totalBytes) / duration.Seconds()

	fmt.Printf("Total bytes received: %d\n", totalBytes)
	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Average speed: %.2f bytes/second\n", avgSpeed)

	return totalBytes, nil
}
