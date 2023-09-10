package api

import (
	"encoding/json"
	"fmt"
	"library/models"
	"library/tools"
	"os"
)

const (
	BookSamplePath              = "json/book_sample.json"
	ListOfBookSamplesPath       = "json/book_list_sample.json"
	CollectionSamplePath        = "json/collection_sample.json"
	ListOfCollectionSamplesPath = "json/collection_list_sample.json"
	GenreSamplePath             = "json/genre_sample.json"
	ListOfGenreSamplesPath      = "json/genre_list_sample.json"
)

func LoadSampleBook() (models.Book, error) {
	var sampleBook models.Book
	rootDir, err := tools.SearchRootDirectory()
	if err != nil {
		return models.Book{}, fmt.Errorf("unable to rerieve root directory. %v", err)
	}

	file, err := os.Open(rootDir + "/" + BookSamplePath)
	if err != nil {
		return models.Book{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a JSON decoder from the file
	decoder := json.NewDecoder(file)

	// Decode the JSON data into a Book struct
	err = decoder.Decode(&sampleBook)
	if err != nil {
		return sampleBook, err
	}

	return sampleBook, nil
}

func LoadListOfBookSamples() ([]models.Book, error) {
	var sampleBooks []models.Book
	rootDir, err := tools.SearchRootDirectory()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve root directory. %v", err)
	}

	file, err := os.Open(rootDir + "/" + ListOfBookSamplesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a JSON decoder from the file
	decoder := json.NewDecoder(file)

	// Decode the JSON data into a slice of Book structs
	err = decoder.Decode(&sampleBooks)
	if err != nil {
		return nil, err
	}

	return sampleBooks, nil
}

/*
func LoadSampleCollection() (*models.Collection, error) {
	// Load sample collection data from a JSON file
	rootDir, err := tools.SearchRootDirectory()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve root directory: %v", err)
	}

	file, err := os.Open(rootDir + "/" + CollectionSamplePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a JSON decoder from the file
	decoder := json.NewDecoder(file)

	// Read the JSON data and extract individual fields
	var collection models.Collection
	if err := decoder.Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to decode collection data: %w", err)
	}

	return &collection, nil
}

func LoadSampleGenre() (*models.Genre, error) {
	// Load sample genre data from a JSON file
	rootDir, err := tools.SearchRootDirectory()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve root directory: %v", err)
	}

	file, err := os.Open(rootDir + "/" + GenreSamplePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a JSON decoder from the file
	decoder := json.NewDecoder(file)

	// Read the JSON data and extract individual fields
	var genre models.Genre
	if err := decoder.Decode(&genre); err != nil {
		return nil, fmt.Errorf("failed to decode genre data: %w", err)
	}

	return &genre, nil
}

func LoadListOfGenreSamples() ([]models.Genre, error) {
	// Load a list of sample genre data from a JSON file
	rootDir, err := tools.SearchRootDirectory()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve root directory: %v", err)
	}

	file, err := os.Open(rootDir + "/" + ListOfGenresSample)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a JSON decoder from the file
	decoder := json.NewDecoder(file)

	// Read the opening bracket of the JSON array
	_, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("error reading opening bracket: %w", err)
	}

	var listOfGenres []models.Genre
	for decoder.More() {
		// Read the JSON data and extract individual fields
		var genre models.Genre
		if err := decoder.Decode(&genre); err != nil {
			return nil, fmt.Errorf("failed to decode genre data: %w", err)
		}
		listOfGenres = append(listOfGenres, genre)
	}

	// Read the closing bracket of the JSON array
	_, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("error reading closing bracket: %w", err)
	}

	return listOfGenres, nil
}

func LoadListOfCollectionSamples() ([]models.Collection, error) {
	// Load a list of sample collection data from a JSON file
	rootDir, err := tools.SearchRootDirectory()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve root directory: %v", err)
	}

	file, err := os.Open(rootDir + "/" + ListOfCollectionsSample)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a JSON decoder from the file
	decoder := json.NewDecoder(file)

	// Read the opening bracket of the JSON array
	_, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("error reading opening bracket: %w", err)
	}

	var listOfCollections []models.Collection
	for decoder.More() {
		// Read the JSON data and extract individual fields
		var collection models.Collection
		if err := decoder.Decode(&collection); err != nil {
			return nil, fmt.Errorf("failed to decode collection data: %w", err)
		}
		listOfCollections = append(listOfCollections, collection)
	}

	// Read the closing bracket of the JSON array
	_, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("error reading closing bracket: %w", err)
	}

	return listOfCollections, nil
}

func UnmarshalListOfBooks(file *os.File) ([]models.Book, error) {
	// Create a JSON decoder from the file
	decoder := json.NewDecoder(file)

	// Read the opening bracket of the JSON array
	_, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("error reading opening bracket: %w", err)

	}

	tempListOfBooks, err := BooksDecoder(decoder)
	if err != nil {
		return nil, fmt.Errorf("failed to read books from file: %w", err)
	}

	// Read the closing bracket of the JSON array
	_, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("error reading closing bracket: %w", err)

	}
	var listOfBooks []models.Book
	for _, tempBook := range tempListOfBooks {
		book, err := NewBook(tempBook.Title, tempBook.Author, tempBook.Published, tempBook.Edition, tempBook.Description, tempBook.GenreName)
		if err != nil {
			return nil, fmt.Errorf("failed to create new Book from Sample: %w", err)
		}
		listOfBooks = append(listOfBooks, *book)
	}

	return listOfBooks, nil
} */
/*
func unmarshalBook(m map[string]interface{}) (*models.Book, error) {
	// Create a Book struct to store the extracted fields
	var book models.Book

	// Extract individual fields from the map and convert types as needed
	if v, ok := m["id"].(float64); ok {
		book.ID = uint(v)
	} else {
		// Set default value for ID
		book.ID = 0
	}
	if v, ok := m["title"].(string); ok {
		book.Title = v
	}
	if v, ok := m["author"].(string); ok {
		book.Author = v
	}
	if v, ok := m["published"].(string); ok {
		publishedTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse published date: %w", err)
		}
		book.Published = publishedTime
	}
	if v, ok := m["edition"].(float64); ok {
		book.Edition = int(v)
	}
	if v, ok := m["description"].(string); ok {
		book.Description = v
	}
	if v, ok := m["genre_name"].(string); ok {
		book.GenreName = v
	}

	if book.Published.IsZero() {
		book.Published = time.Time{}
	}

	return &book, nil
}

func BooksDecoder(decoder *json.Decoder) ([]models.Book, error) {
	var listOfBooks []models.Book

	// Read the JSON data and extract individual book entries
	for decoder.More() {
		// Read the JSON data and extract individual fields
		var m map[string]interface{}
		if err := decoder.Decode(&m); err != nil {
			return nil, err
		}

		book, err := unmarshalBook(m)
		if err != nil {
			return nil, err
		}
		listOfBooks = append(listOfBooks, *book)
	}

	return listOfBooks, nil
}
*/
