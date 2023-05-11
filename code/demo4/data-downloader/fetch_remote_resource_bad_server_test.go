package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startBadTestHTTPServer() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				//time.Sleep(60 * time.Second)
				fmt.Fprint(writer, "Hello World")
			},
		),
	)
	return ts
}

func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func TestFetchBadRemoteResource(t *testing.T) {
	ts := startBadTestHTTPServer()
	defer ts.Close()
	client := http.Client{Timeout: 4 * 1000 * time.Microsecond}
	data, err := fetchRemoteResource(&client, ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	expected := "Hello World"
	got := string(data)

	if expected != got {
		t.Errorf("Expected response to be: %s, Got: %s", expected, got)
	}
}

func main() {

}
