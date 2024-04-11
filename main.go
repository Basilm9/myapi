package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
// tells us where to host our port on our machine
	const filepathRoot = "."
	const port = "8081"
// server multiplexer like a police man at a light directing traffic to the correct urls
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/weather", handlerWeatherAPI)

// wrap our mutiplexer in our cors middleare (middleware adds some piece of information to the request/response)
// think of cors like a bouncer at a club it is the dicider on what request we want to let reach our server
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "hello world" )
}

func handlerWeatherAPI(w http.ResponseWriter, r *http.Request) {
	apiKey := "61e8a860c037438e8fd19sdfjslsdf" // Use your actual API key
	baseURL := "http://api.weatherapi.com/v1/current.json"
	query := "London" // Example query, replace with your desired location

	// Construct the full API request URL with your API key and query
	fullURL := baseURL + "?key=" + apiKey + "&q=" + query

	// Make the HTTP GET request to the WeatherAPI open a response 
	// Since responses arent instant we have to make sure we fully initalize the data before we close the request to the api
	resp, err := http.Get(fullURL)
	if err != nil {
		log.Printf("Error making API request: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Respond with the API response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}