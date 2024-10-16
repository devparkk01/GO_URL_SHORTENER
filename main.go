package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"URL_SHORTENER/controller"
	"URL_SHORTENER/storage"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	routePrefix := "/api/short"
	port := ":8080"
	dbPath := "database.sqlite3"
	_ = os.Setenv("DB_PATH", dbPath)
	// Get a new URL store
	store, err := storage.NewURLStore()
	if err != nil {
		log.Fatal(err)
	}
	controller.Init(store)

	defer store.Close()
	// Initialise Router
	r := mux.NewRouter()
	// Register all the endpoints
	// Handler to shorten the URL
	r.HandleFunc(routePrefix, controller.CreateShortUrl).Methods("POST")
	// Handler to redirect shorten url to the original url
	r.HandleFunc(routePrefix+fmt.Sprintf("/{%s}", controller.PathParamShortUrlId), controller.RedirectUrl).Methods("GET")
	// Handler to update shorten url
	r.HandleFunc(routePrefix+fmt.Sprintf("/{%s}", controller.PathParamShortUrlId), controller.UpdateShortUrl).Methods("PUT")
	// Handler to delete shorten url
	r.HandleFunc(routePrefix+fmt.Sprintf("/{%s}", controller.PathParamShortUrlId), controller.DeleteShortUrl).Methods("DELETE")

	// Listen and Serve the request
	log.Fatal(http.ListenAndServe(port, r))
}
