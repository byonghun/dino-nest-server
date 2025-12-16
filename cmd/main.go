// Entry point for the go application
// python equivalent to main.py
// node equivalent to server.js

// differences: Go compiled language, needs EXPLICIT main() function
// Go typically separates entry point (cmd) from business logic (internal)

package main

import (
	"go-api-server/internal/router" // Import the router package
)

func main() {
    // Initialize the Gin router with all routes
    r := router.SetupRouter()

    // Start the HTTP server on port 8080
    // This will block and keep the server running
    if err := r.Run(":8080"); err != nil {
        // Log fatal error if the server fails to start
        panic("Failed to start server: " + err.Error())
    }
}