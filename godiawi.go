// Package godiawi wraps the diawi api in go
package main

import (
	"log"
)

func main() {
	token := "Your token here"
	file := "DiawiExampleApp.ipa"

	upRequest := UploadRequest{}
	upRequest.Token = token
	upRequest.File = file
	upRequest.FindByUDID = false
	//upRequest.Comment = "Hello, world!"

	upResponse, err := upRequest.Upload()
	if err != nil {
		log.Fatal(err)
	}

	statReq := StatusRequest{}
	statReq.JobIdentifier = upResponse.JobIdentifier
	statReq.Token = token

	sr, err := statReq.WaitForFinishedStatus()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(sr)
}
