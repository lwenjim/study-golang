package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

type requestContextKey struct {
}

type requestContextValue struct {
	requestId string
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}

func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	fmt.Fprintf(w, "Hello world!")
}

func healthHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	if r.Method != "GET" {
		config.logger.Printf("error=\"Invalid request\" path=%s method=%s", r.URL.Path, r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "ok")
}

func panicHanlder(w http.ResponseWriter, r *http.Request, config appConfig) {
	panic("I panicked")
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
	mux.Handle("/panic", &app{config: config, handler: panicHanlder})
}

func loggingMiddleware(h http.Handler) http.Handler {
	var requestId string
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(writer, request)
			ctx := request.Context()
			v := ctx.Value(requestContextKey{})
			if m, ok := v.(requestContextValue); ok {
				requestId = m.requestId
			}
			log.Printf("request_id=%s protocal=%s path=%s method=%s duration=%f", requestId, request.Proto, request.URL.Path, request.Method, time.Now().Sub(startTime).Seconds())
		},
	)
}

func panicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				if rValue := recover(); rValue != nil {
					log.Println("panic detected when handling request", rValue)
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(writer, "Unexpected server error occured")
				}
			}()
			h.ServeHTTP(writer, request)
		},
	)
}

func addRequestIdMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			requestID := uuid.NewString()
			c := requestContextValue{
				requestId: requestID,
			}
			currentCtx := request.Context()

			newCtx := context.WithValue(currentCtx, requestContextKey{}, c)
			rWithContext := request.WithContext(newCtx)

			h.ServeHTTP(writer, rWithContext)
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
	setupHandlers(mux, config)

	m := addRequestIdMiddleware(
		loggingMiddleware(
			panicMiddleware(
				mux,
			),
		),
	)
	log.Fatal(http.ListenAndServe(listenAddr, m))
}
