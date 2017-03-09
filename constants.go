package godiawi

import "time"

type DiawiStatus int

// Diawi Status Codes
const (
	// Field name constants
	fileFieldName  = "file"
	tokenFieldName = "token"
	jobFieldName   = "job"

	// Optional field name constants

	findByUDIDFieldName = "find_by_udid"

	wallOfAppsFieldName       = "wall_of_apps"
	passwordFieldName         = "password"
	commentFieldName          = "comment"
	callbackURLFieldName      = "callback_url"
	callbackEmailsFieldName   = "callback_emails"
	installationNotifications = "installation_notifications"

	// Status Response Codes
	Processing   DiawiStatus = 2001
	Ok           DiawiStatus = 2000
	ErrorOccured DiawiStatus = 4000
)

var (
	// Request urls
	uploadURL = "https://upload.diawi.com/"
	statusURL = "https://upload.diawi.com/status"

	// Timeouts
	UploadTimeoutSeconds time.Duration = 60
	StatusTimeoutSeconds time.Duration = 10
	StatusPollingMax                   = 10
)
