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
	log.Println("ping: GOt a request")
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "pong")
}

func doSomeWork(data []byte) {
	time.Sleep(5 * time.Second)
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

func handleUserApi(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)
		logger.Println("I started processing the request")

		pingServer := r.URL.Query().Get("ping_server")
		if len(pingServer) == 0 {
			pingServer = "http://localhost:8880"
		}
		req, err := createHTTPGetRequestWithTrace(
			r.Context(),
			pingServer+"/ping",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := &http.Client{}
		logger.Println("Outgoing HTTP request")
		resp, err := client.Do(req)
		if err != nil {
			logger.Printf("Error making request:%v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		logger.Println("Processing the response i got")

		go func() {
			doSomeWork(data)
			done <- true
		}()

		select {
		case <-done:
			logger.Println("doSameWork done:Continuing request processing")
		case <-r.Context().Done():
			logger.Printf("Aborting request processing: %v\n", r.Context().Err())
			return
		}

		fmt.Fprintf(w, string(data))
		logger.Println("I finished processing the request")
	}
}

func setupHandlers(mux *http.ServeMux, timeoutDuration time.Duration, logger *log.Logger) {
	userHandler := handleUserApi(logger)
	hTimeout := http.TimeoutHandler(
		userHandler,
		timeoutDuration,
		"I ran out of time",
	)
	mux.Handle("/api/users/", hTimeout)
	mux.HandleFunc("/ping", handlePing)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8880"
	}
	timeoutDuration := 30 * time.Second
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	mux := http.NewServeMux()
	setupHandlers(mux, timeoutDuration, logger)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
