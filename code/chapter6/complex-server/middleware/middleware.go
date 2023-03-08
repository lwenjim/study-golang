package middleware

import (
	"fmt"
	"github.com/lwenjim/code/chapter6/complex-server/config"
	"net/http"
	"time"
)

func loggingMiddleware(h http.Handler, c config.AppConfig) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			t1 := time.Now()
			h.ServeHTTP(writer, request)
			requestDuration := time.Now().Sub(t1).Seconds()
			c.Logger.Printf("protocol=%s path=%s method=%s duration=%f", request.Proto, request.URL.Path, request.Method, requestDuration)
		},
	)
}

func panicMiddleware(h http.Handler, c config.AppConfig) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				if rValue := recover(); rValue != nil {
					c.Logger.Println("panic detected", rValue)
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(writer, "Unexpected server error occurred")
				}
			}()
			h.ServeHTTP(writer, request)
		},
	)
}
