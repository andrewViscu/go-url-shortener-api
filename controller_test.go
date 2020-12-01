package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"testing"
	"log"
	"io/ioutil"

)

type TestServer struct {
	Server *httptest.Server
	Client *http.Client
}

var ts TestServer

func TestDBConnect(*testing.T) {
	handler := ConfigureServer()
	ts.Server = httptest.NewServer(handler)
	ts.Client = ts.Server.Client()

}
func TestShortenURL(*testing.T) {
	var url URL
	jsonValue := []byte(`{"original":"http://google.com"}`)
	request, err := http.NewRequest("POST", ts.Server.URL, bytes.NewBuffer(jsonValue))
	resp, err := ts.Client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Print("Response status: ", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(body))
	json.Unmarshal(body, &url)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(url)

	log.Print("Shorten url: " + url.Shorten)

	request, err = http.NewRequest("GET", ts.Server.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = ts.Client.Do(request)
	log.Print("Response status: ", resp.Status)
	body, _ = ioutil.ReadAll(resp.Body)
	log.Print("Response body:", string(body))
	
	request, err = http.NewRequest("POST", ts.Server.URL + "/" + url.Shorten, nil)
	resp, err = ts.Client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v, %v, %v", resp.Request.URL, resp.Request.Method, resp.Status)
}