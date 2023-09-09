package api_test

/*
import (
	"encoding/json"
	"fmt"
	"library/core"
	"library/models"
	"library/tests"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGenre(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample genre in the database for testing
	genre, err := core.LoadSampleGenre()
	assert.NoError(t, err)

	// Define test cases
	testCases := []struct {
		Name           string
		Genre          *models.Genre
		ExpectedStatus int
		ExpectedError  string
	}{
		{
			Name: "Add Genre with Empty Name",
			Genre: &models.Genre{
				Name: "",
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  "Invalid genre name",
		},
		{
			Name:           "Add Genre Successfully",
			Genre:          genre,
			ExpectedStatus: http.StatusCreated,
			ExpectedError:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Marshal the NewGenre struct to JSON
			requestBody, err := json.Marshal(tc.Genre)
			assert.NoError(t, err)

			// Create a request to AddGenre
			method := "POST"
			url := "/genres"
			var body []byte = requestBody
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Expected status code %d, but got %d", tc.ExpectedStatus, response.Code)
		})
	}
}
func TestGetGenre(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample genre in the database for testing
	genre, err := core.LoadSampleGenre()
	assert.NoError(t, err)

	genreID, err := core.AddGenre(genre.Name)
	assert.NoError(t, err)

	// Define test cases
	testCases := []struct {
		Name           string
		GenreID        int
		GenreName      string
		ExpectedStatus int
	}{
		{
			Name:           "Get Genre Successfully",
			GenreID:        genreID,
			GenreName:      genre.Name,
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "Genre Not Found",
			GenreID:        0,
			GenreName:      "NonExistentGenre",
			ExpectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a request to GetGenre
			method := "GET"
			url := fmt.Sprintf("/genres/%d", tc.GenreID)
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Expected status code %d, but got %d", tc.ExpectedStatus, response.Code)
		})
	}
}

func TestListGenres(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample genre in the database for testing
	genres, err := core.LoadListOfGenreSamples()
	assert.NoError(t, err)

	// Define test cases
	testCases := []struct {
		Name               string
		prepareDB          func() error // Function to prepare the database with test data
		ExpectedStatus     int
		ExpectedGenreCount int
	}{
		{
			Name:               "No Genres Found",
			prepareDB:          func() error { return nil },
			ExpectedStatus:     http.StatusOK,
			ExpectedGenreCount: 0,
		},
		{
			Name: "List Genres Successfully",
			prepareDB: func() error {
				for index, genre := range genres {
					genreID, err := core.AddGenre(genre.Name)
					assert.NoError(t, err)
					genres[index].ID = genreID
				}
				return nil
			},
			ExpectedStatus:     http.StatusOK,
			ExpectedGenreCount: len(genres),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.prepareDB()
			assert.NoError(t, err)
			// Create a request to ListGenres
			method := "GET"
			url := "/genres"
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Expected status code %d, but got %d", tc.ExpectedStatus, response.Code)

			// Unmarshal the response body to get the list of genres
			var genres []models.Genre
			err = json.Unmarshal(response.Body.Bytes(), &genres)
			assert.NoError(t, err)

			// Check the number of genres in the response
			assert.Len(t, genres, tc.ExpectedGenreCount)
		})
	}
}

func TestUpdateGenre(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample genre in the database for testing
	genre, err := core.LoadSampleGenre()
	assert.NoError(t, err)

	genreID, err := core.AddGenre(genre.Name)
	assert.NoError(t, err)

	updatedName := "Science Fiction"

	// Define test cases
	testCases := []struct {
		Name             string
		GenreID          int
		UpdatedGenreName string
		ExpectedStatus   int
	}{
		{
			Name:             "Update Genre Successfully",
			GenreID:          genreID,
			UpdatedGenreName: updatedName,
			ExpectedStatus:   http.StatusOK,
		},
		{
			Name:             "Update Genre With Empty Name",
			GenreID:          genreID,
			UpdatedGenreName: "",
			ExpectedStatus:   http.StatusBadRequest,
		},
		{
			Name:             "Update Invalid Genre",
			GenreID:          -1,
			UpdatedGenreName: updatedName,
			ExpectedStatus:   http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			method := "PUT"
			url := fmt.Sprintf("/genres/%d", tc.GenreID)
			body := []byte(fmt.Sprintf(`{"name": "%s"}`, tc.UpdatedGenreName))
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Expected status code %d, but got %d", tc.ExpectedStatus, response.Code)
		})
	}
}

func TestDeleteGenre(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample genre in the database for testing
	genre, err := core.LoadSampleGenre()
	assert.NoError(t, err)

	genreID, err := core.AddGenre(genre.Name)
	assert.NoError(t, err)

	// Define test cases
	testCases := []struct {
		Name           string
		GenreID        int
		ExpectedStatus int
	}{
		{
			Name:           "Delete Genre Successfully",
			GenreID:        genreID,
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Name:           "Attempt to Delete Non-Existing Genre",
			GenreID:        genreID + 1,
			ExpectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			method := "DELETE"
			url := fmt.Sprintf("/genres/%d", tc.GenreID)
			body := []byte(nil)
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Expected status code %d, but got %d", tc.ExpectedStatus, response.Code)

			if tc.ExpectedStatus == http.StatusNoContent {
				// If the genre was deleted successfully, try fetching it to ensure it is not found
				method = "GET"
				response = sendRequestV1(t, router, method, url, body)
				assert.Equal(t, http.StatusNotFound, response.Code, "Expected status code 404, but got %d", response.Code)
			}
		})
	}
}

func TestCountGenresHandler(t *testing.T) {
	router, db, ctx := tests.SetupMockServer()
	defer tests.TearDownMockServer(db, ctx)

	// Create a sample genre in the database for testing
	genres, err := core.LoadListOfGenreSamples()
	assert.NoError(t, err)

	// Define test cases
	testCases := []struct {
		Name               string
		prepareDB          func() error // Function to prepare the database with test data
		ExpectedStatus     int
		ExpectedGenreCount int
	}{
		{
			Name:               "No Genres Found",
			prepareDB:          func() error { return nil },
			ExpectedStatus:     http.StatusOK,
			ExpectedGenreCount: 0,
		},
		{
			Name: "Count Genres Successfully",
			prepareDB: func() error {
				for index, genre := range genres {
					genreID, err := core.AddGenre(genre.Name)
					assert.NoError(t, err)
					genres[index].ID = genreID
				}
				return nil
			},
			ExpectedStatus:     http.StatusOK,
			ExpectedGenreCount: len(genres),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.prepareDB()
			assert.NoError(t, err)

			// Perform a GET request to the CountGenres endpoint
			method := "GET"
			url := "/genres/count"
			var body []byte = nil
			response := sendRequestV1(t, router, method, url, body)

			// Check the response status code
			assert.Equal(t, tc.ExpectedStatus, response.Code, "Expected status code %d, but got %d", tc.ExpectedStatus, response.Code)

			// Parse the response JSON to get the count
			var responseCount map[string]int
			err = json.NewDecoder(response.Body).Decode(&responseCount)
			assert.NoError(t, err)

			// Check if the response contains the correct count
			count, found := responseCount["count"]
			if !found {
				t.Fatal("Response JSON does not contain 'count'")
			}
			assert.Equal(t, tc.ExpectedGenreCount, count, "Unexpected count response")
		})
	}
}
*/
