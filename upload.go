package godiawi

import (
	"errors"
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

	ds Diawi
}

// NewUploadRequest is the recommended means to create
// UploadRequest structs
func NewUploadRequest(token, file string) UploadRequest {
	ur := UploadRequest{Token: token, File: file}
	ur.ds = NewDiawiService()

	return ur
}

var EmptyFileFieldError = errors.New("File value left blank")
var EmptyTokenFieldError = errors.New("Token value left blank")

func (upRequest *UploadRequest) Upload() (*UploadResponse, error) {

	formWriter := NewFormWriter()

	if upRequest.File != "" {
		formWriter.AddFormFile(FileFieldName, upRequest.File)
	} else {
		return nil, EmptyFileFieldError
	}

	if upRequest.Token != "" {
		formWriter.AddField(TokenFieldName, upRequest.Token)
	} else {
		return nil, EmptyTokenFieldError
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
