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