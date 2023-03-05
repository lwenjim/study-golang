package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startHTTPServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			writer.Header().Set(k, v[0])
		}
		fmt.Fprint(writer, "I am the Request Header echoing program")
	}))
	return ts
}

func TestAddHeaderMiddleware(t *testing.T) {
	testHeaders := map[string]string{
		"X-Client-Id": "test-client",
		"X-Auth-Hash": "random$string",
	}
	client := createClient(testHeaders)
	ts := startHTTPServer()
	defer ts.Close()

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Expected non-nil [AU:JA] error, got: %v", err)
	}
	for k, v := range testHeaders {
		if resp.Header.Get(k) != testHeaders[k] {
			t.Fatalf("Expected header: %s:%s, Got:%s:%s", k, v, k, testHeaders[k])
		}
	}
}
