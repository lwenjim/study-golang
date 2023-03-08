package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func loggingMiddleware(h http.Handler, l *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(writer, request)
			l.Printf("protocol=%s path=%s method=%s duration=%f", request.Proto, request.URL.Path, request.Method, time.Now().Sub(startTime).Seconds())
		},
	)
}

func setupHandlerAndMiddleware(mux *http.ServeMux, l *log.Logger) http.Handler {
	mux.HandleFunc("/api", apiHandler)
	return loggingMiddleware(mux, l)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8880"
	}
	tlsCertFile := os.Getenv("TLS_CERT_FILE_PATH")
	tlsKeyFile := os.Getenv("TLS_KEY_FILE_PATH")

	if len(tlsCertFile) == 0 || len(tlsKeyFile) == 0 {
		log.Fatal("TLS_CERT_FILE_PATH and TLS_KEY_FILE_PATH must be specified")
	}
	mux := http.NewServeMux()

	l := log.New(os.Stdout, "tls-server", log.Lshortfile|log.LstdFlags)
	m := setupHandlerAndMiddleware(mux, l)
	log.Fatal(http.ListenAndServeTLS(listenAddr, tlsCertFile, tlsKeyFile, m))
}
