package benchmark

import (
	"io"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

func BenchmarkApiConcurrent10Request(b *testing.B) {
	client := &http.Client{}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
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
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	wg.Wait()

}

func BenchmarkApiConcurrent100Request(b *testing.B) {
	client := &http.Client{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
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
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	wg.Wait()
}

func BenchmarkApiConcurrent1000Request(b *testing.B) {
	client := &http.Client{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
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
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	wg.Wait()
}

func BenchmarkApiConcurrent10000Request(b *testing.B) {
	client := &http.Client{}
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
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
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	wg.Wait()
}

func BenchmarkApi(b *testing.B) {
	client := &http.Client{}
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
