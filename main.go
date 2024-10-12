package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	routePrefix := "/api/short"
	pathParamShortUrlId := "short_url"
	//port := os.Getenv("PORT")
	port := ":8080"

	// Initialise Router
	r := mux.NewRouter()
	// Register all the endpoints
	// Handler to shorten the URL
	r.HandleFunc(routePrefix, ShortenUrl).Methods("POST")
	// Handler to redirect shorten url to the original url
	r.HandleFunc(routePrefix+fmt.Sprintf("/{%s}", pathParamShortUrlId), FetchUrl).Methods("GET")
	// Handler to update shorten url
	r.HandleFunc(routePrefix+fmt.Sprintf("/{%s}", pathParamShortUrlId), UpdateUrl).Methods("PUT")
	// Handler to delete shorten url
	r.HandleFunc(routePrefix+fmt.Sprintf("/{%s}", pathParamShortUrlId), DeleteUrl).Methods("DELETE")

	// Listen and Serve the request
	log.Fatal(http.ListenAndServe(port, r))
}
