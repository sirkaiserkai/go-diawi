package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"mime/multipart"
	"net/http"
	"time"
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

func (upRequest *UploadRequest) Upload() (*UploadResponse, error) {

	formWriter := NewFormWriter()

	if upRequest.File != "" {
		formWriter.AddFormFile(FileFieldName, upRequest.File)
	} else {
		return nil, fmt.Errorf("File value left blank")
	}

	if upRequest.Token != "" {
		formWriter.AddField(TokenFieldName, upRequest.Token)
	} else {
		return nil, fmt.Errorf("Token value left blank")
	}

	if upRequest.Comment != "" {
		formWriter.AddField(CommentFieldName, upRequest.Comment)
	}

	if upRequest.CallbackUrl != "" {
		formWriter.AddField(CallbackURLFieldName, upRequest.CallbackUrl)
	}

	if len(upRequest.CallbackEmails) != 0 {
		formWriter.AddField(CallbackEmailsFieldName, upRequest.CallbackEmails)
	}

	formWriter.AddField(FindByUDIDFieldName, upRequest.FindByUDID)

	formWriter.AddField(WallOfAppsFieldName, upRequest.WallOfApps)

	formWriter.AddField(InstallationNotifications, upRequest.InstallationNotifcation)

	formWriter.Close()

	b := formWriter.GetBuffer()
	req, err := http.NewRequest("POST", requestURL, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", formWriter.mw.FormDataContentType())

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
