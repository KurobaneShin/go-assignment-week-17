package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// func TestHandleMemoryObjectGet(t *testing.T) {
// 	// Create a new server
// 	server := NewServer(":3334")
//

//
// 	// Create a ResponseRecorder to record the response
// 	rr := httptest.NewRecorder()
//
// 	// Call the handler function with the GET request and ResponseRecorder
// 	server.inMemoryHandler(rr, req)
//
// 	// Check the status code is what we expect
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// Check the response body is what we expect
// expected := `{"key":"testKey","value":"testValue"}`
// if rr.Body.String() != expected {
// 	t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
// }
// }

// You can add more tests for handleMemoryObjectPost and other handlers similarly.

func TestHandleMemoryObjectPost(t *testing.T) {
	// Create a new server
	server := NewServer(":3334")

	// Create a sample MoResponse object
	obj := MoResponse{
		Key:   "testKey",
		Value: "testValue",
	}
	// Encode the object to JSON
	objJSON, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}

	// Create a POST request with the JSON payload
	req, err := http.NewRequest("POST", "/in-memory", bytes.NewBuffer(objJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Set content type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function with the POST request and ResponseRecorder
	server.inMemoryHandler(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Create a request with "key" query parameter
	getReq, err := http.NewRequest("GET", "/in-memory?key=testKey", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	server.inMemoryHandler(rr2, getReq)

	// Check the status code is what we expect
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"key":"testKey","value":"testValue"}`
	if strings.TrimSpace(rr2.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr2.Body.String(), expected)
	}
}

func TestFetchDataHandler(t *testing.T) {
	// Create a new server
	server := NewServer(":3334")

	// Create a sample request payload
	payload := FetchDataParams{
		StartDate: "2016-01-01",
		EndDate:   "2022-01-31",
		MinCount:  2000,
		MaxCount:  3000,
	}
	// Encode the payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	// Create a POST request with the JSON payload
	req, err := http.NewRequest("POST", "/fetch-data", bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Set content type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function with the POST request and ResponseRecorder
	server.fetchDataHandler(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response body
	var response ResponsePayload
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response payload is what we expect
	expected := ResponsePayload{
		Code: 0,
		Msg:  "success",
		Records: []MongoRecord{
			{
				Key:        "exampleKey",
				CreatedAt:  time.Now(),
				TotalCount: 50,
			},
			// Add more expected records as needed
		},
	}

	// Compare response payload with expected payload
	if !reflect.DeepEqual(response.Msg, expected.Msg) {
		t.Errorf("handler returned unexpected response: got %v want %v", response, expected)
	}
}
