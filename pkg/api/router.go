package api

import (
	"library/api/handlers"
	"library/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db db.Database) *gin.Engine {
	router := gin.Default()
	router.Use(DBMiddleware(db.DB))

	// Welcome page route
	router.GET("/", welcomePageHandler)

	// Health check route
	router.GET("/health", healthCheckHandler)
	router.GET("/api/", healthCheckHandler)

	// API routes for version 1
	v1 := router.Group("/api/v1")
	{

		v1.GET("/", healthCheckHandler)
		v1.GET("/health", healthCheckHandler)

		// Books routes
		v1.POST("/books", handlers.AddBook)
		v1.GET("/books/:id", handlers.GetBook)
		v1.GET("/books", handlers.ListBooks)
		v1.PUT("/books/:id", handlers.UpdateBook)
		v1.DELETE("/books/:id", handlers.DeleteBook)
		v1.GET("/books/search", handlers.SearchBooks)
		v1.GET("/books/count", handlers.CountBooks)
		/*
			// Collections routes
			v1.POST("/collections", handlers.AddCollection)
			v1.GET("/collections/:id", handlers.GetCollection)
			v1.GET("/collections", handlers.ListCollections)
			v1.PUT("/collections/:id", handlers.UpdateCollection)
			v1.DELETE("/collections/:id", handlers.DeleteCollection)
			v1.GET("/collections/count", handlers.CountCollections)

			v1.POST("/collections/:id/books/add", handlers.AddBookToCollection)
			v1.GET("/books/:id/collections", handlers.ListCollectionsOfBook)
			v1.GET("/collections/:id/books", handlers.ListBooksInCollection)

			// Genres routes
			v1.POST("/genres", handlers.AddGenre)
			v1.GET("/genres/:id", handlers.GetGenre)
			v1.GET("/genres", handlers.ListGenres)
			v1.PUT("/genres/:id", handlers.UpdateGenre)
			v1.DELETE("/genres/:id", handlers.DeleteGenre)
			v1.GET("/genres/count", handlers.CountGenres) */
	}

	// Default route for 404 Not Found
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page Not Found"})
	})

	return router
}

// Welcome page handler
func welcomePageHandler(c *gin.Context) {
	// You can serve the index.html or any other welcome page here
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Book Management API!",
	})
}

// Health check handler
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
