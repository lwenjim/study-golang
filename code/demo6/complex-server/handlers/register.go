package handlers

import (
	"github.com/lwenjim/study-golang/code/demo6/complex-server/config"
	"net/http"
)

func Register(mux *http.ServeMux, conf config.AppConfig) {
	mux.Handle("/healthz", &app{conf: conf, handler: healthCheckHandler})
	mux.Handle("/api", &app{conf: conf, handler: apiHandler})
	mux.Handle("/panic", &app{conf: conf, handler: panicHandler})
}
