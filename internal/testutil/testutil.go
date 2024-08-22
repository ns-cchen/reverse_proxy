package testutil

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type FakeObject struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Address   string  `json:"address"`
	Company   string  `json:"company"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
}

func StartMockServer(size int, timeUpperBoundInMillisecond int) *http.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")

		gzw, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
		defer func(gzw *gzip.Writer) {
			_ = gzw.Close()
		}(gzw)

		encoder := json.NewEncoder(gzw)
		flusher, _ := w.(http.Flusher)

		_, _ = gzw.Write([]byte("["))
		for i := 0; i < size; i++ {
			if i > 0 {
				_, _ = gzw.Write([]byte(","))
			}
			obj := generateFakeObject(i)
			_ = encoder.Encode(obj)

			_ = gzw.Flush()
			flusher.Flush()
			time.Sleep(time.Duration(rand.Intn(timeUpperBoundInMillisecond)) * time.Millisecond)
		}
		_, _ = gzw.Write([]byte("]"))
		_ = gzw.Flush()
		flusher.Flush()
	})
	server := &http.Server{
		Handler: handler,
	}

	return server
}

func StartTestServer(size int, timeUpperBoundInMillisecond int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")

		gzw, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
		defer func(gzw *gzip.Writer) {
			_ = gzw.Close()
		}(gzw)

		encoder := json.NewEncoder(gzw)
		flusher, _ := w.(http.Flusher)

		_, _ = gzw.Write([]byte("["))
		for i := 0; i < size; i++ {
			if i > 0 {
				_, _ = gzw.Write([]byte(","))
			}
			obj := generateFakeObject(i)
			_ = encoder.Encode(obj)

			_ = gzw.Flush()
			flusher.Flush()
			time.Sleep(time.Duration(rand.Intn(timeUpperBoundInMillisecond)) * time.Millisecond)
		}
		_, _ = gzw.Write([]byte("]"))
		_ = gzw.Flush()
		flusher.Flush()
	}))
}

func generateFakeObject(id int) FakeObject {
	faker := gofakeit.New(0)

	return FakeObject{
		ID:        id,
		Name:      faker.Name(),
		Email:     faker.Email(),
		Phone:     faker.Phone(),
		Address:   faker.Address().Address,
		Company:   faker.Company(),
		Latitude:  faker.Latitude(),
		Longitude: faker.Longitude(),
		Country:   faker.Country(),
	}
}

func ReadBody(r io.Reader) (int64, error) {
	buffer := make([]byte, 4096)
	totalBytes := int64(0)
	startTime := time.Now()
	lastReportTime := startTime
	lastReportBytes := int64(0)

	for {
		n, err := r.Read(buffer)
		totalBytes += int64(n)

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
