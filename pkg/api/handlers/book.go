package handlers

import (
	"errors"
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

// SearchBooks handles the "GET /api/v1/books/search" endpoint to search for books.
func SearchBooks(c *gin.Context) {
	// Get the query parameter from the URL
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var books []models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.Where("title LIKE ? OR author LIKE ?", "%"+query+"%", "%"+query+"%").Find(&books)

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
