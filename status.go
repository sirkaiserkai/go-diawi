package godiawi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// StatusRequest is used to poll diawi to see the status
// of the upload.
type StatusRequest struct {
	// Required parameters
	Token         string
	JobIdentifier string
}

func (s *StatusRequest) GetJobStatus() (*StatusResponse, error) {

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

func (s *StatusRequest) WaitForFinishedStatus() (*StatusResponse, error) {
	sr, err := s.GetJobStatus()
	if err != nil {
		return nil, err
	}

	for {

		switch sr.Status {
		case Processing:
			time.Sleep(1 * time.Second) // diawi documentation recommends 1 second between each status request.
			sr, err = s.GetJobStatus()
			if err != nil {
				return nil, err
			}
		case Ok:
			return sr, nil
		case ErrorOccured:
			return sr, fmt.Errorf("Response included error status = 4000. sr: %s", sr.String())
		default:
			return sr, fmt.Errorf("Unknown status error")
		}
	}
}

// StatusResponse is the response provided by the
type StatusResponse struct {
	Status  DiawiStatus
	Message string

	// Only provided when upload successful
	Hash string
	Link string
}

func (s *StatusResponse) String() string {
	return fmt.Sprintf("status: %d messsage: %s hash: %s link: %s", s.Status, s.Message, s.Hash, s.Link)
}
