package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
func validateConfig(c httpConfig) error {
	allowVerbs := []string{"GET", "POST", "HEAD"}
	for _, v := range allowVerbs {
		if c.verb == v {
			return nil
		}
	}
	return ErrInvalidHTTPMethod
}

func HandleHttp(w io.Writer, args []string) error {
	var err error
	var c = httpConfig{}
	var outputFile string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&outputFile, "output", "", "File path to write the response into")

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

	err = validateConfig(c)
	if err != nil {
		if errors.Is(err, ErrInvalidHTTPMethod) {
			fmt.Fprintln(w, "Invalid HTTP method")
		}
		return err
	}

	c.url = fs.Arg(0)

	data, err := fetchRemoteResource(c.url)
	if err != nil {
		return err
	}
	if len(outputFile) != 0 {
		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Data saved to: %s\n", outputFile)
		return err
	}
	fmt.Fprintln(w, string(data))
	return nil
}

func startTestHttpServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/download", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "this is a response")
	})
	return httptest.NewServer(mux)
}
