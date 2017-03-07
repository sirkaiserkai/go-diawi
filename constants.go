package godiawi

type DiawiStatus int

const (
	// Field name constants
	FileFieldName  = "file"
	TokenFieldName = "token"
	JobFieldName   = "job"

	// Optional field name constants
	FindByUDIDFieldName       = "find_by_udid"
	WallOfAppsFieldName       = "wall_of_apps"
	PasswordFieldName         = "password"
	CommentFieldName          = "comment"
	CallbackURLFieldName      = "callback_url"
	CallbackEmailsFieldName   = "callback_emails"
	InstallationNotifications = "installation_notifications"

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
	UploadTimeoutSeconds = 60
	StatusTimeoutSeconds = 10
	StatusPollingMax     = 10
)
