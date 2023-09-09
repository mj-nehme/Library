package api_test

/*
import (
	"encoding/json"
	"fmt"
	"io"
	"library/core"
	"library/models"
	"library/tests"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCollectionHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	collectionName := "New Collection"

	testCases := []struct {
		Name           string
		CollectionName string
		ExpectedStatus int
	}{
		{
			Name:           "ValidCollection",
			CollectionName: collectionName,
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "DuplicateCollection",
			CollectionName: collectionName, // Duplicate name
			ExpectedStatus: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a new collection payload
			collection := &models.Collection{
				Name: tc.CollectionName,
			}

			payload, err := json.Marshal(collection)
			assert.NoError(t, err, "Failed to marshal collection payload")

			// Create a request to the AddCollection endpoint with the new collection payload
			method := "POST"
			url := "/collections"
			var body []byte = payload
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Unexpected status code")

			if tc.ExpectedStatus == http.StatusCreated {
				// Read the response body
				responseBody, err := io.ReadAll(response.Body)
				assert.Nil(t, err, "Failed to read response body")

				// Ensure that the response contains the success message
				assert.Contains(t, string(responseBody), tc.CollectionName, "Expected success message in the response")

				// Verify that the collection has been created in the database
				collections, err := core.ListCollections()
				assert.Nil(t, err, "Failed to list collections")
				assert.Len(t, collections, 1, "Expected 1 collection after creation")
				assert.Equal(t, tc.CollectionName, collections[0].Name, "Collection name mismatch")
			}
		})
	}
}
func TestGetCollectionHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample collection for testing
	collectionSample, err := core.LoadSampleCollection()
	assert.NoError(t, err)
	collection, err := core.AddCollection(collectionSample.Name)
	assert.NoError(t, err)

	testCases := []struct {
		Name           string
		CollectionID   int
		ExpectedStatus int
		ExpectedName   string
	}{
		{
			Name:           "ValidCollection",
			CollectionID:   collection.ID,
			ExpectedStatus: http.StatusOK,
			ExpectedName:   collection.Name,
		},
		{
			Name:           "NonExistentCollection",
			CollectionID:   collection.ID + 1, // Non-existent ID
			ExpectedStatus: http.StatusNotFound,
			ExpectedName:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a request to the GetCollection endpoint with the collection ID
			method := "GET"
			url := fmt.Sprintf("/collections/%d", tc.CollectionID)
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Unexpected status code")

			if tc.ExpectedStatus == http.StatusOK {
				// Read the response body
				responseBody, err := io.ReadAll(response.Body)
				assert.Nil(t, err, "Failed to read response body")

				// Unmarshal the response JSON into a Collection object
				var gotCollection models.Collection
				err = json.Unmarshal(responseBody, &gotCollection)
				assert.Nil(t, err, "Failed to unmarshal response JSON")

				// Ensure that the response matches the expected collection
				assert.Equal(t, tc.ExpectedName, gotCollection.Name, "Collection name mismatch")
			}
		})
	}
}

func TestListCollectionsHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	collectionSamples, err := core.LoadListOfCollectionSamples()
	assert.NoError(t, err)
	collectionSampleNames := []string{}
	for _, cs := range collectionSamples {
		collectionSampleNames = append(collectionSampleNames, cs.Name)
	}

	testCases := []struct {
		Name           string
		prepareDB      func() error // Function to prepare the database with test data
		ExpectedStatus int
		ExpectedCount  int
	}{
		{
			Name:           "EmptyCollections",
			prepareDB:      func() error { return nil },
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  0, // No collections exist
		},
		{
			Name:           "ValidCollections",
			prepareDB:      AddCollections,
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  len(collectionSamples), // Total count including the sample collection
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.prepareDB()
			// Create a request to the ListCollections endpoint
			method := "GET"
			url := "/collections"
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Unexpected status code")

			if tc.ExpectedStatus == http.StatusOK {
				// Read the response body
				responseBody, err := io.ReadAll(response.Body)
				assert.Nil(t, err, "Failed to read response body")

				// Unmarshal the response JSON into a slice of Collection objects
				var collections []models.Collection
				err = json.Unmarshal(responseBody, &collections)
				assert.Nil(t, err, "Failed to unmarshal response JSON")

				// Ensure that the response contains the expected number of collections
				assert.Len(t, collections, tc.ExpectedCount, "Unexpected number of collections")

				if tc.ExpectedCount > 0 {

					collectionNames := []string{}
					for _, collection := range collections {
						collectionNames = append(collectionNames, collection.Name)
					}
					assert.ElementsMatch(t, collectionNames, collectionSampleNames)
				}
			}
		})
	}
}

func TestUpdateCollectionHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Prepare a sample collection for testing
	collectionSample, err := core.LoadSampleCollection()
	assert.NoError(t, err)
	collection, err := core.AddCollection(collectionSample.Name)
	assert.NoError(t, err)

	testCases := []struct {
		Name                 string
		UpdatedCollection    models.Collection
		ExpectedStatus       int
		ExpectedSuccessMsg   string
		ExpectedCollectionID int
	}{
		{
			Name: "ValidUpdate",
			UpdatedCollection: models.Collection{
				ID:   collection.ID,
				Name: "Updated Fiction Books",
			},
			ExpectedStatus:       http.StatusOK,
			ExpectedSuccessMsg:   "Updated Fiction Books",
			ExpectedCollectionID: collection.ID,
		},
		{
			Name: "InvalidCollectionID",
			UpdatedCollection: models.Collection{
				ID:   99999, // An invalid ID that does not exist
				Name: "Updated Collection Name",
			},
			ExpectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Convert the updated collection payload to JSON
			updatedCollectionJSON, err := json.Marshal(tc.UpdatedCollection)
			assert.Nil(t, err, "Failed to marshal updated collection payload")

			// Create a request to the UpdateCollection endpoint with the updated collection payload
			method := "PUT"
			url := fmt.Sprintf("/collections/%d", tc.UpdatedCollection.ID)
			var body []byte = updatedCollectionJSON
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Unexpected status code")

			if tc.ExpectedStatus == http.StatusOK {
				// Read the response body
				responseBody, err := io.ReadAll(response.Body)
				assert.Nil(t, err, "Failed to read response body")

				// Ensure that the response contains the success message
				assert.Contains(t, string(responseBody), tc.ExpectedSuccessMsg, "Expected success message in the response")

				// Verify that the collection has been updated in the database
				gotCollection, err := core.GetCollection(tc.ExpectedCollectionID)
				assert.Nil(t, err, "Failed to get updated collection")
				assert.Equal(t, gotCollection.Name, tc.ExpectedSuccessMsg, "Collection name mismatch after update")
			}
		})
	}
}

func TestDeleteCollectionHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Prepare a sample collection for testing
	collectionSample, err := core.LoadSampleCollection()
	assert.NoError(t, err)
	collection, err := core.AddCollection(collectionSample.Name)
	assert.NoError(t, err)

	testCases := []struct {
		Name                string
		CollectionID        int
		ExpectedStatus      int
		ExpectedCollections int
	}{
		{
			Name:                "InvalidCollectionID",
			CollectionID:        collection.ID + 1, // An invalid ID that does not exist
			ExpectedStatus:      http.StatusNotFound,
			ExpectedCollections: 1, // Expect the collection to still exist in this case
		},
		{
			Name:                "ValidDelete",
			CollectionID:        collection.ID,
			ExpectedStatus:      http.StatusNoContent,
			ExpectedCollections: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a request to the DeleteCollection endpoint with the collection's ID
			method := "DELETE"
			url := fmt.Sprintf("/collections/%d", tc.CollectionID)
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Unexpected status code")

			// Verify the number of collections in the database
			collections, err := core.ListCollections()
			assert.Nil(t, err, "Failed to list collections")
			assert.Len(t, collections, tc.ExpectedCollections, "Unexpected number of collections after deletion")
		})
	}
}

func TestCountCollectionsHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create some collections and books to add to the database
	sampleCollections, err := core.LoadListOfCollectionSamples()
	assert.NoError(t, err)

	collections := []models.Collection{}
	// Add the collections to the database
	for index := range sampleCollections {
		collection, err := core.AddCollection(sampleCollections[index].Name)
		assert.NoError(t, err)
		collections = append(collections, *collection)
	}

	testCases := []struct {
		Name            string
		ExpectedCount   int
		ExpectedStatus  int
		PostDeletion    bool
		CollectionToDel int
	}{
		{
			Name:            "InitialCount",
			ExpectedCount:   len(collections),
			ExpectedStatus:  http.StatusOK,
			PostDeletion:    false,
			CollectionToDel: -1,
		},
		{
			Name:            "AfterCollectionDeletion",
			ExpectedCount:   len(collections) - 1,
			ExpectedStatus:  http.StatusOK,
			PostDeletion:    true,
			CollectionToDel: 0, // Delete the first collection
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.PostDeletion {
				// Delete a collection if specified
				assert.True(t, tc.CollectionToDel >= 0 && tc.CollectionToDel < len(collections), "Invalid CollectionToDel index")
				err := core.DeleteCollectionByID(collections[tc.CollectionToDel].ID)
				assert.NoError(t, err)
			}

			// Perform a GET request to the CountCollections endpoint
			method := "GET"
			url := "/collections/count"
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code)

			// Parse the response body
			var responseBody map[string]int
			err := json.NewDecoder(response.Body).Decode(&responseBody)
			assert.NoError(t, err)

			// Check if the response count matches the expected count
			assert.Equal(t, tc.ExpectedCount, responseBody["count"], "Expected count does not match")
		})
	}
}

func TestAddBookToCollectionHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Prepare a sample collection
	collectionSample, err := core.LoadSampleCollection()
	assert.NoError(t, err)
	collection, err := core.AddCollection(collectionSample.Name)
	assert.NoError(t, err)
	book, err := core.LoadSampleBook()
	assert.NoError(t, err)
	err = core.AddBook(book)
	assert.NoError(t, err)

	testCases := []struct {
		Name                 string
		CollectionID         int
		BookID               int
		ExpectedStatus       int
		ExpectedMessage      string
		CollectionBooksCount int
	}{
		{
			Name:                 "AddBookToCollectionSuccess",
			CollectionID:         collection.ID,
			BookID:               book.ID, // Pick any book ID from the loaded list of books
			ExpectedStatus:       http.StatusCreated,
			ExpectedMessage:      "Book added to the collection successfully",
			CollectionBooksCount: 1, // One book added to the collection
		},
		{
			Name:                 "AddExistingBookToCollection",
			CollectionID:         collection.ID,
			BookID:               book.ID, // Attempt to add the same book again
			ExpectedStatus:       http.StatusConflict,
			ExpectedMessage:      "Book is already in the collection",
			CollectionBooksCount: 1, // Collection should still have only 1 book
		},
		{
			Name:                 "AddNonexistentBookToCollection",
			CollectionID:         collection.ID,
			BookID:               book.ID + 1, // Nonexistent book ID
			ExpectedStatus:       http.StatusNotFound,
			ExpectedMessage:      "Book not found",
			CollectionBooksCount: 1, // Collection should still have only 1 book
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a new CollectionBook struct
			bookToAdd := models.CollectionBook{
				CollectionID: tc.CollectionID,
				BookID:       tc.BookID,
			}

			// Marshal the CollectionBook struct to JSON
			requestBody, err := json.Marshal(bookToAdd)
			assert.NoError(t, err)

			// Create a request to AddToCollection
			method := "POST"
			url := fmt.Sprintf("/collections/%d/books/add", tc.CollectionID)
			var body []byte = requestBody
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code)

			// Check the collection's book count in the database
			collectionBooksCount, err := core.CountBooksPerCollection(tc.CollectionID)
			assert.NoError(t, err)
			assert.Equal(t, tc.CollectionBooksCount, collectionBooksCount)
		})
	}
}

func TestListBooksInCollectionHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	listOfBooks, err := core.LoadListOfBookSamples()
	assert.NoError(t, err)
	for index, book := range listOfBooks {
		err := core.AddBook(&book)
		assert.NoError(t, err)
		listOfBooks[index].ID = book.ID
	}
	assert.LessOrEqual(t, 10, len(listOfBooks))

	listOfCollections, err := core.LoadListOfCollectionSamples()
	assert.NoError(t, err)
	for index, collection := range listOfCollections {
		newCollection, err := core.AddCollection(collection.Name)
		assert.NoError(t, err)
		listOfCollections[index].ID = newCollection.ID
	}
	assert.LessOrEqual(t, 3, len(listOfCollections))

	// Associate some collections to book number 0
	err = core.AddBookToCollection(listOfCollections[0].ID, listOfBooks[0].ID)
	assert.NoError(t, err)
	err = core.AddBookToCollection(listOfCollections[0].ID, listOfBooks[4].ID)
	assert.NoError(t, err)
	err = core.AddBookToCollection(listOfCollections[0].ID, listOfBooks[9].ID)
	assert.NoError(t, err)

	// Define a set of test cases with different scenarios
	testCases := []struct {
		Name           string // Test case name (for better error reporting)
		collectionID   int
		bookIDs        []int
		ExpectedStatus int
	}{
		{
			Name:           "No Books in Collection",
			collectionID:   listOfCollections[2].ID,
			bookIDs:        []int{},
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name:           "List Books in Collection",
			collectionID:   listOfCollections[0].ID,
			bookIDs:        []int{listOfBooks[0].ID, listOfBooks[4].ID, listOfBooks[9].ID},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "Invalid Collection",
			collectionID:   -1,
			bookIDs:        []int{},
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a request to ListBooksInCollection for the collection
			method := "GET"
			url := fmt.Sprintf("/collections/%d/books", tc.collectionID)
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Expected status code %d, but got %d", tc.ExpectedStatus, response.Code)

			if response.Code == http.StatusOK {
				// Unmarshal the response body
				var retrievedBooks []models.Book
				err := json.Unmarshal(response.Body.Bytes(), &retrievedBooks)
				if err != nil {
					t.Fatalf("Failed to unmarshal response JSON: %v", err)
				}
				assert.Equal(t, len(tc.bookIDs), len(retrievedBooks), "Test case %q failed. Incorrect number of books in the collection.", tc.Name)

				// Check if the retrieved books match the expected books
				bookIDs := []int{}
				for _, book := range retrievedBooks {
					bookIDs = append(bookIDs, book.ID)
				}
				assert.ElementsMatch(t, tc.bookIDs, bookIDs, "Retrieved collections do not match the expected collections")
			}
		})
	}
}

func TestListCollectionsOfBookAPIHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	listOfBooks, err := core.LoadListOfBookSamples()
	assert.NoError(t, err)
	for index, book := range listOfBooks {
		err := core.AddBook(&book)
		assert.NoError(t, err)
		listOfBooks[index].ID = book.ID
	}
	assert.LessOrEqual(t, 10, len(listOfBooks))

	listOfCollections, err := core.LoadListOfCollectionSamples()
	assert.NoError(t, err)
	for index, collection := range listOfCollections {
		newCollection, err := core.AddCollection(collection.Name)
		assert.NoError(t, err)
		listOfCollections[index].ID = newCollection.ID
	}
	assert.LessOrEqual(t, 3, len(listOfCollections))

	// Associate some collections to book number 0
	err = core.AddBookToCollection(listOfCollections[0].ID, listOfBooks[0].ID)
	assert.NoError(t, err)
	err = core.AddBookToCollection(listOfCollections[1].ID, listOfBooks[0].ID)
	assert.NoError(t, err)

	// Define a set of test cases with different scenarios
	testCases := []struct {
		Name                  string // Test case name (for better error reporting)
		bookID                int    // Book ID to list its collections
		expectedCollectionIDs []int  // Expected IDs of collections associated with the book
		ExpectedCode          int    // Expected HTTP response code
	}{
		{
			Name:                  "List Collections of Book",
			bookID:                listOfBooks[0].ID,
			expectedCollectionIDs: []int{listOfCollections[0].ID, listOfCollections[1].ID},
			ExpectedCode:          http.StatusOK,
		},
		{
			Name:                  "List Collections of Non-associated Book",
			bookID:                listOfBooks[9].ID,
			expectedCollectionIDs: nil,
			ExpectedCode:          http.StatusNotFound,
		},
		{
			Name:                  "List Collections of Non-existing Book",
			bookID:                -1,
			expectedCollectionIDs: nil,
			ExpectedCode:          http.StatusBadRequest,
		},
	}

	// Loop over the test cases and run them
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Perform a GET request to the API endpoint
			method := "GET"
			url := fmt.Sprintf("/books/%d/collections", tc.bookID)
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedCode, response.Code)

			if response.Code == http.StatusOK {
				// Unmarshal the response body
				var retrievedCollections []models.Collection
				err := json.Unmarshal(response.Body.Bytes(), &retrievedCollections)
				if err != nil {
					t.Fatalf("Failed to unmarshal response JSON: %v", err)
				}
				assert.Equal(t, len(tc.expectedCollectionIDs), len(retrievedCollections), "Test case %q failed. Incorrect number of books in the collection.", tc.Name)

				// Check if the retrieved collections match the expected collections
				collectionIDs := []int{}
				for _, collection := range retrievedCollections {
					collectionIDs = append(collectionIDs, collection.ID)
				}
				assert.ElementsMatch(t, tc.expectedCollectionIDs, collectionIDs, "Retrieved collections do not match the expected collections")
			}
		})
	}
}

func AddCollections() error {
	// Create a sample collection for testing
	collectionSamples, err := core.LoadListOfCollectionSamples()
	if err != nil {
		return err
	}
	for _, sample := range collectionSamples {
		_, err = core.AddCollection(sample.Name)
		if err != nil {
			return err
		}
	}
	return err
}
*/
