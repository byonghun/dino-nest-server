package router

// What is Gin?
// Gin is a high-performance HTTP web framework for Go (Golang).
// It provides a simple and efficient way to build web applications and APIs.
// Gin is known for its speed, minimalism, and ease of use, making it a popular choice among Go developers.

import (
	"net/http"

	"go-api-server/internal/handler"

	"github.com/gin-gonic/gin"
)

// GetHandler handles HTTP GET requests using Gin context.
// It responds with a JSON object containing a simple greeting message.
// Parameter:
//   - c: *gin.Context containing request/response information
func GetHandler(c *gin.Context) {
    // Prepare the response as a map
    response := map[string]string{"message": "Hello, World!"}
    // Gin's JSON method automatically sets Content-Type and encodes JSON
    c.JSON(http.StatusOK, response)
}

// PostHandler handles HTTP POST requests using Gin context.
// It reads JSON data from the request body and echoes it back.
// Parameter:
//   - c: *gin.Context containing request/response information
func PostHandler(c *gin.Context) {
    var requestData map[string]interface{}
    // Bind JSON from request body to requestData
    if err := c.ShouldBindJSON(&requestData); err != nil {
        // If binding fails, respond with 400 Bad Request
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Prepare the response, echoing back the received data
    response := map[string]interface{}{"received": requestData}
    // Respond with 201 Created status
    c.JSON(http.StatusCreated, response)
}

// SetupRouter initializes and configures the Gin router with all routes.
// Returns a pointer to the configured gin.Engine instance.
func SetupRouter() *gin.Engine {
    r := gin.Default() // Create a new Gin router with default middleware (logger and recovery)

    // Register a GET route at "/get" and associate it with GetHandler.
    r.GET("/get", GetHandler)

    // Register a POST route at "/post" and associate it with PostHandler.
    r.POST("/post", PostHandler)

    // Authentication routes
    // These endpoints handle user registration, login, and logout

    // POST /signup - Register a new user account
    // Expects: { "email": "user@example.com", "password": "password123" }
    // Returns: JWT token and user info
    r.POST("/signup", handler.SignupHandler)

    // POST /login - Authenticate an existing user
    // Expects: { "email": "user@example.com", "password": "password123" }
    // Returns: JWT token and user info
    r.POST("/login", handler.LoginHandler)

    // POST /logout - Log out the current user
    // Expects: Authorization header with Bearer token
    // Returns: Success message
    r.POST("/logout", handler.LogoutHandler)

    // GET /users - List all registered users
    // Returns: List of users and count of users
    r.GET("/users", handler.ListUsersHandler)

    // GET /users/search - Get user by email using query parameter
    // Example: /users/search?email=user@example.com
    r.GET("/users/search", handler.GetUserByEmailHandler)

    // GET /users/:id - Get user by their unique ID
    // Example: /users/123e4567-e89b-12d3-a456-426614174000
    r.GET("/users/:id", handler.GetUserByIDHandler)

    // Return the configured router so it can be used to start the HTTP server.
    return r
}