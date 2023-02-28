package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
Entry point handler for collection information
*/
func CollectionHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		handleMapPostRequest(w, r)
	case http.MethodGet:
		handleMapGetRequest(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' and '"+http.MethodPost+"' are supported.", http.StatusNotImplemented)
		return
	}

}

/*
Dedicated handler for POST requests
*/
func handleMapPostRequest(w http.ResponseWriter, r *http.Request) {

	// Instantiate decoder
	decoder := json.NewDecoder(r.Body)

	// Open-ended map structure
	mp := map[string]interface{}{}

	// Decode location instance --> Alternative: "err := json.NewDecoder(r.Body).Decode(&mp)"
	err := decoder.Decode(&mp)
	if err != nil {
		log.Println("Error during encoding: " + err.Error())
		http.Error(w, "Error during decoding", http.StatusBadRequest)
		return
	}

	// Flat printing
	fmt.Println("Received following request:")
	fmt.Println(mp)

	// Pretty printing
	output, err := json.MarshalIndent(mp, "", "  ")
	if err != nil {
		log.Println("Error during pretty printing of output: " + err.Error())
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
	}

	fmt.Println("Pretty printing:")
	fmt.Println(string(output))

	// TODO: Handle content (e.g., writing to DB, process, etc.)

	// Return status code (good practice) - note that no content is provided in the body (and indicated via status code)
	http.Error(w, "", http.StatusNoContent)
}

/*
Dedicated handler for GET requests
*/
func handleMapGetRequest(w http.ResponseWriter, r *http.Request) {

	// Create collection (open-ended, untyped value)
	collection := map[string]interface{}{
		"first":  "firstValue",
		"second": 2,
		"third":  3.1,
	}

	// Write content type header (best practice)
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	// Encode specific content --> Alternative: "err := json.NewEncoder(w).Encode(collection)"
	err := encoder.Encode(collection)
	if err != nil {
		log.Println("Error during encoding: " + err.Error())
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}

	// Explicit specification of return status code --> will default to 200 if not provided.
	http.Error(w, "", http.StatusNoContent)
}
