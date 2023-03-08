package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println("ping: Got a request")
	fmt.Fprintf(w, "pong")
}

func doSomeWork() {
	time.Sleep(20 * time.Second)
}

func createHTTPGetRequestWithTrace(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			fmt.Printf("DNS start info:%+v\n", info)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Done Info:%+v\n", info)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			fmt.Printf("Got conn:%+v\n", info)
		},
		TLSHandshakeStart: func() {
			fmt.Printf("TLS HandShake Start\n")
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			fmt.Printf("TLS HandShake Done\n")
		},
		PutIdleConn: func(err error) {
			fmt.Printf("Put Idle Conn Error:%+v\n", err)
		},
	}
	ctxTrace := httptrace.WithClientTrace(req.Context(), trace)
	req = req.WithContext(ctxTrace)
	return req, err
}

func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("I starting processing the request")
	doSomeWork()

	req, err := createHTTPGetRequestWithTrace(r.Context(), "http://localhost:8880/ping")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	log.Println("Outgoing HTTP request")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request:%v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(data))
	log.Println("I finished processing the request")
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8880"
	}
	mux := http.NewServeMux()

	timeoutDuration := 1 * time.Second
	userHandler := http.HandlerFunc(handleUserAPI)
	hTimeout := http.TimeoutHandler(userHandler, timeoutDuration, "I can out of time")

	mux.Handle("/api/users", hTimeout)
	mux.HandleFunc("/ping", handlePing)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
