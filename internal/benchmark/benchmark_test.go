package benchmark

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

func BenchmarkApiConcurrent(b *testing.B) {
	requestCounts := []int{1}

	for _, count := range requestCounts {
		b.Run(fmt.Sprintf("Requests-%d", count), func(b *testing.B) {
			client := &http.Client{}

			var wg sync.WaitGroup

			for i := 0; i < count; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					response, err := client.Get("http://localhost:8000")
					if err != nil {
						b.Error(err)
						return
					}

					// Read response in a streaming manner
					buffer := make([]byte, 4096)
					for {
						_, err := response.Body.Read(buffer)

						if err == io.EOF {
							break
						}
						if err != nil {
							b.Error(err)
							return
						}
					}

					err = response.Body.Close()
					if err != nil {
						b.Error(err)
						return
					}
				}()
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			}

			wg.Wait()
		})
	}
}

func BenchmarkApi(b *testing.B) {
	client := &http.Client{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response, err := client.Get("http://localhost:8000")
		if err != nil {
			b.Fatal(err)
		}

		buffer := make([]byte, 4096)
		for {
			_, err := response.Body.Read(buffer)

			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}

		_ = response.Body.Close()
	}
}
