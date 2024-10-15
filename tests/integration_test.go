package tests

import (
	"URL_SHORTENER/controller"
	"URL_SHORTENER/models"
	"URL_SHORTENER/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var endpoint = "/api/short"

func SetupTestDB(t *testing.T) (*httptest.Server, *storage.URLStore) {
	// Create an in-memory database
	dbPath := ":memory:"
	store, err := storage.NewURLStore(dbPath)
	require.NoError(t, err)

	// Initialize controllers with the store
	controller.Init(store)

	// Initialize the router
	router := mux.NewRouter()
	routePrefix := "/api/short"
	router.HandleFunc(routePrefix, controller.CreateShortUrl).Methods("POST")
	router.HandleFunc(routePrefix+"/{short_url}", controller.RedirectUrl).Methods("GET")
	router.HandleFunc(routePrefix+"/{short_url}", controller.UpdateShortUrl).Methods("PUT")
	router.HandleFunc(routePrefix+"/{short_url}", controller.DeleteShortUrl).Methods("DELETE")

	// Return a test server running the router
	return httptest.NewServer(router), store
}

func TestCreateShortUrlIntegration(t *testing.T) {
	server, store := SetupTestDB(t)
	defer store.Close()

	originalUrl := "http://example.com"
	params := &controller.CreateShortUrlRequestParams{
		OriginalUrl: originalUrl,
	}
	jsonBody, _ := json.Marshal(params)

	// Perform the POST request
	resp, err := http.Post(server.URL+endpoint, "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert the response
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Decode the response
	var createShortResp controller.ShortUrlResponse
	err = json.NewDecoder(resp.Body).Decode(&createShortResp)
	require.NoError(t, err)
	require.Equal(t, originalUrl, createShortResp.OriginalUrl)
	require.NotEmpty(t, createShortResp.ShortUrl)
	require.NotEmpty(t, createShortResp.CreatedAt)
}

func TestRedirectUrlIntegration(t *testing.T) {
	server, store := SetupTestDB(t)
	defer store.Close()

	originalUrl := "http://example.com"
	shortUrl := "esd87df7"
	createdAt := time.Now().Format(controller.YYYYMMDDhhmmss)

	// Manually insert the test URL into the database to simulate a previously created short URL
	err := store.InsertUrl(&models.Url{
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
		CreatedAt:   createdAt,
	})
	require.NoError(t, err)

	// Perform the GET request to redirect
	resp, err := http.Get(server.URL + fmt.Sprintf("%s/%s", endpoint, shortUrl))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert the response
	require.Equal(t, http.StatusOK, resp.StatusCode)
	redirectUrlResp := &controller.ShortUrlResponse{}
	_ = json.NewDecoder(resp.Body).Decode(&redirectUrlResp)

	// Verify that the redirect location matches the original URL
	require.Equal(t, originalUrl, redirectUrlResp.OriginalUrl)
	require.Equal(t, shortUrl, redirectUrlResp.ShortUrl)
	parsedTime, _ := time.Parse(time.RFC3339, redirectUrlResp.CreatedAt)
	require.Equal(t, createdAt, parsedTime.Format(controller.YYYYMMDDhhmmss))
}

func TestDeleteShortUrlIntegration(t *testing.T) {
	server, store := SetupTestDB(t)
	defer store.Close()

	shortUrl := "esd87df7"
	originalUrl := "http://example.com"

	// Insert the test URL into the database to simulate a previously created short URL
	err := store.InsertUrl(&models.Url{
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
		CreatedAt:   time.Now().Format(controller.YYYYMMDDhhmmss),
	})
	require.NoError(t, err)

	// Perform the DELETE request
	req, err := http.NewRequest(http.MethodDelete, server.URL+fmt.Sprintf("%s/%s", endpoint, shortUrl), nil)
	require.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert the response
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// Check the URL has been deleted from the database
	exists := store.CheckShortUrlExists(shortUrl)
	require.False(t, exists)
}

func TestUpdateShortUrlIntegration(t *testing.T) {
	server, store := SetupTestDB(t)
	defer store.Close()

	shortUrl := "esd87df7"
	originalUrl := "http://newexample.com"
	createdAt := time.Now().Format(controller.YYYYMMDDhhmmss)

	// Insert the test URL into the database to simulate a previously created short URL
	err := store.InsertUrl(&models.Url{
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
		CreatedAt:   createdAt,
	})
	require.NoError(t, err)

	// Perform the PUT request
	req, err := http.NewRequest(http.MethodPut, server.URL+fmt.Sprintf("%s/%s", endpoint, shortUrl), nil)
	require.NoError(t, err)
	//req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	updateShortUrlRes := &controller.UpdateShortUrlResponse{}
	_ = json.NewDecoder(resp.Body).Decode(updateShortUrlRes)

	updatedShortUrl := updateShortUrlRes.UpdatedShortUrl

	// Assert the response
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NotEqual(t, shortUrl, updatedShortUrl)

	// Ensure the previous short url has been removed from the database
	url, err := store.GetOriginalUrl(shortUrl)
	require.Nil(t, url)
	require.NotNil(t, err)

	// Ensure the previous short url has been removed from the database
	url, err = store.GetOriginalUrl(updatedShortUrl)
	require.NotNil(t, url)
	require.Nil(t, err)
	require.Equal(t, url.OriginalUrl, originalUrl)
}
