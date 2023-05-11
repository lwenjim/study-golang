package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}
func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	config.logger.Println("Handling API request")
	fmt.Fprintf(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	config.logger.Println("Handling healthcheck request")
	fmt.Fprintf(w, "ok")
}

func panicHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	panic("I panicked")
}

func setupHandles(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthCheckHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
	mux.Handle("/panic", &app{config: config, handler: panicHandler})
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(writer, request)
			log.Printf("path=%s method=%s duration=%f", request.URL.Path, request.Method, time.Now().Sub(startTime).Seconds())
		},
	)
}

func panicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				if rValue := recover(); rValue != nil {
					log.Println("panic detected", rValue)
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(writer, "Unexpected server error")
				}
			}()
			h.ServeHTTP(writer, request)
		},
	)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8880"
	}
	config := appConfig{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	mux := http.NewServeMux()
	setupHandles(mux, config)

	m := loggingMiddleware(panicMiddleware(mux))
	log.Fatal(http.ListenAndServe(listenAddr, m))
}
