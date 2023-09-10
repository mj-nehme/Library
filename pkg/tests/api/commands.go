package api

import (
	"encoding/json"
	"library/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func CreateBookTemplate(t *testing.T, router *gin.Engine) models.Book {
	// Create a sample book in the database for testing
	book, err := LoadSampleBook()
	assert.NoError(t, err)

	response, err := SendAddBookRequest(router, &book)
	assert.NoError(t, err)

	// Read the response body and unmarshal it into a book
	var createdBook models.Book
	err = json.Unmarshal(response.Body.Bytes(), &createdBook)
	if err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}

	return createdBook
}

func CreateListOfBookTemplates(t *testing.T, router *gin.Engine) []models.Book {
	// Load a list of sample books
	sampleBooks, err := LoadListOfBookSamples()
	assert.NoError(t, err)

	createdBooks := []models.Book{}
	for _, book := range sampleBooks {
		response, err := SendAddBookRequest(router, &book)
		assert.NoError(t, err)

		// Read the response body and unmarshal it into a book
		var createdBook models.Book
		err = json.Unmarshal(response.Body.Bytes(), &createdBook)
		if err != nil {
			t.Fatalf("Failed to unmarshal response JSON: %v", err)
		}
		createdBooks = append(createdBooks, createdBook)
	}

	return createdBooks
}
