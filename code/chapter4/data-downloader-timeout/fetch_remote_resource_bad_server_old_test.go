package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func startBadTestHTTPServerV1() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				time.Sleep(60 * time.Second)
				fmt.Fprint(writer, "Hello World")
			},
		),
	)
	return ts
}

func startBadTestHTTPServerV2(shutdownServer chan struct{}) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		<-shutdownServer
		fmt.Fprint(writer, "Hello World")
	}))
	return ts
}

func TestFetchBadRemoteResourceV1(t *testing.T) {
	ts := startBadTestHTTPServerV1()
	defer ts.Close()

	client := createHTTPClientWithTimeout(200 * time.Microsecond)
	_, err := fetchRemoteResource(client, ts.URL)
	if err == nil {
		t.Fatal("Expected non-nil error")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("Expected error to contain: context deadline exceeded, Got: %v", err.Error())
	}
}

func TestFetchBadRemoteResourceV2(t *testing.T) {
	shutdownServer := make(chan struct{})
	ts := startBadTestHTTPServerV2(shutdownServer)
	defer ts.Close()
	defer func() {
		shutdownServer <- struct{}{}
	}()

	client := createHTTPClientWithTimeout(200 * time.Microsecond)
	_, err := fetchRemoteResource(client, ts.URL)
	if err == nil {
		t.Log("Expected non-nil error")
		t.Fail()
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("Expected error to contain: context deadline exceeded, Got: %v", err.Error())
	}
}
