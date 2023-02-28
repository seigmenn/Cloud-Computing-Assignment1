package main

import (
	"assignment_1/handler"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var global = time.Now()

type Uni struct {
	Name         string      `json:"name"`
	Country      string      `json:"country"`
	AlphaTwoCode string      `json:"alpha_two_code"`
	WebPages     []string    `json:"web_pages"`
	Languages    interface{} `json:"languages"`
	Map          MapS        `json:"maps"`
}
type MapS struct {
	OpenStreet string `json:"openStreetMaps"`
}
type diagnostics struct {
	UniAPI     int     `json:"universitiesapi"`
	CountryAPI int     `json:"countriesapi"`
	Version    string  `json:"version"`
	Uptime     float64 `json:"uptime"`
}

func uniHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	search := "http://universities.hipolabs.com/search?name=" + parts[3]

	infoUni, err := http.Get(search)
	if err != nil {
		log.Println("Error:" + err.Error())
		return
	}
	unis := []Uni{}

	decoder := json.NewDecoder(infoUni.Body)
	//decoder.DisallowUnknownFields()

	err = decoder.Decode(&unis)
	if err != nil {
		http.Error(w, "Error during decoding: "+err.Error(), http.StatusBadRequest)
		return
	}

	for i, un := range unis {
		search = "https://restcountries.com/v3.1/name/" + un.Country + "?fields=languages,maps"
		univers := []Uni{}
		info, feil := http.Get(search)
		if feil != nil {
			http.Error(w, "Error during request: "+feil.Error(), http.StatusBadRequest)
			return
		}

		decoder = json.NewDecoder(info.Body)

		err = decoder.Decode(&univers)
		if err != nil {
			http.Error(w, "Error during decoding: "+err.Error(), http.StatusBadRequest)
			return
		}

		unis[i].Languages = univers[0].Languages
		unis[i].Map = univers[0].Map
	}
	output, err := json.MarshalIndent(unis, "", "	")
	_, err = fmt.Fprintf(w, "%v", string(output))
	if err != nil {
		http.Error(w, "Error when returning output:", http.StatusInternalServerError)
	}
	http.Error(w, "OK", http.StatusOK)
}

func neighbourUnis(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	type Names struct {
		Common string `json:"common"`
	}

	type Country struct {
		Name      Names       `json:"name"`
		Languages interface{} `json:"languages"`
		Borders   []string    `json:"borders"`
		Map       MapS        `json:"maps"`
	}

	type Border struct {
		Name Names `json:"name"`
	}
	parts := strings.Split(r.URL.Path, "/")
	cntryName := parts[4]
	uniName := parts[5]

	limit := 0
	query := r.URL.Query()
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	log.Println(limit)

	search := "http://restcountries.com/v3.1/name/" + cntryName
	info, err := http.Get(search)
	if err != nil {
		http.Error(w, "Error 1:", http.StatusBadRequest)
		return
	}

	countries := []Country{}
	decoder := json.NewDecoder(info.Body)
	err = decoder.Decode(&countries)

	if err != nil {
		http.Error(w, "Error 2:"+err.Error(), http.StatusBadRequest)
		return
	}

	ngbourUnis := []Uni{}
	counter := 0

	for _, border := range countries[0].Borders {

		search := "http://restcountries.com/v3.1/alpha/" + border
		info, err = http.Get(search)
		if err != nil {
			http.Error(w, "Error 3:", http.StatusBadRequest)
			return
		}

		borderCountries := []Border{}
		decoder = json.NewDecoder(info.Body)
		err = decoder.Decode(&borderCountries)

		if err != nil {
			http.Error(w, "Error 4:"+err.Error(), http.StatusBadRequest)
			return
		}

		search = "http://universities.hipolabs.com/search?country=" + borderCountries[0].Name.Common + "&name=" + uniName
		info, err = http.Get(search)
		if err != nil {
			http.Error(w, "Error 5:", http.StatusBadRequest)
			return
		}

		decoder = json.NewDecoder(info.Body)
		err = decoder.Decode(&ngbourUnis)
		if err != nil {
			http.Error(w, "Error 6:"+err.Error(), http.StatusBadRequest)
			return
		}

		for i, un := range ngbourUnis {
			if counter == limit {
				break
			}
			search = "https://restcountries.com/v3.1/name/" + un.Country + "?fields=languages,maps"
			univers := []Uni{}
			info, feil := http.Get(search)
			if feil != nil {
				http.Error(w, "Error during request: "+feil.Error(), http.StatusBadRequest)
				return
			}

			decoder = json.NewDecoder(info.Body)

			err = decoder.Decode(&univers)
			if err != nil {
				http.Error(w, "Error during decoding: "+err.Error(), http.StatusBadRequest)
				return
			}
			counter++
			ngbourUnis[i].Languages = univers[0].Languages
			ngbourUnis[i].Map = univers[0].Map
			log.Println(ngbourUnis[i])
			log.Println(counter)
			if counter == limit {
				break
			}
		}
		if counter == limit {
			break
		}
	}
	lim := []Uni{}
	for i, u := range ngbourUnis {
		if i == limit {
			break
		}
		lim = append(lim, u)
	}
	output, err := json.MarshalIndent(lim, "", "		")
	if err != nil {
		http.Error(w, "Error 7:", http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(w, "%v", string(output))
	if err != nil {
		http.Error(w, "Error 8:", http.StatusInternalServerError)
		return
	}

	http.Error(w, "OK", http.StatusOK)
}

func diagHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	uniResp, err := http.Get("http://universities.hipolabs.com/")
	if err != nil {
		http.Error(w, "Error with uni API:", http.StatusInternalServerError)
		return
	}

	countryResp, err := http.Get("https://restcountries.com/")
	if err != nil {
		http.Error(w, "Error with countries API:", http.StatusInternalServerError)
		return
	}

	diag := diagnostics{
		UniAPI:     uniResp.StatusCode,
		CountryAPI: countryResp.StatusCode,
		Uptime:     time.Since(global).Seconds(),
		Version:    "v1",
	}

	output, err := json.MarshalIndent(diag, "", "  ")
	if err != nil {
		http.Error(w, "Error:", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(output)
	if err != nil {
		http.Error(w, "Error when returning output:", http.StatusInternalServerError)
		return
	}

	http.Error(w, "OK", http.StatusOK)
}

func main() {
	// Handle port assignment (either based on environment variable, or local override)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}
	// Set up handler endpoints
	http.HandleFunc(handler.DEFAULT_PATH, handler.EmptyHandler)
	http.HandleFunc(handler.LOCATION_PATH, handler.LocationHandler)
	http.HandleFunc(handler.COLLECTION_PATH, handler.CollectionHandler)
	http.HandleFunc("/uniinfo/v1/", uniHandler)
	http.HandleFunc("/unisearcher/v1/diag/", diagHandler)
	http.HandleFunc("/unisearcher/v1/neighbourunis/", neighbourUnis)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
