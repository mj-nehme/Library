package handlers

/*
import (
	"library/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddGenre handles the "POST /api/v1/genres" endpoint to create a new genre.
func AddGenre(c *gin.Context) {
	// Bind the JSON request body to a Genre struct
	var newGenre models.Genre
	if err := c.ShouldBindJSON(&newGenre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data. " + err.Error()})
		return
	}

	// Validate the required fields (e.g., title and author)
	err := validators.AddGenre(&newGenre)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Perform genre creation in the database
	genreID, err := core.AddGenre(newGenre.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre" + err.Error()})
		return
	}
	genre := models.Genre{ID: genreID, Name: newGenre.Name}

	// Return the newly created genre as the response
	c.JSON(http.StatusCreated, genre)
}

// GetGenre handles the "GET /api/v1/genres/:id" endpoint to retrieve a specific genre by its ID.
func GetGenre(c *gin.Context) {
	genreIDParam := c.Param("id")
	genreID, err := strconv.Atoi(genreIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID" + err.Error()})
		return
	}

	genre, err := core.GetGenre(genreID)
	if err != nil {
		if _, ok := err.(*core.DoesNotExist); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre ID not found" + err.Error()})
			return
		} else if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genre" + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, genre)
}

// ListGenres handles the "GET /api/v1/genres" endpoint to retrieve a list of all genres.
func ListGenres(c *gin.Context) {
	genres, err := core.ListGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve genres" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, genres)
}

// UpdateGenre updates the information for an existing genre.
func UpdateGenre(c *gin.Context) {
	// Get the genre ID from the URL parameter
	genreIDStr := c.Param("id")
	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID" + err.Error()})
		return
	}

	// Check if the genre with the given ID exists in the database
	existingGenre, err := core.GetGenre(genreID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found" + err.Error()})
		return
	}

	// Bind the request JSON body to a Genre object
	var updatedGenre models.Genre
	if err := c.ShouldBindJSON(&updatedGenre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body" + err.Error()})
		return
	}
	updatedGenre.ID = genreID

	// Save the updated genre in the database
	affectedRows, err := core.UpdateGenre(&updatedGenre)
	if err != nil {
		if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update genre" + err.Error()})
		return
	}
	if affectedRows == 0 {
		c.JSON(http.StatusNotModified, gin.H{"error": "No affected rows" + err.Error()})
		return
	}

	// Update the genre information
	existingGenre.Name = updatedGenre.Name
	// Update other fields as needed

	// Return the updated genre as the response
	c.JSON(http.StatusOK, existingGenre)
}

// DeleteGenre deletes a genre from the system.
func DeleteGenre(c *gin.Context) {
	// Get the genre ID from the URL parameter
	genreIDStr := c.Param("id")
	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID" + err.Error()})
		return
	}

	// Delete the genre from the database
	err = core.DeleteGenre(genreID)
	if err != nil {
		if _, ok := err.(*core.DoesNotExist); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre ID not found" + err.Error()})
			return
		} else if _, ok := err.(*core.InvalidRequest); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request" + err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete genre" + err.Error()})
			return
		}
	}

	// Return a success response
	c.JSON(http.StatusNoContent, gin.H{"message": "Genre deleted successfully"})
}

// CountGenres returns the total number of genres in the database.
func CountGenres(c *gin.Context) {
	count, err := core.CountGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count genres" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
*/
