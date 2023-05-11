package main

import (
	"fmt"
	"io"
	"time"
)

func longRunningProcess(w *io.PipeWriter) {
	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "hello")
		time.Sleep(1 * time.Second)
	}
	w.Close()
}
func main() {

}
