package handler

import (
    "net/http"
    "encoding/json"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "Hello, World!"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
    var requestData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    response := map[string]interface{}{"received": requestData}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}