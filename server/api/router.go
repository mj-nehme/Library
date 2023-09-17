// @title		Library API Swagger
// @version	2.0
// @description     A book management API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Mohamad Jaafar Nehme
// @contact.url    https://www.linkedin.com/in/mjnehme/
// @contact.email  Mohamad.jaafar.nehme@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8090
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

package api

import (
	"library/api/handlers"
	"library/db"
	"library/docs"
	"net/http"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/gin-gonic/gin"
)

const htmlFiles = "templates/*"

func SetupRouter(db db.Database, initialPath ...string) *gin.Engine {
	router := gin.Default()
	relativePath := getRelativePath(initialPath...)
	router.LoadHTMLGlob(relativePath)
	docs.SwaggerInfo.BasePath = "/api/v1"
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

		// Books routes
		v1.POST("/books", handlers.AddBook)
		v1.GET("/books/:id", handlers.GetBook)
		v1.GET("/books", handlers.ListBooks)
		v1.PUT("/books/:id", handlers.UpdateBook)
		v1.PATCH("/books/:id", handlers.PatchBook)
		v1.DELETE("/books/:id", handlers.DeleteBook)
		v1.GET("/books/search", handlers.SearchBooks)
		v1.GET("/books/count", handlers.CountBooks)
	}

	// Serve Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Default route for 404 Not Found
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, handlers.MessageResponse{Message: "Page Not Found"})
	})

	return router
}

//	@Summary		Welcome page
//	@Description	Welcome page for the Book Management API
//	@Tags			info
//	@Produce		json
//	@Success		200	{object}	handlers.MessageResponse	"Returns the homepage"
//	@Router			/ [get]
//
// Welcome page handler
func welcomePageHandler(c *gin.Context) {
	// Serve the welcome page HTML file
	c.HTML(http.StatusOK, "welcome.html", nil)
}

//	@Summary		Health check
//	@Description	Check the status of the Book Management API
//	@Tags			info
//	@Produce		json
//	@Success		200	{object}	handlers.MessageResponse	"Returns the status message"
//	@Router			/health [get]
//
// Health check handler
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, handlers.StatusResponse{
		Status: "ok",
	})
}

func getRelativePath(initialPath ...string) string {
	var relativePath string
	if len(initialPath) == 0 {
		relativePath = htmlFiles
	} else {
		relativePath = initialPath[0] + htmlFiles
	}
	return relativePath
}
