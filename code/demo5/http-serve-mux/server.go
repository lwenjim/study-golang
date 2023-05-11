package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type logLine struct {
	URL           string `json:"URL"`
	Method        string `json:"method"`
	ContentLength int64  `json:"content_length"`
	Protocol      string `json:"protocol"`
}

func logRequest(req *http.Request) {
	l := logLine{
		URL:           req.URL.String(),
		Method:        req.Method,
		ContentLength: req.ContentLength,
		Protocol:      req.Proto,
	}
	data, err := json.Marshal(&l)
	if err != nil {
		panic(err)
	}
	log.Println(string(data))
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello world!")
	logRequest(req)
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
	logRequest(req)
}

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/api", apiHandler)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8880"
	}
	mux := http.NewServeMux()
	setupHandlers(mux)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
