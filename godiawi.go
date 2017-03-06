// Package godiawi wraps the diawi api in go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	// Field name constants
	FileFieldName  = "file"
	TokenFieldName = "token"
	JobFieldName   = "job"

	// Request url constants
	requestURL = "https://upload.diawi.com/"
	statusURL  = "https://upload.diawi.com/status"

	// Status Response Codes
	StatusProcessing   = 2001
	StatusOk           = 2000
	StatusErrorOccured = 4000
)

// UploadRequest is used to upload apps to diawi.
type UploadRequest struct {
	// Required parameters
	Token string
	File  string

	// Optional parameters
	WallOfApps              bool
	FindByUDID              bool
	InstallationNotifcation bool
	Password                string
	Comment                 string
	CallbackUrl             string
	CallbackEmails          []string
}

type FormWriter struct {
	buff bytes.Buffer
	mw   *multipart.Writer
}

func (fw *FormWriter) AddFormFile(fieldName, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fieldWriter, err := fw.mw.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}

	if _, err = io.Copy(fieldWriter, f); err != nil {
		return err
	}

	return nil
}

func (fw *FormWriter) AddField(fieldName, fieldValue string) error {
	fieldWriter, err := fw.mw.CreateFormField(fieldName)
	if err != nil {
		return err
	}

	if _, err = fieldWriter.Write([]byte(fieldValue)); err != nil {
		return err
	}

	return nil
}

func (fw *FormWriter) GetBuffer() bytes.Buffer {
	return fw.buff
}

func (fw *FormWriter) Close() {
	fw.mw.Close()
}

func (upRequest *UploadRequest) Upload() (*UploadResponse, error) {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)

	formWriter := FormWriter{buff: b, mw: mw}

	formWriter.AddFormFile(FileFieldName, upRequest.File)

	formWriter.AddField(TokenFieldName, upRequest.Token)

	req, err := http.NewRequest("POST", requestURL, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

	// Submit the request
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %s", res.Status)
	}

	resData, err := ioutil.ReadAll(res.Body)

	uploadRes := UploadResponse{}
	err = json.Unmarshal(resData, &uploadRes)
	if err != nil {
		return nil, err
	}

	return &uploadRes, nil
}

// UploadResponse contains the response provided by diawi
// following an upload request. Contains the job identifier
// for the upload.
type UploadResponse struct {
	JobIdentifier string `json:"job"`
}

// StatusRequest is used to poll diawi to see the status
// of the upload.
type StatusRequest struct {
	// Required parameters
	Token         string
	JobIdentifier string
}

func (s *StatusRequest) getJobStatus() (*StatusResponse, error) {

	url := statusURL + "?token=" + s.Token + "&job=" + s.JobIdentifier
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Submit the request
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %s", res.Status)
	}

	resData, err := ioutil.ReadAll(res.Body)

	statusRes := StatusResponse{}
	err = json.Unmarshal(resData, &statusRes)
	if err != nil {
		return nil, err
	}

	return &statusRes, nil
}

func (s *StatusRequest) waitForFinishedStatus() (*StatusResponse, error) {
	sr, err := s.getJobStatus()
	if err != nil {
		return nil, err
	}

	for {
		// log.Println("waitForFinishedStatus: ")
		// log.Println(sr)

		switch sr.Status {
		case StatusProcessing:
			time.Sleep(1 * time.Second) // diawi documentation recommends 1 second between each status request.
			sr, err = s.getJobStatus()
			if err != nil {
				return nil, err
			}
		case StatusOk:
			return sr, nil
		case StatusErrorOccured:
			return sr, fmt.Errorf("Response included error status = 4000. sr: %s", sr.String())
		default:
			return sr, fmt.Errorf("Unknown status error")
		}
	}
}

// StatusResponse is the response provided by the
type StatusResponse struct {
	Status  int
	Message string

	// Only provided when upload successful
	Hash string
	Link string
}

func (s *StatusResponse) String() string {
	//return "status: " + String(s.Status) + " message: " + s.Message + " hash: " + s.Hash + " link: " + s.Link
	return fmt.Sprintf("status: %d messsage: %s hash: %s link: %s", s.Status, s.Message, s.Hash, s.Link)
}

func main() {
	token := "YourTokenHere"
	file := "DiawiExampleApp.ipa"

	upRequest := UploadRequest{}
	upRequest.Token = token
	upRequest.File = file

	upResponse, err := upRequest.Upload()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("Status Response")
	// log.Println(statusRes)

	statReq := StatusRequest{}
	statReq.JobIdentifier = upResponse.JobIdentifier
	statReq.Token = token

	sr, err := statReq.waitForFinishedStatus()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(sr)
}
