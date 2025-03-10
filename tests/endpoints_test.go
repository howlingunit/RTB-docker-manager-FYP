package main_test

import (
	"encoding/json"
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

type Challenge struct {
	Name       string `json:"name"`
	Difficulty string `json:"difficulty"`
}

func TestGetChallenges(t *testing.T) {
	res, err := http.Get("http://localhost:8080/get-challenges")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var challenges []Challenge
	err = json.NewDecoder(res.Body).Decode(&challenges)
	if err != nil {
		t.Fatal(err)
	}
	// Check the response body
	expected := Challenge{Name: "sambaCry", Difficulty: "Easy"}
	assert.Contains(t, challenges, expected)
}
