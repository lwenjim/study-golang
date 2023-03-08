package middleware

import (
	"github.com/lwenjim/code/chapter6/complex-server/config"
	"net/http"
)

func RegisterMiddleware(mux *http.ServeMux, c config.AppConfig) http.Handler {
	return loggingMiddleware(panicMiddleware(mux, c), c)
}
