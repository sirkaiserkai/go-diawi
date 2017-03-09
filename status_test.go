package godiawi

import (
	//"flag"
	//"fmt"
	//"log"
	//"os"
	"testing"
)

//var TestToken = "12345abcd"
var TestJobIdentifer = "abcd1234"

// Unit Tests
func TestGetStatusNoToken(t *testing.T) {
	sr := StatusRequest{Token: "", JobIdentifier: TestJobIdentifer}
	_, err := sr.GetJobStatus()
	if err == nil {
		t.Error("Should be no token error")
	}
}

func TestGetStatusNoJobId(t *testing.T) {
	sr := StatusRequest{Token: TestToken, JobIdentifier: ""}
	_, err := sr.GetJobStatus()
	if err == nil {
		t.Error("Should be no job ID error")
	}
}

func TestGetStatusSuccess(t *testing.T) {
	sr := NewStatusRequest(TestToken, TestAppFile)
	sr.ds = DiawiTestSerivce{returnFailStatus: false, returnSuccessStatus: true}

	sResponse, err := sr.GetJobStatus()
	if err != nil {
		t.Errorf("Received error: %s", err.Error())
	}

	if sResponse.Status != Ok {
		t.Errorf("Status did not equal 'Ok'. Was instead: %d", sResponse.Status)
	}
}

func TestGetStatusError(t *testing.T) {
	sr := NewStatusRequest(TestToken, TestAppFile)
	sr.ds = DiawiTestSerivce{returnFailStatus: true, returnSuccessStatus: false}

	sResponse, err := sr.GetJobStatus()
	if err != nil {
		t.Errorf("Received error: %s", err.Error())
	}

	if sResponse.Status != ErrorOccured {
		t.Errorf("Status did not equal 'ErrorOccured'. Was instead: %d", sResponse.Status)
	}
}

func TestGetStatusProcessing(t *testing.T) {
	sr := NewStatusRequest(TestToken, TestAppFile)
	sr.ds = DiawiTestSerivce{returnFailStatus: false, returnSuccessStatus: false}

	sResponse, err := sr.GetJobStatus()
	if err != nil {
		t.Errorf("Received error: %s", err.Error())
	}

	if sResponse.Status != Processing {
		t.Errorf("Status did not equal 'Processing'. Was instead: %d", sResponse.Status)
	}
}

func WaitForFinished(t *testing.T, expectedStatus DiawiStatus) {
	sr := NewStatusRequest(TestToken, TestAppFile)
	if expectedStatus == Ok {
		sr.ds = DiawiTestSerivce{returnFailStatus: false, returnSuccessStatus: true}
	}

	if expectedStatus == Processing {
		sr.ds = DiawiTestSerivce{returnFailStatus: false, returnSuccessStatus: false}
	}

	if expectedStatus == ErrorOccured {
		sr.ds = DiawiTestSerivce{returnFailStatus: true, returnSuccessStatus: false}
	}

	sResponse, err := sr.WaitForFinishedStatus()
	if err != nil {
		if sResponse != nil {
			if sResponse.Status != ErrorOccured {
				t.Errorf("Received error: %s", err.Error())
			}
		}
	}

	if sResponse.Status != expectedStatus {
		t.Errorf("Status did not equal '%d'. Was instead: %d", expectedStatus, sResponse.Status)
	}

}

func TestWaitForFinishedStatusSuccess(t *testing.T) {
	WaitForFinished(t, Ok)
}

func TestWaitForFinishedStatusError(t *testing.T) {
	WaitForFinished(t, ErrorOccured)
}

// Integration Tests //
func PerformUpload(t *testing.T) *UploadResponse {
	ur := NewUploadRequest(*diawiToken, *appFile)
	uploadResponse, err := ur.Upload()
	if err != nil {
		t.Log("Error occured. Did you provide the absolute path to the app file?")
		t.Fatalf("Error: %s", err.Error())
	}

	if uploadResponse.JobIdentifier == "" {
		t.Fatalf("Job identifier blank")
	}

	return uploadResponse
}

func TestIntegrationGetStatusSuccess(t *testing.T) {
	if !(*integration) {
		t.Skip()
	}
	ur := PerformUpload(t)
	sr := NewStatusRequest(*diawiToken, ur.JobIdentifier)

	status, err := sr.GetJobStatus()
	if err != nil {
		t.Errorf("Received Error: %s", err.Error())
	}

	if status != nil {
		if status.Status == 0 || status.Status == ErrorOccured {
			t.Errorf("Received Error status: %d", status.Status)
		}
	}
}

func TestIntegrationWaitForFinishedStatusSuccess(t *testing.T) {
	if *integration {
		t.Skip()
	}

	ur := PerformUpload(t)
	sr := NewStatusRequest(*diawiToken, ur.JobIdentifier)

	status, err := sr.WaitForFinishedStatus()
	if err != nil {
		t.Errorf("Received Error: %s", err.Error())
	}

	if status == nil {
		t.Fatal("StatusResponse struct is nil")
	}

	if status.Status == ErrorOccured {
		t.Fatal("StatusResponse status == ErrorOccured (4000)")
	}
}
