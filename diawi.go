package godiawi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type MultipartForm bytes.Buffer

type Diawi interface {
	UploadApp(fw FormWriter, responseStruct interface{}) error
	GetStatus(token, job string, responseStruct interface{}) error
}

func NewDiawiService() Diawi {
	return DiawiService{uploadURL: uploadURL, statusURL: statusURL}
}

type DiawiService struct {
	uploadURL string
	statusURL string
}

func (d DiawiService) UploadApp(fw FormWriter, responseStruct interface{}) error {
	return makeRequest(d.uploadURL, "POST", fw.mw.FormDataContentType(), fw.GetBuffer(), UploadTimeoutSeconds, responseStruct)
}

func (d DiawiService) GetStatus(token, job string, responseStruct interface{}) error {
	url := d.statusURL + "?token=" + token + "&job=" + job
	return makeRequest(url, "GET", "application/json", nil, StatusTimeoutSeconds, responseStruct)
}

func makeRequest(url, requestMethod, contentType string, body *bytes.Buffer, timeout time.Duration, responseStruct interface{}) error {
	var req *http.Request
	var err error

	// Causes error when nil interface passed in
	if body == nil {
		req, err = http.NewRequest(requestMethod, url, nil)
	} else {
		req, err = http.NewRequest(requestMethod, url, body)
	}

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	// submit the request
	client := &http.Client{Timeout: timeout * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("bad status: %d\n for request: %s\n Cannot parse body", res.StatusCode, url)
		}
		return fmt.Errorf("bad status: %d\n for request: %s\n Response Body: %s", res.StatusCode, url, string(body))
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resData, responseStruct)
	if err != nil {
		return err
	}

	return nil
}
