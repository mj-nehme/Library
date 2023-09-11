package api_test

import (
	"encoding/json"
	"fmt"
	"library/models"
	"library/tests"
	"library/tests/api"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestAddBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a valid book JSON
	sampleBook, err := api.LoadSampleBook()
	if err != nil {
		slog.Error("Unable to load sample book. " + err.Error())
	}

	testCases := []struct {
		Description string
		Book        models.Book
		Expected    int // Expected HTTP status code
	}{
		{
			Description: "Add Valid Book",
			Book:        sampleBook,
			Expected:    http.StatusCreated,
		},
		{
			Description: "Add Invalid Book",
			Book:        models.Book{},
			Expected:    http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			response, err := api.SendAddBookRequest(router, &tc.Book)
			assert.NoError(t, err)

			// Check the response status code
			assert.Equal(t, tc.Expected, response.Code, "Expected status code %d, but got %d", tc.Expected, response.Code)
		})
	}
}

func TestGetBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	book := api.CreateBookTemplate(t, router)

	testCases := []struct {
		Description string
		BookID      uint
		Expected    int // Expected HTTP status code
	}{
		{
			Description: "Get Existing Book",
			BookID:      book.ID,
			Expected:    http.StatusOK,
		},
		{
			Description: "Get Non-Existent Book",
			BookID:      book.ID + 1, // Use a non-existent ID
			Expected:    http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			response, err := api.SendGetBookRequest(router, tc.BookID)
			assert.NoError(t, err)

			// Check the response status code
			assert.Equal(t, tc.Expected, response.Code, "Expected status code %d, but got %d", tc.Expected, response.Code)

			if tc.Expected == http.StatusOK {
				// Read the response body and unmarshal it into a book
				var responseBook models.Book
				err := json.Unmarshal(response.Body.Bytes(), &responseBook)
				if err != nil {
					t.Fatalf("Failed to unmarshal response JSON: %v", err)
				}

				// Convert expectedTime to UTC timezone
				actualTimeUTC := book.Published.In(time.UTC).Round(time.Hour)
				expectedTimeUTC := responseBook.Published.In(time.UTC).Round(time.Hour)

				// Add assertions to verify that the response book matches the expected book
				assert.Equal(t, book.Title, responseBook.Title, "Title mismatch")
				assert.Equal(t, book.Author, responseBook.Author, "Author mismatch")
				assert.Equal(t, expectedTimeUTC, actualTimeUTC, "PublishedDate mismatch")
				assert.Equal(t, book.Edition, responseBook.Edition, "Edition mismatch")
				assert.Equal(t, book.GenreName, responseBook.GenreName, "Genre mismatch")
			}
		})
	}
}

func TestListBooksHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	books := api.CreateListOfBookTemplates(t, router)

	// Perform a GET request to the "ListBooks" endpoint
	response, err := api.SendListBooksRequest(router)
	assert.NoError(t, err)

	// Check the response status code
	assert.Equal(t, http.StatusOK, response.Code, "Expected status code 200, but got %d", response.Code)

	listOfBooks := []models.Book{}
	err = json.Unmarshal(response.Body.Bytes(), &listOfBooks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}
	assert.Len(t, books, len(listOfBooks))
	for index := range books {
		listOfBooks[index].Published = listOfBooks[index].Published.UTC().Round(time.Hour)
		books[index].Published = books[index].Published.UTC().Round(time.Hour)
		assert.Equal(t, books[index], listOfBooks[index])
	}
}

func TestUpdateBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample book in the database for testing
	book := api.CreateBookTemplate(t, router)

	updatedBook := api.CopyBook(&book)
	updatedTitle := "Updated Book Title"
	updatedAuthor := "Updated Book Author"
	updatedBook.Title = updatedTitle
	updatedBook.Author = updatedAuthor

	updatedBookTitle := api.CopyBook(updatedBook)
	updatedBookTitle.Title = updatedTitle
	updatedBookTitle.Author = book.Author

	nonExistingBook := api.CopyBook(updatedBookTitle)
	nonExistingBook.ID = book.ID - 1

	// Define test cases for updating books
	testCases := []struct {
		Description      string
		UpdatedBook      models.Book
		ExpectedTitle    string
		ExpectedAuthor   string
		ExpectedHTTPCode int
		ExpectedDBTitle  string
		ExpectedDBAuthor string
		ShouldFail       bool // Whether the update should fail
	}{
		{
			Description:      "Update Book Title and Author",
			UpdatedBook:      *updatedBook,
			ExpectedTitle:    updatedTitle,
			ExpectedAuthor:   updatedAuthor,
			ExpectedHTTPCode: http.StatusOK,
			ExpectedDBTitle:  updatedTitle,
			ExpectedDBAuthor: updatedAuthor,
			ShouldFail:       false,
		},
		{
			Description:      "Update Book Title Only",
			UpdatedBook:      *updatedBookTitle,
			ExpectedTitle:    updatedTitle,
			ExpectedAuthor:   book.Author, // Author should remain unchanged
			ExpectedHTTPCode: http.StatusOK,
			ExpectedDBTitle:  updatedTitle,
			ExpectedDBAuthor: book.Author, // Author in the database should remain unchanged
			ShouldFail:       false,
		},
		{
			Description:      "Update Non-Existent Book",
			UpdatedBook:      *nonExistingBook,
			ExpectedHTTPCode: http.StatusNotFound,
			ShouldFail:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			response, err := api.SendUpdateBookRequest(router, &tc.UpdatedBook)
			assert.NoError(t, err)

			// Check the response status code
			assert.Equal(t, tc.ExpectedHTTPCode, response.Code, "Expected status code %d, but got %d", tc.ExpectedHTTPCode, response.Code)

			if !tc.ShouldFail && tc.ExpectedHTTPCode == http.StatusOK {
				// Read the response body
				var responseBook models.Book
				err = json.Unmarshal(response.Body.Bytes(), &responseBook)
				if err != nil {
					t.Fatalf("Failed to unmarshal response JSON: %v", err)
				}

				// Verify that the response book matches the expected updated book data
				assert.Equal(t, tc.ExpectedTitle, responseBook.Title, "Title mismatch")
				assert.Equal(t, tc.ExpectedAuthor, responseBook.Author, "Author mismatch")

				// Fetch the book from the database to ensure it was updated
				response, err := api.SendGetBookRequest(router, tc.UpdatedBook.ID)
				assert.NoError(t, err)

				// Check the response status code
				assert.Equal(t, http.StatusOK, response.Code, "Expected status code %d, but got %d", http.StatusOK, response.Code)

				// Read the response body and unmarshal it into a book
				var updatedBookFromDB models.Book
				err = json.Unmarshal(response.Body.Bytes(), &updatedBookFromDB)
				assert.NoError(t, err)

				// Verify that the book in the database matches the expected updated book data
				assert.Equal(t, tc.ExpectedDBTitle, updatedBookFromDB.Title, "Title mismatch in the database")
				assert.Equal(t, tc.ExpectedDBAuthor, updatedBookFromDB.Author, "Author mismatch in the database")
			}
		})
	}
}

func TestPatchBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample book in the database for testing
	book := api.CreateBookTemplate(t, router)

	updatedBook := api.CopyBook(&book)
	updatedTitle := "Updated Book Title"
	updatedAuthor := "Updated Book Author"
	updatedBook.Title = updatedTitle
	updatedBook.Author = updatedAuthor

	updatedBookTitle := api.CopyBook(updatedBook)
	updatedBookTitle.Title = updatedTitle
	updatedBookTitle.Author = book.Author

	nonExistingBook := api.CopyBook(updatedBookTitle)
	nonExistingBook.ID = book.ID + 1
	updatedBookTitle.Title = book.Title
	updatedBookTitle.Author = book.Author

	// Define test cases for patching books
	testCases := []struct {
		Description      string
		PatchedBook      models.Book
		ExpectedHTTPCode int
		ShouldFail       bool // Whether the patch should fail
	}{
		{
			Description:      "Patch Book Title and Author",
			PatchedBook:      *updatedBook,
			ExpectedHTTPCode: http.StatusOK,
			ShouldFail:       false,
		},
		{
			Description:      "Patch Book Title Only",
			PatchedBook:      *updatedBookTitle,
			ExpectedHTTPCode: http.StatusOK,
			ShouldFail:       false,
		},
		{
			Description:      "Patch Non-Existent Book",
			PatchedBook:      *nonExistingBook,
			ExpectedHTTPCode: http.StatusNotFound,
			ShouldFail:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			response, err := api.SendPatchBookRequest(router, &tc.PatchedBook)
			assert.NoError(t, err)

			// Check the response status code
			assert.Equal(t, tc.ExpectedHTTPCode, response.Code, "Expected status code %d, but got %d", tc.ExpectedHTTPCode, response.Code)

			if !tc.ShouldFail && tc.ExpectedHTTPCode == http.StatusOK {
				// Read the response body
				var responseBook models.Book
				err = json.Unmarshal(response.Body.Bytes(), &responseBook)
				if err != nil {
					t.Fatalf("Failed to unmarshal response JSON: %v", err)
				}

				// Verify that the response book matches the expected patched book data
				assert.Equal(t, tc.PatchedBook.Title, responseBook.Title, "Title mismatch")
				assert.Equal(t, tc.PatchedBook.Author, responseBook.Author, "Author mismatch")

				// Fetch the book from the database to ensure it was patched
				response, err := api.SendGetBookRequest(router, tc.PatchedBook.ID)
				assert.NoError(t, err)

				// Check the response status code
				assert.Equal(t, http.StatusOK, response.Code, "Expected status code %d, but got %d", http.StatusOK, response.Code)

				// Read the response body and unmarshal it into a book
				var patchedBookFromDB models.Book
				err = json.Unmarshal(response.Body.Bytes(), &patchedBookFromDB)
				assert.NoError(t, err)

				// Verify that the book in the database matches the expected patched book data
				assert.Equal(t, tc.PatchedBook.Title, patchedBookFromDB.Title, "Title mismatch in the database")
				assert.Equal(t, tc.PatchedBook.Author, patchedBookFromDB.Author, "Author mismatch in the database")
			}
		})
	}
}

func TestDeleteBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample book in the database for testing
	book := api.CreateBookTemplate(t, router)

	// Define test cases for deleting books
	testCases := []struct {
		Description      string
		BookID           uint
		ExpectedHTTPCode int
		ShouldBeDeleted  bool
	}{
		{
			Description:      "Delete Existing Book",
			BookID:           book.ID,
			ExpectedHTTPCode: http.StatusOK,
			ShouldBeDeleted:  false,
		},
		{
			Description:      "Delete Invalid Book",
			BookID:           book.ID + 1, // Use an invalid ID
			ExpectedHTTPCode: http.StatusNotFound,
			ShouldBeDeleted:  false, // Invalid book should still not exist
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			response, err := api.SendDeleteBookRequest(router, tc.BookID)
			assert.NoError(t, err)

			// Check the response status code
			assert.Equal(t, tc.ExpectedHTTPCode, response.Code, "Expected status code %d, but got %d", tc.ExpectedHTTPCode, response.Code)

			// Verify the book's existence in the database
			response, err = api.SendGetBookRequest(router, tc.BookID)
			assert.Equal(t, http.StatusNotFound, response.Code, "Expected status code %d, but got %d", tc.ExpectedHTTPCode, response.Code)

			if tc.ShouldBeDeleted {
				assert.Error(t, err, "Expected book to be deleted from the database")
			} else {
				assert.NoError(t, err, "Expected book to exist in the database")

			}
		})
	}
}

func TestSearchBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	books := api.CreateListOfBookTemplates(t, router)

	// Create test cases for different search queries
	testCases := []struct {
		Description      string
		QueryParams      map[string]string
		ExpectedHTTPCode int
		ExpectedCount    int
	}{
		{
			Description: "Search by Author and Genre",
			QueryParams: map[string]string{
				"author": books[0].Author,
				"genre":  books[0].GenreName,
			},
			ExpectedHTTPCode: http.StatusOK,
			ExpectedCount:    1,
		},
		{
			Description: "Search by Author Only",
			QueryParams: map[string]string{
				"author": books[0].Author,
			},
			ExpectedHTTPCode: http.StatusOK,
			ExpectedCount:    1, // Adjust the expected count as needed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			// Create a URL with query parameters
			values := url.Values{}
			for key, value := range tc.QueryParams {
				values.Add(key, value)
			}
			query := values.Encode()
			fmt.Println("Query: ", query)

			response, err := api.SendSearchBooksRequest(router, query)
			assert.NoError(t, err)

			// Check the response status code
			assert.Equal(t, tc.ExpectedHTTPCode, response.Code, "Expected status code %d, but got %d", tc.ExpectedHTTPCode, response.Code)

			// Read the response body
			var responseBooks []models.Book
			err = json.Unmarshal(response.Body.Bytes(), &responseBooks)
			if err != nil {
				t.Fatalf("Failed to unmarshal response JSON: %v", err)
			}

			// Verify the number of books in the response
			assert.Len(t, responseBooks, tc.ExpectedCount, "Expected %d books in the response", tc.ExpectedCount)

			// Add assertions to compare the response books with the expected books
			for i, expectedBook := range books {
				if i >= len(responseBooks) {
					break // Avoid index out of range error
				}
				expectedTimeUTC := expectedBook.Published.In(time.UTC).Round(time.Hour)
				publishedTimeUTC := responseBooks[i].Published.In(time.UTC).Round(time.Hour)

				assert.Equal(t, expectedBook.Title, responseBooks[i].Title, "Title mismatch")
				assert.Equal(t, expectedBook.Author, responseBooks[i].Author, "Author mismatch")
				assert.Equal(t, expectedBook.Edition, responseBooks[i].Edition, "Edition mismatch")
				assert.Equal(t, expectedTimeUTC, publishedTimeUTC, "Published mismatch")
				assert.Equal(t, expectedBook.Description, responseBooks[i].Description, "Description mismatch")
				assert.Equal(t, expectedBook.GenreName, responseBooks[i].GenreName, "Genre mismatch")
			}
		})
	}
}

func TestCountBooksHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	books := api.CreateListOfBookTemplates(t, router)

	// Define test cases for counting books
	testCases := []struct {
		Description      string
		ExpectedHTTPCode int
		ExpectedCount    int64
	}{
		{
			Description:      "Count All Books",
			ExpectedHTTPCode: http.StatusOK,
			ExpectedCount:    int64(len(books)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			response, err := api.SendCountBooksRequest(router)
			assert.NoError(t, err)

			// Check the response status code
			assert.Equal(t, tc.ExpectedHTTPCode, response.Code, "Expected status code %d, but got %d", tc.ExpectedHTTPCode, response.Code)

			// Parse the response JSON to get the count
			var responseCount int64 // Change the type to int64
			err = json.NewDecoder(response.Body).Decode(&responseCount)
			assert.NoError(t, err)

			// Check if the response contains the correct count
			assert.Equal(t, tc.ExpectedCount, responseCount, "Unexpected count response")
		})
	}
}
