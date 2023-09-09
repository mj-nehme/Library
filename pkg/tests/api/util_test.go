package api_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	apiPath      = "/api"
	apiVersionV1 = "v1"
	v1Prefix     = apiPath + "/" + apiVersionV1
)

func sendRequest(t *testing.T, router *gin.Engine, method string, path string, requestBody []byte) *httptest.ResponseRecorder {
	var body io.Reader
	if requestBody == nil {
		body = nil
	} else {
		body = bytes.NewBuffer(requestBody)
	}
	request, err := http.NewRequest(method, path, body)
	assert.NoError(t, err)

	request.Header.Set("Accept-Version", apiVersionV1)
	request.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	response := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(response, request)

	return response
}

func sendRequestV1(t *testing.T, router *gin.Engine, method string, path string, requestBody []byte) *httptest.ResponseRecorder {
	url := v1Prefix + path
	response := sendRequest(t, router, method, url, requestBody)

	return response
}
