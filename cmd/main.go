// Entry point for the go application
// python equivalent to main.py
// node equivalent to server.js

// differences: Go compiled language, needs EXPLICIT main() function
// Go typically separates entry point (cmd) from business logic (internal)

package main

import (
	"go-api-server/internal/database" // Import the database package
	"go-api-server/internal/handler"  // Import the handler package
	"go-api-server/internal/router"   // Import the router package
)

func main() {
    // Initialize the in-memory database
    // This creates a new instance of our database to store users
    // In production, you'd connect to a real database here (PostgreSQL, MySQL, MongoDB, etc.)
    handler.DB = database.NewInMemoryDB()

    // Initialize the Gin router with all routes
    // This sets up all our API endpoints (/get, /post, /signup, /login, /logout)
    r := router.SetupRouter()

    // Start the HTTP server on port 8080
    // This will block and keep the server running until interrupted
    if err := r.Run(":8080"); err != nil {
        // Log fatal error if the server fails to start
        // In production, use a proper logger instead of panic
        panic("Failed to start server: " + err.Error())
    }
}