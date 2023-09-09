package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

const (
	apiPath      = "/api"
	apiVersionV1 = "v1"
	v1Prefix     = apiPath + "/" + apiVersionV1
)

func SendRequest(router *gin.Engine, method string, path string, requestBody []byte) (*httptest.ResponseRecorder, error) {
	var body io.Reader
	if requestBody == nil {
		body = nil
	} else {
		body = bytes.NewBuffer(requestBody)
	}
	request, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept-Version", apiVersionV1)
	request.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	response := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(response, request)

	return response, nil
}

func SendRequestV1(router *gin.Engine, method string, path string, requestBody []byte) (*httptest.ResponseRecorder, error) {
	url := v1Prefix + path

	return SendRequest(router, method, url, requestBody)
}
