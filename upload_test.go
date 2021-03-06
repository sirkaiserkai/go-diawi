package godiawi

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
)

// Test flags
var (
	diawiToken  = flag.String("token", "", "token to be used in integraton tests")
	appFile     = flag.String("app", "", "app to be uploaded for integration tests")
	integration = flag.Bool("integration", false, "run integration tests")
)

// Used for testing purposes
type diawiTestSerivce struct {
	failUpload          bool // Flag set when we want the service to return a negative for upload
	returnFailStatus    bool
	returnSuccessStatus bool
}

func (d diawiTestSerivce) UploadApp(fw formWriter, responseStruct interface{}) error {
	if d.failUpload {
		return fmt.Errorf("Expected Error")
	}

	ur, ok := responseStruct.(*UploadResponse)
	if !ok {
		return fmt.Errorf("response struct is wrong type: %T", responseStruct)
	}

	ur.JobIdentifier = "123456"

	return nil
}

func (d diawiTestSerivce) GetStatus(token, job string, responseStruct interface{}) error {

	sr, ok := responseStruct.(*StatusResponse)
	if !ok {
		return fmt.Errorf("response struct is wrong type: %T", responseStruct)
	}

	if d.returnFailStatus {
		sr.Status = ErrorOccured
		sr.Message = "Error occured"
		return nil
	} else if d.returnSuccessStatus {
		sr.Status = Ok
		sr.Message = "Hello, world!"
		sr.Hash = "12345"
		sr.Link = "http://fake.come"
		return nil
	} else {
		sr.Status = Processing
		sr.Message = "Processing request"
	}
	return nil
}

var TestToken = "1234567abcdef"
var TestAppFile = "app.ipa"

// TESTS //

func TestMain(m *testing.M) {
	flag.Parse()
	if *integration {
		if *diawiToken == "" || *appFile == "" {
			log.Fatal("In order to perform integration test token and app (file) flags need to be set")
		}
	}

	result := m.Run()
	os.Exit(result)
}

// Unit Tests
func TestUploadMissingFile(t *testing.T) {
	ur := UploadRequest{Token: TestToken, File: ""}

	if _, err := ur.Upload(); err == nil {
		t.Error("Should receive error due to lack of file")
	}
}

func TestUploadMissingToken(t *testing.T) {
	ur := UploadRequest{Token: "", File: TestAppFile}

	if _, err := ur.Upload(); err == nil {
		t.Error("Should receive error due to lack of token")
	}
}

func TestUploadSuccess(t *testing.T) {
	ur := NewUploadRequest(TestToken, TestAppFile)
	ur.ds = diawiTestSerivce{}

	uploadResponse, err := ur.Upload()
	if err != nil {
		t.Fatal("Error occured: %s", err)
	}

	if uploadResponse.JobIdentifier == "" {
		t.Error("Job identifer blank")
	}
}

func TestVerboseUploadSuccess(t *testing.T) {
	ur := NewUploadRequest(TestToken, TestAppFile)
	ur.ds = diawiTestSerivce{}

	ur.WallOfApps = true
	ur.FindByUDID = true
	ur.InstallationNotifcation = true
	ur.Password = "Password"
	ur.Comment = "Hello, world!"
	ur.CallbackUrl = "http://example.com"
	ur.CallbackEmails = []string{"example@example.com", "fake@example.com"}

	uploadResponse, err := ur.Upload()
	if err != nil {
		t.Fatal("Error occured: %s", err)
	}

	if uploadResponse.JobIdentifier == "" {
		t.Error("Job identifer blank")
	}
}

func TestExpectedFailure(t *testing.T) {
	ur := NewUploadRequest(TestToken, TestAppFile)
	ur.ds = diawiTestSerivce{failUpload: true}

	_, err := ur.Upload()
	if err == nil {
		t.Error("Expected error but did not return one")
	}
}

// Integration Tests //
func TestIntegrationUploadSuccess(t *testing.T) {
	if !(*integration) {
		t.Skip()
	}

	ur := NewUploadRequest(*diawiToken, *appFile)

	uploadResponse, err := ur.Upload()
	if err != nil {
		t.Log("Error occured. Did you provide the absolute path to the app file?")
		t.Fatalf("Error: %s", err.Error())
	}

	if uploadResponse.JobIdentifier == "" {
		t.Error("Job identifier blank")
	}
}

func TestIntegrationUploadUnAuthError(t *testing.T) {
	if !(*integration) {
		t.Skip()
	}

	ur := NewUploadRequest("not a legit token", *appFile)

	_, err := ur.Upload()
	if err == nil {
		t.Error("No error occured despite lack of token")
	}
}

func TestIntegrationUploadFileError(t *testing.T) {
	if !(*integration) {
		t.Skip()
	}

	ur := NewUploadRequest(*diawiToken, "Not a legit file name")
	_, err := ur.Upload()
	if err == nil {
		t.Error("No error occured despite lack of file")
	}

}
