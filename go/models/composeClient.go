package models

import (
	"bytes"
	"fmt"
	"net/http"
)

const composeDomain string = "compose"
const composePort string = "5000"
const composeCreatePath = "/composeWithWebApp"

func SendPostAsync(body []byte, rc chan *http.Response) {
	composeURL := "http://" + composeDomain + ":" + composePort + composeCreatePath
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