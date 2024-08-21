package testutil

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
)

type TestObject struct {
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

func StartTestServer(size int, timeUpperBoundInMillisecond int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)

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
			obj := generateRandomObject(i)
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

func generateRandomObject(id int) TestObject {
	faker := gofakeit.NewFaker(source.NewCrypto(), true)

	return TestObject{
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

func GenerateTestData(size int) []byte {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	encoder := json.NewEncoder(gzw)

	_, _ = gzw.Write([]byte("["))
	for i := 0; i < size; i++ {
		if i > 0 {
			_, _ = gzw.Write([]byte(","))
		}
		obj := generateRandomObject(i)
		_ = encoder.Encode(obj)
	}
	_, _ = gzw.Write([]byte("]"))
	_ = gzw.Close()
	return buf.Bytes()
}
