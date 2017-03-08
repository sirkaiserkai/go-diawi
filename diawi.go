package godiawi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	""
)

type MultipartForm bytes.Buffer

type Diawi interface {
	GetStatus(fw FormWriter, responseStruct *interface{}) error
	UploadApp(fw FormWriter) (*http.Response, error)
}

type DiawiService struct {
	uploadUrl string
	statusUrl string
}

func NewDiawiService() DiawiService {
	return DiawiService{uploadUrl: uploadURL, statusUrl: statusURL}
}

func (d *DiawiService) GetStatus(fw FormWriter, responseStruct *interface{}) error {
	req, err := http.NewRequest("POST", uploadURL, fw.GetBuffer())
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", fw.mw.FormDataContentType())

	// submit the request
	client := &http.Client{Timeout: UploadTimeoutSeconds * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	resData, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(resData, responseStruct)
	if err != nil {
		return nil, err
	}

	return nil
}
