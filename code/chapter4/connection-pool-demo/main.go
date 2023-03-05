package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	transport := &http.Transport{
		IdleConnTimeout: 10 * time.Second,
	}
	client := http.Client{
		Timeout:   d,
		Transport: transport,
	}
	return &client
}

func createHTTPGetRequestWithTrace(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	trace := &httptrace.ClientTrace{
		DNSDone: func(info httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", info)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", info)
		},
		TLSHandshakeStart: func() {
			fmt.Printf("TLS HandShake Start\n")
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			fmt.Printf("TLS HandShake DOne\n")
		},
		PutIdleConn: func(err error) {
			fmt.Printf("Put Idle Conn Error: %+v\n", err)
		},
	}
	ctxTrace := httptrace.WithClientTrace(req.Context(), trace)
	req = req.WithContext(ctxTrace)
	return req, err
}

func main() {
	d := 5 * time.Second
	ctx := context.Background()
	client := createHTTPClientWithTimeout(d)

	req, err := createHTTPGetRequestWithTrace(ctx, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for {
		resp, _ := client.Do(req)
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		fmt.Printf("Resp protocol:%#v\n", resp.Proto)
		time.Sleep(1 * time.Second)
		fmt.Println("------------")
	}
}
