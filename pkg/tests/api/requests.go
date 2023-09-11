package api

import (
	"encoding/json"
	"fmt"
	"library/models"
	"net/http/httptest"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func SendAddBookRequest(router *gin.Engine, book *models.Book) (*httptest.ResponseRecorder, error) {
	jsonData, err := json.Marshal(book)
	if err != nil {
		slog.Error("Unable to marshal book in JSON")
	}

	method := "POST"
	url := "/books"
	body := jsonData
	return SendRequestV1(router, method, url, body)
}

func SendGetBookRequest(router *gin.Engine, ID uint) (*httptest.ResponseRecorder, error) {
	method := "GET"
	url := "/books/" + strconv.Itoa(int(ID))
	var body []byte = nil
	return SendRequestV1(router, method, url, body)
}

func SendListBooksRequest(router *gin.Engine) (*httptest.ResponseRecorder, error) {
	method := "GET"
	url := "/books"
	var body []byte = nil
	return SendRequestV1(router, method, url, body)
}

func SendUpdateBookRequest(router *gin.Engine, book *models.Book) (*httptest.ResponseRecorder, error) {
	jsonData, err := json.Marshal(book)
	if err != nil {
		slog.Error("Unable to marshal book in JSON")
	}

	method := "PUT"
	url := fmt.Sprintf("/books/%d", book.ID)
	body := jsonData
	return SendRequestV1(router, method, url, body)
}

func SendDeleteBookRequest(router *gin.Engine, ID uint) (*httptest.ResponseRecorder, error) {
	// Perform a DELETE request to the "DeleteBook" endpoint with the book ID
	method := "DELETE"
	url := fmt.Sprintf("/books/%d", ID)
	var body []byte = nil
	return SendRequestV1(router, method, url, body)
}
