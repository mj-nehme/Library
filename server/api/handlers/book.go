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

//	@Summary		Add a new book
//	@Description	Add a new book to the library
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			newBook	body		models.Book		true	"New Book details"
//	@Success		201		{object}	models.Book		"Returns the newly created book"
//	@Failure		400		{object}	ErrorResponse	"Invalid JSON data or validation error"
//	@Failure		500		{object}	ErrorResponse	"Failed to create book"
//	@Router			/books [post]
//
// AddBook handles the "POST /books" endpoint to create a new book.
func AddBook(c *gin.Context) {
	// Bind the JSON request body to a Book struct
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON data. " + err.Error()})
		return
	}

	if err := validate.Struct(newBook); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: getValidationErrors(err)})
		return
	}

	// Create a new record in the database
	db := c.MustGet("db").(*gorm.DB)
	err := db.Create(&newBook).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create book" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

//	@Summary		Get a book by ID
//	@Description	Retrieve a book by its ID
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Book ID"
//	@Success		200	{object}	models.Book		"Returns the requested book"
//	@Failure		400	{object}	ErrorResponse	"Invalid book ID"
//	@Failure		404	{object}	ErrorResponse	"Book not found"
//	@Failure		500	{object}	ErrorResponse	"Failed to fetch book"
//	@Router			/books/{id} [get]
//
// GetBook handles the "GET /books/:id" endpoint to retrieve a specific book by its ID.
func GetBook(c *gin.Context) {
	bookIDParam := c.Param("id")

	bookID, err := strconv.Atoi(bookIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid book ID" + err.Error()})
		return
	}

	// Use the validator to check if the book ID is a valid integer
	err = validate.Var(bookID, "required,numeric")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid book ID" + err.Error()})
		return
	}

	var book models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.First(&book, bookID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Book not found" + result.Error.Error()})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch book" + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// @Summary		List books
// @Description	Retrieve a list of all books
// @Tags			books
// @Accept			json
// @Produce		json
// @Success		200	{array}		models.Book		"Returns the list of books"
// @Failure		500	{object}	ErrorResponse	"Failed to retrieve books"
// @Router			/books [get]
func ListBooks(c *gin.Context) {
	var books []models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.Find(&books)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve books" + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// @Summary		Update a book
// @Description	Update a book's details
// @Tags			books
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"Book ID"
// @Param			book	body		models.Book		true	"Updated Book details"
// @Success		200		{object}	models.Book		"Returns the updated book"
// @Failure		400		{object}	ErrorResponse	"Invalid JSON data or validation error"
// @Failure		404		{object}	ErrorResponse	"Book not found"
// @Failure		500		{object}	ErrorResponse	"Failed to update book"
// @Router			/books/{id} [put]
func UpdateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON data. " + err.Error()})
		return
	}

	// Validate the required fields (e.g., title and author)
	err := validate.Struct(book)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var existingBook models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.First(&existingBook, book.ID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Book not found" + result.Error.Error()})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update book" + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, existingBook)
}

// @Summary		Patch a book
// @Description	Partially update a book's details
// @Tags			books
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"Book ID"
// @Param			book	body		models.Book		true	"Updated Book details"
// @Success		200		{object}	models.Book		"Returns the updated book"
// @Failure		400		{object}	ErrorResponse	"Invalid JSON data or validation error"
// @Failure		404		{object}	ErrorResponse	"Book not found"
// @Failure		500		{object}	ErrorResponse	"Failed to update book"
// @Router			/books/{id} [patch]
func PatchBook(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON data. " + err.Error()})
		return
	}

	var existingBook models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.First(&existingBook, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Book not found" + result.Error.Error()})
		return
	}

	result = db.Model(&existingBook).Updates(updates)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, existingBook)
}

// @Summary		Delete a book
// @Description	Delete a book by its ID
// @Tags			books
// @Accept			json
// @Produce		json
// @Param			id	path		int				true	"Book ID"
// @Success		200	{object}	MessageResponse	"Returns a success message"
// @Failure		400	{object}	ErrorResponse	"Invalid book ID"
// @Failure		404	{object}	ErrorResponse	"Book not found"
// @Failure		500	{object}	ErrorResponse	"Failed to delete book"
// @Router			/books/{id} [delete]
func DeleteBook(c *gin.Context) {
	// Get the book ID from the URL parameter
	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid book ID" + err.Error()})
		return
	}

	var existingBook models.Book
	db := c.MustGet("db").(*gorm.DB)
	result := db.First(&existingBook, bookID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Book not found" + result.Error.Error()})
		return
	}

	result = db.Delete(&existingBook)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Book deleted successfully"})
}

//	@Summary		Search for books
//	@Description	Search for books based on various criteria
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			title		query		string			false	"Title of the book"
//	@Param			author		query		string			false	"Author of the book"
//	@Param			from		query		string			false	"Published date range start (YYYY-MM-DD)"
//	@Param			to			query		string			false	"Published date range end (YYYY-MM-DD)"
//	@Param			description	query		string			false	"Description of the book"
//	@Param			genre		query		string			false	"Genre of the book"
//	@Success		200			{array}		models.Book		"Returns the list of matching books"
//	@Failure		500			{object}	ErrorResponse	"Failed to fetch books"
//	@Router			/books/search [get]
//
// SearchBooks handles the "GET /books/search" endpoint to search for books.
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

	var books []models.Book
	result := query.Find(&books)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch books" + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

// @Summary		Count books
// @Description	Get the total count of books
// @Tags			books
// @Accept			json
// @Produce		json
// @Success		200	{integer}	int64			"Returns the total count of books"
// @Failure		500	{object}	ErrorResponse	"Failed to retrieve books count"
// @Router			/books/count [get]
func CountBooks(c *gin.Context) {
	var count int64
	db := c.MustGet("db").(*gorm.DB)
	result := db.Model(&models.Book{}).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve books count" + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, count)
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
