package main_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	expected := Challenge{Name: "sambacry", Difficulty: "Easy"}
	assert.Contains(t, challenges, expected)
}

type DockerResponse struct {
	Name string `json:"name"`
	Flag string `json:"flag"`
	Ip   string `json:"ip"`
}

func TestCreateChallenges(t *testing.T) {
	data := `[{"name":"sambacry", "flag":"1234"},{"name":"dvwa","flag":"1234"}]` // JSON body as a string
	req, err := http.NewRequest("POST", "http://localhost:8080/create-challenges", strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var dockers []DockerResponse
	err = json.NewDecoder(res.Body).Decode(&dockers)
	if err != nil {
		t.Fatal(err)
	}

	expected := DockerResponse{Name: "sambacry", Flag: "1234", Ip: "172.18.0.2"}

	assert.Contains(t, dockers, expected)

}

func TestRemoveChallenges(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/remove-challenges", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	assert.Equal(t, "Removed Challenges", res.Body)
}

type CreatePlatformResponse struct {
	User string `json:"user"`
	Ip   string `json:"ip"`
}

func TestCreatePlatforms(t *testing.T) {
	data := `["1234", "5678"]`
	req, err := http.NewRequest("POST", "http://localhost:8080/create-platforms", strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var platforms []CreatePlatformResponse
	err = json.NewDecoder(res.Body).Decode(&platforms)
	if err != nil {
		t.Fatal(err)
	}

	expected := []CreatePlatformResponse{
		{User: "1234", Ip: "172.18.0.3"},
		{User: "5678", Ip: "172.18.0.4"},
	}

	assert.Contains(t, platforms, expected)
}

type GetPlatformResponse struct {
	User string `json:"user"`
	Ip   string `json:"ip"`
}

func TestGetPlatform(t *testing.T) {
	userId := "1234"
	res, err := http.Get(fmt.Sprint("http://localhost:8080/get-platform/", userId))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var platform []GetPlatformResponse
	err = json.NewDecoder(res.Body).Decode(&platform)
	if err != nil {
		t.Fatal(err)
	}
	// Check the response body
	expected := GetPlatformResponse{User: "1234", Ip: "172.18.0.3"}
	assert.Contains(t, platform, expected)
}

func TestRemovePlatforms(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/remove-platforms", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	assert.Equal(t, "Removed Platforms", res.Body)
}