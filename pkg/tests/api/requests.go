package api

import (
	"encoding/json"
	"library/models"
	"net/http/httptest"

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
