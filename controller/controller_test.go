package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"URL_SHORTENER/models"
	"URL_SHORTENER/storage"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCreateShortUrl(t *testing.T) {

	var routePrefix = "/api/short"
	t.Run("Empty Original URL", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		params := &CreateShortUrlRequestParams{}
		params.OriginalUrl = ""
		jsonBody, _ := json.Marshal(params)
		req := httptest.NewRequest(http.MethodPost, routePrefix, bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		// Set up the router
		router := mux.NewRouter()
		router.HandleFunc(routePrefix, CreateShortUrl).Methods("POST")
		router.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("Original url already exists", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		originalUrl := "http://example.com"

		// API request body
		params := &CreateShortUrlRequestParams{
			OriginalUrl: originalUrl,
		}
		jsonBody, _ := json.Marshal(params)
		mockErr := errors.New(storage.ErrURLAlreadyShortened)

		// Setup expectations
		resources.MockDb.EXPECT().InsertUrl(gomock.Any()).Times(1).Return(mockErr)

		req := httptest.NewRequest(http.MethodPost, routePrefix, bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		// Create the API router
		router := mux.NewRouter()
		router.HandleFunc(routePrefix, CreateShortUrl).Methods("POST")
		router.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusConflict, res.StatusCode)
	})

	t.Run("Successful short URL generation", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		originalUrl := "http://example.com"
		createdAt := time.Now().Format(YYYYMMDDhhmmss)

		// API request body
		params := &CreateShortUrlRequestParams{
			OriginalUrl: originalUrl,
		}
		jsonBody, _ := json.Marshal(params)

		// Setup expectations
		resources.MockDb.EXPECT().InsertUrl(gomock.Any()).Times(1).Return(nil)

		req := httptest.NewRequest(http.MethodPost, routePrefix, bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		// Set up the router
		router := mux.NewRouter()
		router.HandleFunc(routePrefix, CreateShortUrl).Methods("POST")
		router.ServeHTTP(w, req)

		res := w.Result()
		var responseBody ShortUrlResponse
		err := json.NewDecoder(res.Body).Decode(&responseBody)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		require.Equal(t, originalUrl, responseBody.OriginalUrl)
		require.Equal(t, createdAt, responseBody.CreatedAt)
		require.NotNil(t, responseBody.ShortUrl)
	})

}

func TestRedirectUrl(t *testing.T) {
	var endpoint = "/api/short/{short_url}"

	t.Run("Original URL not found", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		shortUrl := "esd87df7"
		mockErr := errors.New(storage.ErrShortURLDoesNotExist)
		resources.MockDb.EXPECT().GetOriginalUrl(shortUrl).Times(1).Return(nil, mockErr)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/short/%s", shortUrl), nil)
		w := httptest.NewRecorder()

		// Set up the router
		router := mux.NewRouter()
		router.HandleFunc(endpoint, RedirectUrl)
		router.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusNotFound, res.StatusCode)

	})

	t.Run("Successful URL redirection", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		shortUrl := "esd87df7"
		originalUrl := "http://example.com"
		createdAt := time.Now().Format(YYYYMMDDhhmmss)

		mockUrlRes := &models.Url{
			ShortUrl:    shortUrl,
			OriginalUrl: originalUrl,
			CreatedAt:   createdAt,
		}

		resources.MockDb.EXPECT().GetOriginalUrl(shortUrl).Times(1).Return(mockUrlRes, nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/short/%s", shortUrl), nil)
		w := httptest.NewRecorder()

		// Set up the router
		router := mux.NewRouter()
		router.HandleFunc(endpoint, RedirectUrl)
		router.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		var responseBody ShortUrlResponse
		err := json.NewDecoder(res.Body).Decode(&responseBody)
		require.NoError(t, err)
		require.Equal(t, originalUrl, responseBody.OriginalUrl)
		require.Equal(t, shortUrl, responseBody.ShortUrl)
	})

}

func TestDeleteShortUrl(t *testing.T) {
	var endPoint = "/api/short/{short_url}"

	t.Run("Original URL not found", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		mockShortUrl := "esd87df7"
		// Setup expectations
		mockErr := errors.New(storage.ErrShortURLDoesNotExist)
		resources.MockDb.EXPECT().DeleteShortUrl(mockShortUrl).Times(1).Return(mockErr)
		// Create API request
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/short/%s", mockShortUrl), nil)
		w := httptest.NewRecorder()
		// Set up the router
		router := mux.NewRouter()
		router.HandleFunc(endPoint, DeleteShortUrl).Methods("DELETE")
		router.ServeHTTP(w, req)

		// Validate the response
		res := w.Result()
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("Successful Short URL deletion", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		mockShortUrl := "esd87df7"
		// Setup expectations
		resources.MockDb.EXPECT().DeleteShortUrl(mockShortUrl).Times(1).Return(nil)
		// Create API request
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/short/%s", mockShortUrl), nil)
		w := httptest.NewRecorder()
		// Set up the router
		router := mux.NewRouter()
		router.HandleFunc(endPoint, DeleteShortUrl).Methods("DELETE")
		router.ServeHTTP(w, req)

		// Validate the response
		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

}

func TestUpdateShortUrl(t *testing.T) {
	var endPoint = "/api/short/{short_url}"

	t.Run("Original URL not found", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		mockShortUrl := "esd87df7"
		// Setup expectations
		mockErr := errors.New(storage.ErrShortURLDoesNotExist)
		resources.MockDb.EXPECT().UpdateShortUrl(gomock.Any(), mockShortUrl, gomock.Any()).Times(1).Return(mockErr)
		// Create API request
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/short/%s", mockShortUrl), nil)
		w := httptest.NewRecorder()
		// Set up the router
		router := mux.NewRouter()
		router.HandleFunc(endPoint, UpdateShortUrl).Methods("PUT")
		router.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("Successful short url updation", func(t *testing.T) {
		resources := SetupTestDB(t)
		defer resources.TearDown()

		mockShortUrl := "esd87df7"
		// Setup expectations
		resources.MockDb.EXPECT().UpdateShortUrl(gomock.Any(), mockShortUrl, gomock.Any()).Times(1).Return(nil)
		// Create API request
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/short/%s", mockShortUrl), nil)
		w := httptest.NewRecorder()
		// Set up the router
		router := mux.NewRouter()
		router.HandleFunc(endPoint, UpdateShortUrl).Methods("PUT")
		router.ServeHTTP(w, req)

		res := w.Result()
		responseBody := &UpdateShortUrlResponse{}
		err := json.NewDecoder(res.Body).Decode(res)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		require.NotNil(t, responseBody)
	})
}

func TestGenerateUniqueShortUrl(t *testing.T) {
	var length = 8
	shortUrl1 := generateUniqueShortUrl(length)
	shortUrl2 := generateUniqueShortUrl(length)

	// Check if the length of the generated short URL is correct
	if len(shortUrl1) != length {
		t.Errorf("Expected URL length of %d, but got %d", length, len(shortUrl1))
	}

	// Check if the generated URLs are unique
	if shortUrl1 == shortUrl2 {
		t.Errorf("Expected unique URLs, but got identical URLs: %s and %s", shortUrl1, shortUrl2)
	}

	// Check if the URL contains only allowed characters
	for _, char := range shortUrl1 {
		if !strings.ContainsRune(charSet, char) {
			t.Errorf("Generated URL contains invalid character: %c", char)
		}
	}
}
