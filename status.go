package godiawi

import (
	"errors"
	"fmt"
	"time"
)

// StatusRequest is used to poll diawi to see the status
// of the upload.
type StatusRequest struct {
	// Required parameters
	Token         string
	JobIdentifier string

	ds diawi
}

// NewStatusRequest is the method to be used to create the
// StatusRequest struct
func NewStatusRequest(token, jobIdentifer string) StatusRequest {
	sr := StatusRequest{Token: token, JobIdentifier: jobIdentifer}
	sr.ds = newDiawiService()

	return sr
}

var EmptyTokenField = errors.New("Token field is blank")
var EmptyJobField = errors.New("Job identifier field is blank")

// GetJobStatus makes a status request to see the current progress
// for the app upload.
func (s *StatusRequest) GetJobStatus() (*StatusResponse, error) {

	if s.Token == "" {
		return nil, EmptyTokenField
	}

	if s.JobIdentifier == "" {
		return nil, EmptyJobField
	}

	statusRes := StatusResponse{}
	s.ds.GetStatus(s.Token, s.JobIdentifier, &statusRes)

	return &statusRes, nil
}

var MaxPollsReachedError = errors.New("Exceeded max number of polls to get upload status.")
var UnknownStatusError = errors.New("Unknown status error")

// WaitForFinishedStatus continually polls diawi using GetJobStatus
// until the service provides a DiawiStatus other than
// Processing (2001)
func (s *StatusRequest) WaitForFinishedStatus() (*StatusResponse, error) {
	sr, err := s.GetJobStatus()
	if err != nil {
		return nil, err
	}

	numberOfPolls := 0
	for {
		switch sr.Status {
		case Processing:

			if numberOfPolls > StatusPollingMax {
				return nil, MaxPollsReachedError
			}

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
			return sr, UnknownStatusError
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
