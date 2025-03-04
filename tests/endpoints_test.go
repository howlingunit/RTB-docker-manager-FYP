package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTest(t *testing.T) {
	// Create a request to the endpoint
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the response body
	expected := "hi"
	if rr.Body.String() != expected {
		t.Errorf("expected response body %q, got %q", expected, rr.Body.String())
	}
}
