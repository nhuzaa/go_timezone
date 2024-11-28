package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCurrentTime(t *testing.T) {
	req, err := http.NewRequest("GET", "/current-time", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getCurrentTime)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expected)
	}
}
