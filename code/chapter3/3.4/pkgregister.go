package pkgregister

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/textproto"
)

type pkgData struct {
	Name     string
	Version  string
	Filename string
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
