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

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api", apiHandler)
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(writer, request)
			log.Printf("protocol=%s path=%s method=%s duration=%f",
				request.Proto,
				request.URL.Path,
				request.Method,
				time.Now().Sub(startTime).Seconds(),
			)
		},
	)
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
	setupHandlers(mux)

	m := loggingMiddleware(mux)

	log.Fatal(
		http.ListenAndServeTLS(listenAddr, tlsCertFile, tlsKeyFile, m),
	)

}
