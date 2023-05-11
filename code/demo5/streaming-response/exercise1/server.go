package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	fileName := v.Get("filename")

	if len(fileName) == 0 {
		http.Error(w, "fileName query parameter not specified", http.StatusBadRequest)
		return
	}

	fileName = filepath.Base(fileName)
	if strings.HasPrefix(fileName, ".") {
		http.Error(w, "Invalid fileName specified", http.StatusBadRequest)
		return
	}

	f, err := os.Open(fileName)
	if err != nil {
		http.Error(w, "Enter opening file", http.StatusInternalServerError)
		return
	}

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)

	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	f.Seek(0, 0)

	contextType := http.DetectContentType(buffer)

	log.Println(contextType)

	w.Header().Set("Content-Type", contextType)

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8880"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/download", downloadFileHandler)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s, Error: %v", listenAddr, err)
	}
}
