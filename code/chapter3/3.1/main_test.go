package main

import (
	"testing"
)

func Test_fetchRemoteResource(t *testing.T) {
	ts := startTestHTTPServer()
	defer ts.Close()
	expected := "Hello World"
	data, err := fetchRemoteResource(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(data) {
		t.Errorf("Expected to be: %s, Got: %s", expected, data)
	}
}
