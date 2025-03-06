package main_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest(t *testing.T) {
	// Create a request to the endpoint
	res, err := http.Get("http://localhost:8080/test") // Real API call
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	// Check the response body
	expected := "hi"
	assert.Equal(t, expected, string(body))
}

