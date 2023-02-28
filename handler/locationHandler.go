package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
Entry point handler for Location information
*/
func LocationHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		handlePostRequest(w, r)
	case http.MethodGet:
		handleGetRequest(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' and '"+http.MethodPost+"' are supported.", http.StatusNotImplemented)
		return
	}

}

/*
Dedicated handler for POST requests
*/
func handlePostRequest(w http.ResponseWriter, r *http.Request) {

	// TODO: Check for content type

	// Instantiate decoder
	decoder := json.NewDecoder(r.Body)
	// Ensure parser fails on unknown fields (baseline way of detecting different structs than expected ones)
	// Note: This does not lead to a check whether an actually provided field is empty!
	decoder.DisallowUnknownFields()

	// Prepare empty struct to populate
	location := Location{}

	// Decode location instance --> Alternative: "err := json.NewDecoder(r.Body).Decode(&location)"
	err := decoder.Decode(&location)
	if err != nil {
		// Note: more often than not is this error due to client-side input, rather than server-side issues
		http.Error(w, "Error during decoding: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validation of input (Golang does not do that itself :()

	// TODO: Write convenience function for validation

	if location.Name == "" {
		http.Error(w, "Invalid input: Field 'Name' is empty.", http.StatusBadRequest)
		return
	}

	if location.Postcode == "" {
		http.Error(w, "Invalid input: Field 'Postcode' not found.", http.StatusBadRequest)
		return
	}

	// Field country is not required, hence no check

	emptyCoords := Coordinates{}
	if location.Geolocation == emptyCoords {
		http.Error(w, "Invalid input: Field 'Geolocation' not found.", http.StatusBadRequest)
		return
	}

	// Flat printing
	fmt.Println("Received following request:")
	fmt.Println(location)

	// Pretty printing
	output, err := json.MarshalIndent(location, "", "  ")
	if err != nil {
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
	}

	fmt.Println("Pretty printing:")
	fmt.Println(string(output))

	// TODO: Handle content (e.g., writing to DB, process, etc.)

	// Return status code (good practice)
	http.Error(w, "OK", http.StatusOK)
}

/*
Dedicated handler for GET requests
*/
func handleGetRequest(w http.ResponseWriter, r *http.Request) {

	// Create instance of content (could be read from DB, file, etc.)
	location := Location{
		Name:     "Gjøvik",
		Postcode: "2815",
		Country:  "Norway",
		Geolocation: Coordinates{
			Latitude:  60.7847024,
			Longitude: 10.6891797}}

	// Write content type header (best practice)
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	// Encode specific content --> Alternative: "err := json.NewEncoder(w).Encode(location)"
	err := encoder.Encode(location)
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Anything missing here?
}
