package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEvents_stub(t *testing.T) {
	req, err := http.NewRequest("GET", "/events?n=1&filter=down&file=junk.log", nil)
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
	expected := `{"Events":["prefix:./"]}`
	if recorder.Body.String() != expected {
		t.Errorf("GetEvents returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}
