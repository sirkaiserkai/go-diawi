package godiawi

import (
	"testing"
)

// Unit Tests
func TestErrorUpload(t *testing.T) {
	ur := UploadRequest{Token: "", File: ""}

	if _, err := ur.Upload(); err == nil {
		t.Error("Should receive error due to lack of file")
	}

	ur.File = "abcd"

	if _, err := ur.Upload(); err == nil {
		t.Error("Should receive error due to lack of token")
	}
}

// TODO: Complete this test
func TestUploadSuccess(t *testing.T) {
	ur := UploadRequest{Token: "", File: "app.ipa"}
}
