package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type pkgRegisterResult struct {
	ID string `json:"id"`
}
type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {
	p := pkgRegisterResult{}
	b, err := json.Marshal(data)
	if err != nil {
		return p, err
	}
	reader := bytes.NewReader(b)
	r, err := http.Post(url, "application/json", reader)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()
	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	if r.StatusCode != http.StatusOK {
		return p, errors.New(string(responseData))
	}
	fmt.Println(string(responseData))
	err = json.Unmarshal(responseData, &p)
	return p, err
}

func packageRegHander(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		p := pkgData{}
		d := pkgRegisterResult{}

		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(data, &p)
		if err != nil || len(p.Name) == 0 || len(p.Version) == 0 {
			return
		}
		d.ID = p.Name + "-" + p.Version

		jsonData, err := json.Marshal(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonData))
	} else {
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
		return
	}
}

func main() {

}
