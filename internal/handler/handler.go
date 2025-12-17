package handler

import (
	"encoding/json"
	"net/http" // Go's built-in net/http package for HTTP handling
)

// GetHandler handles HTTP GET requests.
// It responds with a JSON object containing a simple greeting message.
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response
//   - r: *http.Request representing the incoming HTTP request
func GetHandler(w http.ResponseWriter, r *http.Request) {
    // Prepare the response as a map (will be encoded to JSON)
    response := map[string]string{"message": "Hello, World!"}
    // Set the Content-Type header to indicate JSON response
    w.Header().Set("Content-Type", "application/json")
    // Set the HTTP status code to 200 OK
    w.WriteHeader(http.StatusOK)
    // Encode the response map to JSON and write it to the response
    json.NewEncoder(w).Encode(response)
}

// PostHandler handles HTTP POST requests.
// It reads JSON data from the request body, decodes it into a map,
// and responds with a JSON object echoing the received data.
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response
//   - r: *http.Request representing the incoming HTTP request
func PostHandler(w http.ResponseWriter, r *http.Request) {
    var requestData map[string]interface{}
    // Decode the JSON request body into requestData map
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        // If decoding fails, respond with a 400 Bad Request and the error message
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Prepare the response, echoing back the received data
    response := map[string]interface{}{"received": requestData}
    // Set the Content-Type header to indicate JSON response
    w.Header().Set("Content-Type", "application/json")
    // Set the HTTP status code to 201 Created
    w.WriteHeader(http.StatusCreated)
    // Encode the response map to JSON and write it to the response
    json.NewEncoder(w).Encode(response)
}