package main

import (
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

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusNotFound %d, received %d. ", http.StatusNotFound, resp.StatusCode)
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


// ***TESTING THE STATUS NOT IMPLEMENTED FOR EACH FUNCTION*** //

func Test_getApi_NotImplemented(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getApi))
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


func Test_getApiIgc_NotImplemented(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getApiIgc))
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


func Test_getApiIgcId_NotImplemented(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getApiIgcId))
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


func Test_getApiIgcField_NotImplemented(t *testing.T) {


	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getApiIgcField))
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



// ***TESTING MALFORMED URL FOR EACH FUNCTION*** //

// NOTE: I believe that while using the Gorilla Mux Router it handles itself the malformed URL by returning 404 File Not Found status code.
// 		 But anyways i tried to do some testing about it.

func Test_getApi_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getApi))
	defer ts.Close()

	testCases := []string {
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


func Test_getApiIgc_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getApiIgc))
	defer ts.Close()

	testCases := []string {
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


func Test_getApiIgcId_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getApiIgcId))
	defer ts.Close()

	testCases := []string {
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


func Test_getApiIgcField_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getApiIgcField))
	defer ts.Close()

	testCases := []string {
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



/*
// Testing the getApiIgc function for getting an empty array
func Test_getApiIgc_Empty(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getApiIgc))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}

	var a []interface{}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	if len(a) != 0 {
		t.Errorf("Expected empty array, got %s", a)
	}
}
*/
