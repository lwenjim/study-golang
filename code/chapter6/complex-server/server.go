package main

import (
	"github.com/lwenjim/code/chapter6/complex-server/config"
	"github.com/lwenjim/code/chapter6/complex-server/handlers"
	"github.com/lwenjim/code/chapter6/complex-server/middleware"
	"io"
	"log"
	"net/http"
	"os"
)

func setupServer(mux *http.ServeMux, w io.Writer) http.Handler {
	conf := config.InitConfig(w)
	handlers.Register(mux, conf)
	return middleware.RegisterMiddleware(mux, conf)
}
func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8880"
	}
	mux := http.NewServeMux()
	wrappedMux := setupServer(mux, os.Stdout)
	log.Fatal(http.ListenAndServe(listenAddr, wrappedMux))
}
