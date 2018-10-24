package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getAPITickerLatest(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPITickerLatest))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}


	req, err = http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusNotFound %d, received %d. ", http.StatusNotFound, resp.StatusCode)
		return
	}

}


func Test_getAPITicker(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPITicker))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}


	req, err = http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusNotFound %d, received %d. ", http.StatusNotFound, resp.StatusCode)
		return
	}

}


func Test_getAPITickerTimestamp(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPITickerTimestamp))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusBadRequest %d, received %d. ", http.StatusBadRequest, resp.StatusCode)
		return
	}

}