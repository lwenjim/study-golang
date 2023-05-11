package main

import (
	"testing"
)

func Test_fetchPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	packages, err := fetchPackageData(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 2 {
		t.Fatalf("Expected 2 packages, Got back: %d", len(packages))
	}
}
