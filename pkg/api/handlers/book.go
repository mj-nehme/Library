package handlers

import (
	"errors"
	"fmt"
	"library/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// AddBook handles the "POST /api/v1/books" endpoint to create a new book.
func AddBook(c *gin.Context) {
	// Bind the JSON request body to a Book struct
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data. " + err.Error()})
		return
	}

	if err := validate.Struct(newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getValidationErrors(err)})
		return
	}

	// Create a new record in the database
	db := c.MustGet("db").(*gorm.DB)
	err := db.Create(&newBook).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

// GetBook handles the "GET /api/v1/books/:id" endpoint to retrieve a specific book by its ID.
func GetBook(c *gin.Context) {
	bookIDParam := c.Param("id")

	bookID, err := strconv.Atoi(bookIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID" + err.Error()})
		return
	}

	// Use the validator to check if the book ID is a valid integer
	err = validate.Var(bookID, "required,numeric")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID" + err.Error()})
		return
	}

	var book models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.First(&book, bookID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found" + result.Error.Error()})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book" + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func ListBooks(c *gin.Context) {
	var books []models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.Find(&books)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books" + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func UpdateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data. " + err.Error()})
		return
	}

	// Validate the required fields (e.g., title and author)
	err := validate.Struct(book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var existingBook models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.First(&existingBook, book.ID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found" + result.Error.Error()})
		return
	}

	existingBook.Title = book.Title
	existingBook.Author = book.Author
	existingBook.Published = book.Published
	existingBook.Edition = book.Edition
	existingBook.Description = book.Description
	existingBook.GenreName = book.GenreName

	result = db.Save(&existingBook)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book" + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, existingBook)
}

func PatchBook(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data. " + err.Error()})
		return
	}

	var existingBook models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.First(&existingBook, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found" + result.Error.Error()})
		return
	}

	result = db.Model(&existingBook).Updates(updates)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, existingBook)
}

func DeleteBook(c *gin.Context) {
	// Get the book ID from the URL parameter
	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID" + err.Error()})
		return
	}

	var existingBook models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.First(&existingBook, bookID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found" + result.Error.Error()})
		return
	}

	result = db.Delete(&existingBook)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// buildSearchQuery constructs a GORM query based on provided parameters.
func buildSearchQuery(db *gorm.DB, params map[string]string) *gorm.DB {
	query := db
	for key, value := range params {
		switch key {
		case "author":
			query = query.Where("author LIKE ?", "%"+value+"%")
		case "genre":
			query = query.Where("genre_name LIKE ?", "%"+value+"%")
		case "title":
			query = query.Where("title LIKE ?", "%"+value+"%")
		case "from":
			// Assuming "from" is the parameter for the start of the date range
			if value != "" {
				query = query.Where("published >= ?", value)
			}
		case "to":
			// Assuming "to" is the parameter for the end of the date range
			if value != "" {
				query = query.Where("published <= ?", value)
			}
		case "description":
			query = query.Where("description LIKE ?", "%"+value+"%")
		}
	}
	return query
}

// SearchBooks handles the "GET /api/v1/books/search" endpoint to search for books.
func SearchBooks(c *gin.Context) {
	// Get the query parameters from the URL
	params := map[string]string{
		"title":       c.Query("title"),
		"author":      c.Query("author"),
		"from":        c.Query("from"),
		"to":          c.Query("to"),
		"description": c.Query("description"),
		"genre":       c.Query("genre"),
	}

	// Build the search query
	db := c.MustGet("db").(*gorm.DB)
	query := buildSearchQuery(db, params)
	fmt.Println("Query: ", query)

	var books []models.Book
	result := query.Find(&books)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books" + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func CountBooks(c *gin.Context) {
	var count int64
	db := c.MustGet("db").(*gorm.DB)
	result := db.Model(&models.Book{}).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books count" + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, count)
}
