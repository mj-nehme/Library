package api_test

import (
	"encoding/json"
	"library/models"
	"library/tests"
	"library/tests/api"
	"net/http"
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
		assert.Equal(t, books[index], listOfBooks[index])
	}
}

/*
func TestUpdateBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample book in the database for testing
	book, err := api.LoadSampleBook()
	assert.NoError(t, err)

	err = core.AddBook(book)
	assert.NoError(t, err)

	book, err = core.GetBook(book.ID)
	assert.NoError(t, err)
	assert.Positive(t, book.ID)

	updatedBook, err := core.CopyBook(book)
	assert.NoError(t, err)
	updatedTitle := "Updated Book Title"
	updatedAuthor := "Updated Book Author"
	updatedBook.Title = updatedTitle
	updatedBook.Author = updatedAuthor

	updatedBookTitle, err := core.CopyBook(updatedBook)
	assert.NoError(t, err)
	updatedBookTitle.Title = updatedTitle
	updatedBookTitle.Author = book.Author

	nonExistingBook, err := core.CopyBook(updatedBookTitle)
	assert.NoError(t, err)
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
		{
			Description:      "Update Bad Request Book",
			UpdatedBook:      models.Book{ID: -1},
			ExpectedHTTPCode: http.StatusBadRequest,
			ShouldFail:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			// Convert the updated book payload to JSON
			updatedBookJSON, err := json.Marshal(tc.UpdatedBook)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Perform a PUT request to the "UpdateBook" endpoint
			method := "PUT"
			url := fmt.Sprintf("/books/%d", tc.UpdatedBook.ID)
			body := updatedBookJSON
			response, err := SendRequestV1(router, method, url, body)

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
				updatedBookFromDB, err := core.GetBook(tc.UpdatedBook.ID)
				if err != nil {
					t.Fatalf("Failed to fetch updated book from the database: %v", err)
				}

				// Verify that the book in the database matches the expected updated book data
				assert.Equal(t, tc.ExpectedDBTitle, updatedBookFromDB.Title, "Title mismatch in the database")
				assert.Equal(t, tc.ExpectedDBAuthor, updatedBookFromDB.Author, "Author mismatch in the database")
			}
		})
	}
}

func TestDeleteBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample book in the database for testing
	book, err := api.LoadSampleBook()
	assert.NoError(t, err)

	err = core.AddBook(book)
	if err != nil {
		t.Fatalf("Failed to create sample book: %v", err)
	}

	// Define test cases for deleting books
	testCases := []struct {
		Description      string
		BookID           int
		ExpectedHTTPCode int
		ShouldExist      bool
	}{
		{
			Description:      "Delete Existing Book",
			BookID:           book.ID,
			ExpectedHTTPCode: http.StatusOK,
			ShouldExist:      false,
		},
		{
			Description:      "Delete Invalid Book",
			BookID:           -1, // Use an invalid ID
			ExpectedHTTPCode: http.StatusBadRequest,
			ShouldExist:      false, // Invalid book should still not exist
		},
		{
			Description:      "Delete Non-Existent Book",
			BookID:           book.ID - 1, // Use a non-existent ID
			ExpectedHTTPCode: http.StatusNotFound,
			ShouldExist:      false, // Non-existent book still not exist
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			// Perform a DELETE request to the "DeleteBook" endpoint with the book ID
			method := "DELETE"
			url := fmt.Sprintf("/books/%d", tc.BookID)
			var body []byte = nil
			response, err := SendRequestV1(router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedHTTPCode, response.Code, "Expected status code %d, but got %d", tc.ExpectedHTTPCode, response.Code)

			// Verify the book's existence in the database
			_, err := core.GetBook(tc.BookID)
			if tc.ShouldExist {
				assert.NoError(t, err, "Expected book to exist in the database")
			} else {
				assert.Error(t, err, "Expected book to be deleted from the database")
			}
		})
	}
}
func TestSearchBookHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	listOfBooks, err := api.LoadListOfBookSamples()
	assert.NoError(t, err)
	numberOfSampleBooks := len(listOfBooks)

	for _, book := range listOfBooks {
		err = core.AddBook(&book)
		if err != nil {
			t.Fatalf("Failed to create sample book: %v", err)
		}
	}

	listOfBooks, err = core.ListBooks()
	assert.NoError(t, err)
	assert.Equal(t, numberOfSampleBooks, len(listOfBooks), "Expected %d books, but got %d", numberOfSampleBooks, len(listOfBooks))

	count, err := core.CountBooks()
	assert.NoError(t, err)
	assert.Equal(t, numberOfSampleBooks, count, "Expected %d books, but got %d", numberOfSampleBooks, len(listOfBooks))

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
				"author": listOfBooks[0].Author,
				"genre":  listOfBooks[0].GenreName,
			},
			ExpectedHTTPCode: http.StatusOK,
			ExpectedCount:    1,
		},
		{
			Description: "Search by Author Only",
			QueryParams: map[string]string{
				"author": listOfBooks[0].Author,
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

			// Create the URL for the GET request
			method := "GET"
			url := fmt.Sprintf("/books/search?%s", query)
			var body []byte = nil
			response, err := SendRequestV1(router, method, url, body)

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
			for i, expectedBook := range listOfBooks {
				if i >= len(responseBooks) {
					break // Avoid index out of range error
				}
				assert.Equal(t, expectedBook.Title, responseBooks[i].Title, "Title mismatch")
				assert.Equal(t, expectedBook.Author, responseBooks[i].Author, "Author mismatch")
				assert.Equal(t, expectedBook.Edition, responseBooks[i].Edition, "Edition mismatch")
				assert.Equal(t, expectedBook.Published.Format(time.DateOnly), responseBooks[i].Published.Format(time.DateOnly), "Published mismatch")
				assert.Equal(t, expectedBook.Description, responseBooks[i].Description, "Description mismatch")
				assert.Equal(t, expectedBook.GenreName, responseBooks[i].GenreName, "Genre mismatch")
			}
		})
	}
}
func TestCountBooksHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Load a list of sample books
	books, err := api.LoadListOfBookSamples()
	assert.NoError(t, err)

	for _, book := range books {
		err := core.AddBook(&book)
		assert.NoError(t, err)
	}

	// Define test cases for counting books
	testCases := []struct {
		Description      string
		RequestURL       string
		ExpectedHTTPCode int
		ExpectedCount    int
		ExpectedResponse map[string]int
	}{
		{
			Description:      "Count All Books",
			RequestURL:       "/books/count",
			ExpectedHTTPCode: http.StatusOK,
			ExpectedCount:    len(books),
			ExpectedResponse: map[string]int{"count": len(books)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			// Prepare the request to the "CountBooks" endpoint
			method := "GET"
			url := tc.RequestURL
			var body []byte = nil
			response, err := SendRequestV1(router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedHTTPCode, response.Code, "Expected status code %d, but got %d", tc.ExpectedHTTPCode, response.Code)

			// Parse the response JSON to get the count
			var responseCount map[string]int
			err = json.NewDecoder(response.Body).Decode(&responseCount)
			assert.NoError(t, err)

			// Check if the response contains the correct count
			assert.Equal(t, tc.ExpectedResponse, responseCount, "Unexpected count response")
		})
	}
}

*/
