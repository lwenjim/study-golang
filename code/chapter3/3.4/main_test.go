package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func startTestPackageServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(packageRegHandler))
	return ts
}

func Test_registerPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	p := pkgData{
		Name:     "mypackage",
		Version:  "0.1",
		FileName: "mypackage-0.1,tar.gz",
		Bytes:    strings.NewReader("data"),
	}
	pResult, err := registerPackageData(ts.URL, p)
	if err != nil {
		t.Fatal(err)
	}
	if pResult.Id != fmt.Sprintf("%s-%s", p.Name, p.Version) {
		if pResult.Id != fmt.Sprintf("%s-%s", p.Name, p.Version) {
			t.Errorf("Expected package ID to be %s-%s, Got: %s", p.Name, p.Version, pResult.Id)
		}
		if pResult.Filename != p.FileName {
			t.Errorf("Expected package filename to be %s, Got: %s", p.FileName, pResult.Filename)
		}
		if pResult.Size != 4 {
			t.Errorf("Expected package size to be 4, Got: %d", pResult.Size)
		}
	}
}
