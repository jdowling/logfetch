package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEvents_lastLineOnly(t *testing.T) {
	req, err := http.NewRequest("GET", "/events?n=1&filter=LOG&file=test.log", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	server := Server{"./"}
	handler := http.HandlerFunc(server.GetEvents)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("GetEvents returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"Events":["2019-05-14 15:53:16.215 EDT [609] LOG:  database system is shut down"]}`
	if recorder.Body.String() != expected {
		t.Errorf("GetEvents returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}

func TestGetEvents_filterWorks(t *testing.T) {
	req, err := http.NewRequest("GET", "/events?n=1&filter=aborting&file=test.log", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	server := Server{"./"}
	handler := http.HandlerFunc(server.GetEvents)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("GetEvents returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"Events":["2019-05-14 15:53:16.131 EDT [609] LOG:  aborting any active transactions"]}`
	if recorder.Body.String() != expected {
		t.Errorf("GetEvents returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}

func TestGetEvents_pathologicalFilter(t *testing.T) {
	req, err := http.NewRequest("GET", "/events?n=1&filter=\\L&file=test.log", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	server := Server{"./"}
	handler := http.HandlerFunc(server.GetEvents)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("GetEvents returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	expected := `Invalid regex:\L`
	if recorder.Body.String() != expected {
		t.Errorf("GetEvents returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}

func TestGetEvents_defaultPrefix(t *testing.T) {
	// TODO: works on local windows machine but it's trying to do real IO
	// so this might pass on a linux host. Figure out how to properly mock
	// out a filesystem in Go. Maybe something like:
	// https://stackoverflow.com/questions/16742331/how-to-mock-abstract-filesystem-in-go
	req, err := http.NewRequest("GET", "/events?n=1&filter=down&file=test.log", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	server := NewServer()
	handler := http.HandlerFunc(server.GetEvents)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("GetEvents returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Check the response body is what we expect.
	expected := `Error opening file:\var\log\test.log`
	if recorder.Body.String() != expected {
		t.Errorf("GetEvents returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}
