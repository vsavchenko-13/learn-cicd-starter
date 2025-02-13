package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal the payload to JSON
	dat, err := json.Marshal(payload)
	if err != nil {
		// Log the error with more context (like the payload, if needed)
		log.Printf("Error marshalling JSON: %s, Payload: %#v", err, payload)

		// Send a response with a 500 error and a detailed error message in JSON
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": "Internal Server Error", "message": "Failed to marshal the response payload"}
		// If marshaling the error response also fails, just write a generic message
		dat, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			http.Error(w, "Failed to generate error response", http.StatusInternalServerError)
			return
		}
		// Check for error while writing the response
		if _, writeErr := w.Write(dat); writeErr != nil {
			// Log the write error
			log.Printf("Error writing error response: %s", writeErr)
		}
		return
	}

	// Set the correct status code and write the marshaled data to the response body
	w.WriteHeader(code)
	if _, writeErr := w.Write(dat); writeErr != nil {
		// Log the write error
		log.Printf("Error writing response: %s", writeErr)
	}
}
