package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func startTestPackageServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(packageRegHander))
	return ts
}

func Test_packageRegHander(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	p := pkgData{
		Name:    "mypackage",
		Version: "0.1",
	}
	resp, err := registerPackageData(ts.URL, p)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "mypackage-0.1" {
		t.Errorf("Expected package id to be mypackage-0.1, Got: %s", resp.ID)
	}
}

func TestRegisterEmptyPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	p := pkgData{
		Name: "123",
	}
	resp, err := registerPackageData(ts.URL, p)
	if err != nil {
		t.Fatalf("Expected error to be non-nil, got nil")
	}
	if len(resp.ID) != 0 {
		t.Errorf("Expected package ID to be empty, got: %s", resp.ID)
	}
}
