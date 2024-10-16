package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"URL_SHORTENER/controller"
	"URL_SHORTENER/models"
	"URL_SHORTENER/storage"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var endpoint = "/api/short"

func SetupTestDB(t *testing.T) (*mux.Router, *storage.URLStore) {
	// Create an in-memory database
	dbPath := ":memory:"
	_ = os.Setenv("DB_PATH", dbPath)
	store, err := storage.NewURLStore()
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

	_ = httptest.NewServer(router)
	return router, store
}

func TestCreateShortUrlIntegration(t *testing.T) {
	router, store := SetupTestDB(t)
	defer store.Close()

	originalUrl := "http://example.com"
	params := &controller.CreateShortUrlRequestParams{
		OriginalUrl: originalUrl,
	}
	jsonBody, _ := json.Marshal(params)

	// Perform the POST request
	req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()

	// Unmarshal the response
	var createShortResp controller.ShortUrlResponse
	_ = json.NewDecoder(result.Body).Decode(&createShortResp)

	// Verify the response
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.Equal(t, originalUrl, createShortResp.OriginalUrl)
	require.NotEmpty(t, createShortResp.ShortUrl)
	require.NotEmpty(t, createShortResp.CreatedAt)
}

func TestRedirectUrlIntegration(t *testing.T) {
	router, store := SetupTestDB(t)
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

	// Perform the GET request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", endpoint, shortUrl), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()

	// Unmarshal the response
	redirectUrlResp := &controller.ShortUrlResponse{}
	_ = json.NewDecoder(result.Body).Decode(&redirectUrlResp)

	// Verify the response
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Equal(t, originalUrl, redirectUrlResp.OriginalUrl)
	require.Equal(t, shortUrl, redirectUrlResp.ShortUrl)
	require.Equal(t, createdAt, redirectUrlResp.CreatedAt)
}

func TestDeleteShortUrlIntegration(t *testing.T) {
	router, store := SetupTestDB(t)
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
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", endpoint, shortUrl), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()

	// Verify the response
	require.Equal(t, http.StatusOK, result.StatusCode)
	// Check the URL has been deleted from the database
	exists := store.CheckShortUrlExists(shortUrl)
	require.False(t, exists)
}

func TestUpdateShortUrlIntegration(t *testing.T) {
	router, store := SetupTestDB(t)
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
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", endpoint, shortUrl), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()

	// Unmarshal the response
	updateShortUrlRes := &controller.UpdateShortUrlResponse{}
	_ = json.NewDecoder(result.Body).Decode(updateShortUrlRes)
	updatedShortUrl := updateShortUrlRes.UpdatedShortUrl

	// Verify the response
	require.Equal(t, http.StatusCreated, result.StatusCode)
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

func TestCreateUpdateRetrieveShortUrl(t *testing.T) {
	router, store := SetupTestDB(t)
	defer store.Close()

	// Step 1: Create a new short URL
	originalUrl := "http://example.com"
	createShortUrlParams := &controller.CreateShortUrlRequestParams{
		OriginalUrl: originalUrl,
	}
	jsonBody, _ := json.Marshal(createShortUrlParams)
	req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()

	// Unmarshal the response
	createShortUrlResp := &controller.ShortUrlResponse{}
	_ = json.NewDecoder(result.Body).Decode(createShortUrlResp)

	// Verify the response
	require.Equal(t, http.StatusCreated, result.StatusCode)
	require.Equal(t, originalUrl, createShortUrlResp.OriginalUrl)
	shortUrl := createShortUrlResp.ShortUrl
	require.NotNil(t, shortUrl)

	// Step 2: Update the short url
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", endpoint, shortUrl), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result = w.Result()

	// Unmarshal the response
	updateShortUrlResp := &controller.UpdateShortUrlResponse{}
	_ = json.NewDecoder(result.Body).Decode(updateShortUrlResp)

	// Verify the response
	require.Equal(t, http.StatusCreated, result.StatusCode)
	updatedShortUrl := updateShortUrlResp.UpdatedShortUrl
	require.NotNil(t, updatedShortUrl)

	// Step 3: Retrieve the updated short url
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", endpoint, updatedShortUrl), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result = w.Result()

	// Unmarshal the response
	getShortUrlResp := &controller.ShortUrlResponse{}
	_ = json.NewDecoder(result.Body).Decode(getShortUrlResp)

	// Verify the response
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Equal(t, updatedShortUrl, getShortUrlResp.ShortUrl)
	require.Equal(t, originalUrl, getShortUrlResp.OriginalUrl)
}

func TestConcurrentAccess(t *testing.T) {
	router, urlStore := SetupTestDB(t)
	defer urlStore.Close()

	var mutex sync.Mutex
	var wg sync.WaitGroup
	var shortUrl = ""
	var updatedShortUrl = ""
	var errorCount = 0
	originalUrl := "http://example.com"
	concurrencyCount := 5

	// Concurrently insert URL
	for i := 0; i < concurrencyCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			createShortUrlParams := &controller.CreateShortUrlRequestParams{
				OriginalUrl: originalUrl,
			}
			jsonBody, _ := json.Marshal(createShortUrlParams)

			req := httptest.NewRequest(http.MethodPost, "/api/short", bytes.NewBuffer(jsonBody))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()
			shortUrlRes := &controller.ShortUrlResponse{}
			_ = json.NewDecoder(res.Body).Decode(shortUrlRes)
			// Count the no of instances error is Encountered
			if res.StatusCode == http.StatusConflict {
				mutex.Lock()
				errorCount++
				mutex.Unlock()
			} else {
				shortUrl = shortUrlRes.ShortUrl
			}
		}(i)
	}
	wg.Wait()

	// Verify no of instances error encountered is (total concurrency - 1)
	require.Equal(t, concurrencyCount-1, errorCount)

	// Verify if the original URL exists
	require.True(t, urlStore.CheckOriginalUrlExists(originalUrl))
	// Verify if the short URL exists
	require.True(t, urlStore.CheckShortUrlExists(shortUrl))

	// reset error count
	errorCount = 0

	// Concurrently update short URL
	for i := 0; i < concurrencyCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodPut, "/api/short/"+shortUrl, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			res := w.Result()
			updateShortUrlResp := &controller.UpdateShortUrlResponse{}
			_ = json.NewDecoder(res.Body).Decode(updateShortUrlResp)
			// Count the no of instances error is Encountered
			if res.StatusCode == http.StatusNotFound {
				mutex.Lock()
				errorCount++
				mutex.Unlock()
			} else {
				updatedShortUrl = updateShortUrlResp.UpdatedShortUrl
			}
		}(i)
	}
	wg.Wait()

	// Verify no of instances error encountered is (total concurrency - 1)
	require.Equal(t, concurrencyCount-1, errorCount)

	// Verify if the updated URL exists
	require.True(t, urlStore.CheckShortUrlExists(updatedShortUrl))

	// reset count to 0
	errorCount = 0

	// Concurrently Delete short URL
	for i := 0; i < concurrencyCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			deleteReq := httptest.NewRequest(http.MethodDelete, "/api/short/"+updatedShortUrl, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, deleteReq)
			res := w.Result()
			// Count the no of instances error is Encountered
			if res.StatusCode == http.StatusNotFound {
				mutex.Lock()
				errorCount++
				mutex.Unlock()
			}
		}(i)
	}
	wg.Wait()
	// Verify no of instances error encountered is (total concurrency - 1)
	require.Equal(t, concurrencyCount-1, errorCount)
	// Verify the URL has been deleted
	require.False(t, urlStore.CheckShortUrlExists(shortUrl))
	require.False(t, urlStore.CheckShortUrlExists(updatedShortUrl))
	require.False(t, urlStore.CheckOriginalUrlExists(originalUrl))
}
