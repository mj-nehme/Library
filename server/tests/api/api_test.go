package api_test

import (
	"bytes"
	"fmt"
	"library/tests"
	"library/tests/api"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	ctx, router, db := tests.SetupMockServer()
	defer tests.TearDownMockServer(ctx, db)

	method := "GET"
	url := "/health"
	var body []byte = nil
	response, err := api.SendRequest(router, method, url, body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Code, "Unexpected status code for health check")
}

func TestVersionHandler(t *testing.T) {
	ctx, router, db := tests.SetupMockServer()
	defer tests.TearDownMockServer(ctx, db)

	testCases := []struct {
		Version     string
		Expectation int
	}{
		{Version: "1", Expectation: http.StatusOK},
		{Version: "2", Expectation: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Version %s.0", tc.Version), func(t *testing.T) {
			method := "GET"
			url := fmt.Sprintf("/api/v%s/", tc.Version)
			var body []byte = nil
			response, err := api.SendRequest(router, method, url, body)
			assert.NoError(t, err)

			assert.Equal(t, tc.Expectation, response.Code, "Unexpected status code for version: %s", tc.Version)
		})
	}
}

func TestWelcomePageHandler(t *testing.T) {
	ctx, router, db := tests.SetupMockServer()
	defer tests.TearDownMockServer(ctx, db)

	// Create a new HTTP request to the root path "/"
	method := "GET"
	url := "/"
	var body []byte = nil
	contentType := "text/html"
	response, err := api.SendRequest(router, method, url, body, contentType)
	assert.NoError(t, err)

	// Check the response status code
	assert.Equal(t, http.StatusOK, response.Code, "Unexpected status code")

	// Check the Content-Type header
	assert.Equal(t, "text/html; charset=utf-8", response.Header().Get("Content-Type"), "Unexpected Content-Type")

	// Read the response body into a buffer
	var buf bytes.Buffer
	_, err = buf.ReadFrom(response.Body)
	assert.NoError(t, err)

	// Convert the buffer to a string
	responseBody := buf.String()

	// Verify the response body contains the expected message
	expectedMessage := "Welcome to the Book Management API!"
	assert.Contains(t, responseBody, expectedMessage, "Unexpected message in response body")
}
