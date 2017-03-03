// Package godiawi wraps the diawi api in go
package godiawi

import (
	"net/http"
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

// UploadResponse contains the response provided by diawi
// following an upload request. Contains the job identifier
// for the upload.
type UploadResponse struct {
	JobIdentifier string
}

// StatusRequest is used to poll diawi to see the status
// of the upload.
type StatusRequest struct {
	// Required parameters
	Token         string
	JobIdentifier string
}

// StatusResponse is the response provided by the
type StatusResponse struct {
	Status  int
	Message string

	// Only provided when upload successful
	Hash string
	Link string
}
