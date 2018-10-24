package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_igcInfo(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(igcInfo))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, _ := client.Do(req)

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Expected StatusFound %d, received %d. ", http.StatusFound, resp.StatusCode)
		return
	}

	// Testing the Status Not Implemented Yet
	req, err = http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the Delete request, %s", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Error executing the Delete request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusNotImplemented, resp.StatusCode)
		return
	}

}


func Test_getAPI_NotImplemented(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPI))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusNotImplemented, resp.StatusCode)
		return
	}
}

func Test_getAPIIgc_NotImplemented(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPIIgc))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusNotImplemented, resp.StatusCode)
		return
	}

}

func Test_getAPIIgcId_NotImplemented(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPIIgcID))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusNotImplemented, resp.StatusCode)
		return
	}

}

func Test_getAPIIgcField_NotImplemented(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPIIgcField))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusNotImplemented, resp.StatusCode)
		return
	}

}


func Test_getAPI_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPI))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/something/",
		ts.URL + "/something/123/",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}

func Test_getAPIIgc_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPIIgc))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/something/",
		ts.URL + "/something/123/",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}

func Test_getAPIIgcId_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPIIgcID))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/something/",
		ts.URL + "/something/123/",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}

func Test_getAPIIgcField_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPIIgcField))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/something/",
		ts.URL + "/something/123/",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}


func Test_getAPIIgc_Post(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPIIgc))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	apiURLTest := URLStruct{}
	apiURLTest.URL = "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"

	jsonData, _ := json.Marshal(apiURLTest)

	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Error making the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	if resp.StatusCode == 200 {
		assert.Equal(t, 200 , resp.StatusCode, "OK response is expected")
	} else {
		assert.Equal(t, 409 , resp.StatusCode, "OK response is expected")
	}

}


func Test_getAPIIgc_Post_Empty(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPIIgc))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	apiURLTest := URLStruct{}
	apiURLTest.URL = ""

	jsonData, _ := json.Marshal(apiURLTest)

	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Error making the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	assert.Equal(t, 400, resp.StatusCode, "OK response is expected")

}


//func Test_getAPIIgc_Get(t *testing.T) {
//
//	ts := httptest.NewServer(http.HandlerFunc(getAPIIgc))
//	defer ts.Close()
//
//	//create a request to our mock HTTP server
//	client := &http.Client{}
//
//	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
//	if err != nil {
//		t.Errorf("Error making the GET request, %s", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Errorf("Error executing the GET request, %s", err)
//	}
//
//	countAtFunction := make([]string, 0, 0)
//
//	json.NewDecoder(resp.Body).Decode(&countAtFunction)
//
//	assert.Equal(t, 200, resp.StatusCode, "OK response is expected")
//
//
//	clientDB := mongoConnect()
//	trackCount := int(countAllTracks(clientDB))
//
//	assert.Equal(t, trackCount, len(countAtFunction), "Expected same number of tracks")
//
//}