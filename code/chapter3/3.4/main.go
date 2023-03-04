package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

type pkgData struct {
	Name     string
	Version  string
	FileName string
	Bytes    io.Reader
}

type pkgRegisterResult struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
}

type Form struct {
	Value map[string][]string
	File  map[string][]*FileHeader
}

func createMultipartMessage(data pkgData) ([]byte, string, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer

	mw := multipart.NewWriter(&b)

	fw, err = mw.CreateFormField("name")

	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Name)

	fw, err = mw.CreateFormField("version")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Version)

	fw, err = mw.CreateFormFile("fieldData", data.FileName)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(fw, data.Bytes)
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}
	contentType := mw.FormDataContentType()
	return b.Bytes(), contentType, nil
}

func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {
	p := pkgRegisterResult{}
	payload, contentType, err := createMultipartMessage(data)
	if err != nil {
		return p, err
	}
	reader := bytes.NewReader(payload)
	r, err := http.Post(url, contentType, reader)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()

	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(respData, &p)
	return p, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}

func packageRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		d := pkgRegisterResult{}
		err := r.ParseMultipartForm(5000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mForm := r.MultipartForm
		f := mForm.File["filedata"][0]
		d.Id = fmt.Sprintf("%s-%s", mForm.Value["name"][0], mForm.Value["version"][0])
		d.Filename = f.Filename
		d.Size = f.Size
		jsonData, err := json.Marshal(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonData))

	} else {
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
		return
	}

}

func main() {

}
