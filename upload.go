package godiawi

import (
	"errors"
)

// UploadRequest is used to wrap the parameters for a diawi upload request
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

	ds diawi
}

// NewUploadRequest is the recommended means to create
// UploadRequest structs
func NewUploadRequest(token, file string) UploadRequest {
	ur := UploadRequest{Token: token, File: file}
	ur.ds = newDiawiService()

	return ur
}

var EmptyFileFieldError = errors.New("File value left blank")
var EmptyTokenFieldError = errors.New("Token value left blank")

// Upload requests the diawi service using the values set in the
// respective UploadRequest. Returns an UploadResponse provided diawi
// was able to process the request (Not guaranteeing it was successful).
// Returns an error object if an error is encountered.
func (upRequest *UploadRequest) Upload() (*UploadResponse, error) {

	formWriter := newformWriter()

	if upRequest.File != "" {
		formWriter.AddFormFile(fileFieldName, upRequest.File)
	} else {
		return nil, EmptyFileFieldError
	}

	if upRequest.Token != "" {
		formWriter.AddField(tokenFieldName, upRequest.Token)
	} else {
		return nil, EmptyTokenFieldError
	}

	if upRequest.Comment != "" {
		formWriter.AddField(commentFieldName, upRequest.Comment)
	}

	if upRequest.CallbackUrl != "" {
		formWriter.AddField(callbackURLFieldName, upRequest.CallbackUrl)
	}

	if len(upRequest.CallbackEmails) != 0 {
		formWriter.AddField(callbackEmailsFieldName, upRequest.CallbackEmails)
	}

	formWriter.AddField(findByUDIDFieldName, upRequest.FindByUDID)

	formWriter.AddField(wallOfAppsFieldName, upRequest.WallOfApps)

	formWriter.AddField(installationNotifications, upRequest.InstallationNotifcation)

	formWriter.Close()

	ur := UploadResponse{}
	err := upRequest.ds.UploadApp(formWriter, &ur)
	if err != nil {
		return nil, err
	}

	return &ur, nil
}

// UploadResponse contains the response provided by diawi
// following an upload request. Contains the job identifier
// for the upload.
type UploadResponse struct {
	JobIdentifier string `json:"job"`
}
