package handlers

/*
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddCollection handles the "POST /api/v1/collections" endpoint to create a new collection.
func AddCollection(c *gin.Context) {
	// Bind the JSON request body to a Collection struct
	var collection = core.NewCollectionTemplate()
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data. " + err.Error()})
		return
	}

	// Validate the required fields (e.g., title and author)
	err := validators.AddCollection(collection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Perform collection creation in the database
	newCollection, err := core.AddCollection(collection.Name)
	if err != nil {
		if _, ok := err.(*core.AlreadyExists); ok {
			c.JSON(http.StatusConflict, gin.H{"error": "collection already exists" + err.Error()})
			return
		} else if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection" + err.Error()})
			return
		}
	}

	// Return the newly created collection as the response
	c.JSON(http.StatusCreated, newCollection)
}

// GetCollection retrieves a specific collection by its ID.
func GetCollection(c *gin.Context) {
	collectionIDParam := c.Param("id")
	collectionID, err := strconv.Atoi(collectionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID" + err.Error()})
		return
	}

	collection, err := core.GetCollection(collectionID)
	if err != nil {
		if _, ok := err.(*core.DoesNotExist); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "collection not found" + err.Error()})
			return
		} else if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collection" + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, collection)
}

// ListCollections retrieves a list of all collections.
func ListCollections(c *gin.Context) {
	collections, err := core.ListCollections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve collections" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, collections)
}

// UpdateCollection handles the "PUT /api/v1/collections/:id" endpoint to update an existing collection.
func UpdateCollection(c *gin.Context) {
	// Get the collection ID from the URL parameter
	collectionIDParam := c.Param("id")
	collectionID, err := strconv.Atoi(collectionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID" + err.Error()})
		return
	}

	// Check if the collection with the given ID exists in the database
	existingCollection, err := core.GetCollection(collectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found" + err.Error()})
		return
	}

	// Bind the request JSON body to a Collection object
	var updatedCollection models.Collection
	if err := c.ShouldBindJSON(&updatedCollection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body" + err.Error()})
		return
	}

	// Update the collection information
	existingCollection.Name = updatedCollection.Name
	// Update other fields as needed

	// Save the updated collection in the database
	affectedRows, err := core.UpdateCollection(existingCollection)
	if err != nil {
		if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collection" + err.Error()})
		return
	}
	if affectedRows == 0 {
		c.JSON(http.StatusNotModified, gin.H{"error": "No affected rows" + err.Error()})
		return
	}

	// Return the updated collection as the response
	c.JSON(http.StatusOK, existingCollection)
}

// DeleteCollection deletes a collection from the system.
func DeleteCollection(c *gin.Context) {
	// Get the collection ID from the URL parameter
	collectionIDParam := c.Param("id")
	collectionID, err := strconv.Atoi(collectionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID" + err.Error()})
		return
	}

	// Delete the collection from the database
	err = core.DeleteCollectionByID(collectionID)
	if err != nil {
		if _, ok := err.(*core.DoesNotExist); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "collection not found" + err.Error()})
			return
		} else if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book" + err.Error()})
			return
		}
	}

	// Return a success response
	c.JSON(http.StatusNoContent, gin.H{"message": "Collection deleted successfully"})
}

// CountCollections retrieves the total number of collections.
func CountCollections(c *gin.Context) {
	count, err := core.CountCollections()
	if err != nil {
		if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the number of collections" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// AddBookToCollection adds a book to a collection.
func AddBookToCollection(c *gin.Context) {
	// Bind the JSON request body to a CollectionBook struct
	var collectionBook = models.CollectionBook{CollectionID: -1, BookID: -1}
	if err := c.ShouldBindJSON(&collectionBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data. " + err.Error()})
		return
	}

	// Validate the required fields (e.g., CollectionID and BookID)
	if collectionBook.CollectionID == 0 || collectionBook.BookID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CollectionID and BookID are required"})
		return
	}

	// Perform adding the book to the collection in the database
	err := core.AddBookToCollection(collectionBook.CollectionID, collectionBook.BookID)
	if err != nil {
		if _, ok := err.(*core.DoesNotExist); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "collection or book doesn't exist" + err.Error()})
			return
		} else if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		} else if _, ok := err.(*core.AlreadyExists); ok {
			c.JSON(http.StatusConflict, gin.H{"error": "collection-book combination already exists" + err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate book with a collection" + err.Error()})
			return
		}
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Book added to the collection successfully"})
}

// ListBooksInCollection retrieves all books in a specific collection.
func ListBooksInCollection(c *gin.Context) {
	collectionIDParam := c.Param("id")
	collectionID, err := strconv.Atoi(collectionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID" + err.Error()})
		return
	}

	books, err := core.ListBooksInCollection(collectionID)
	if err != nil {
		if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books in the collection" + err.Error()})
		return
	}

	if len(books) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No associated books for the collection"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// ListCollectionsOfBook retrieves the list of collections that a book belongs to by book ID.
func ListCollectionsOfBook(c *gin.Context) {
	bookIDParam := c.Param("id")
	bookID, err := strconv.Atoi(bookIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID" + err.Error()})
		return
	}

	collections, err := core.ListCollectionsOfBook(bookID)
	if err != nil {
		if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve collections for the book" + err.Error()})
		return
	}

	if len(collections) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No associated collections for the book"})
		return
	}

	c.JSON(http.StatusOK, collections)
}
*/
