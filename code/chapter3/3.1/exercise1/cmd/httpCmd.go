package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

type httpConfig struct {
	url  string
	verb string
}

func fetchRemoteResource(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
func HandleHttp(w io.Writer, args []string) error {
	var err error
	var v string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "GET", "HTTP method")

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.
		
http: <options> server`
		fmt.Fprintln(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	err = fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}
	c := httpConfig{verb: v}
	c.url = fs.Arg(0)

	data, err := fetchRemoteResource(c.url)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, string(data))
	return nil
}
