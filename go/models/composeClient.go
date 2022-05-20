package models

import (
	"bytes"
	"fmt"
	"net/http"
)

const composeDomain string = "compose" // This is the name of the container in docker.
const composeDomainAWS string = "http://ec2-18-204-20-12.compute-1.amazonaws.com:8010/composeWithWebApp"
const composePort string = "8010"
const composeCreatePath = "composeWithWebApp"

func SendPostAsync(body []byte, rc chan *http.Response) {
	composeURL := "http://" + composeDomain + ":" + composePort + "/" + composeCreatePath
	//
	//writer := multipart.NewWriter(body)
	//
	//writer.Close()
	response, err := http.Post(composeURL, "application/json", bytes.NewReader(body))
	//response, err := http.Post(composeURL, "multipart/form-data", bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	rc <- response
}

func SendPostAsyncWithURL(body []byte, rc chan *http.Response, composeURL string) {
	response, err := http.Post(composeURL, "application/json", bytes.NewReader(body))
	//response, err := http.Post(composeURL, "multipart/form-data", bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	rc <- response
}
