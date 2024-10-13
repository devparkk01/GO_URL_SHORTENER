package controller

import (
	"URL_SHORTENER/models"
	"URL_SHORTENER/storage"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

const (
	HOSTURL             = "http://localhost:"
	PathParamShortUrlId = "short_url"
	YYYYMMDDhhmmss      = "2006-01-02 15:04:05"
	HostUrl
)

var store *storage.URLStore

func Init(urlStore *storage.URLStore) {
	store = urlStore
}

type ShortUrlResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
	CreatedAt   string `json:"created_at"`
}

type CreateShortUrlRequestParams struct {
	OriginalUrl string `json:"original_url"`
}

type UpdateShortUrlResponse struct {
	UpdatedShortUrl string `json:"updated_short_url"`
}

func CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	params := new(CreateShortUrlRequestParams)
	err := json.NewDecoder(r.Body).Decode(&params)
	fmt.Printf("Long url : {%s}", params.OriginalUrl)
	if err != nil {
		ServerResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request parameters"})
		return
	}
	err = validateShortenUrlParams(params)
	if err != nil {
		ServerResponse(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	shortUrl := generateUniqueShortUrl(8)
	createdAt := time.Now().Format(YYYYMMDDhhmmss)

	// Convert to DB request
	url := &models.Url{
		OriginalUrl: params.OriginalUrl,
		ShortUrl:    shortUrl,
		CreatedAt:   createdAt,
	}

	err = store.InsertUrl(url)
	if err != nil {
		ServerResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating url"})
		return
	}
	// convert DB response to API response
	response := ShortUrlResponse{
		OriginalUrl: url.OriginalUrl,
		ShortUrl:    url.ShortUrl,
		CreatedAt:   url.CreatedAt,
	}
	ServerResponse(w, http.StatusCreated, response)
}

func RedirectUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl, err := ParsePathParam(r, PathParamShortUrlId)
	if err != nil {
		ServerResponse(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	fmt.Printf("short url %s\n", shortUrl)
	url, err := store.GetOriginalUrl(shortUrl)
	if err != nil {
		ServerResponse(w, http.StatusNotFound, ErrorResponse{Error: "Short Url does not exist"})
		return
	}
	fmt.Println(url)

	// convert DB response to API response
	response := ShortUrlResponse{
		OriginalUrl: url.OriginalUrl,
		ShortUrl:    url.ShortUrl,
		CreatedAt:   url.CreatedAt,
	}
	fmt.Println(response)
	ServerResponse(w, http.StatusOK, response)
}

func UpdateShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl, err := ParsePathParam(r, PathParamShortUrlId)
	if err != nil {
		ServerResponse(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	log.Println(shortUrl)

	newShortUrl := generateUniqueShortUrl(8)
	// find if the URL exists in the DB
	if !store.CheckShortUrlExists(shortUrl) {
		ServerResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Short Url does not exist"})
		return
	}
	created_at := time.Now().Format(YYYYMMDDhhmmss)
	err = store.UpdateShortUrl(newShortUrl, shortUrl, created_at)
	if err != nil {
		ServerResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Error updating short url."})
		return
	}
	// convert DB response to API response
	response := UpdateShortUrlResponse{
		UpdatedShortUrl: newShortUrl,
	}
	ServerResponse(w, http.StatusCreated, response)

}

func DeleteShortUrl(w http.ResponseWriter, r *http.Request) {
	SetHeader(w, contentType, applicationJson)
	shortUrl, err := ParsePathParam(r, PathParamShortUrlId)
	if err != nil {
		ServerResponse(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	log.Println(shortUrl)
	// find if the URL exists in the DB
	if !store.CheckShortUrlExists(shortUrl) {
		ServerResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Short Url does not exist"})
		return
	}
	err = store.DeleteShortUrl(shortUrl)
	if err != nil {
		ServerResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Error deleting short url."})
		return
	}
	ServerResponse(w, http.StatusOK, "Deletion Successful.")
}

// generateUniqueShortUrl generates unique short Url of given length
func generateUniqueShortUrl(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortUrl := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		shortUrl[i] = charset[charIndex.Int64()]
	}
	return fmt.Sprintf(HOSTURL, string(shortUrl))
}

func checkOriginalUrlExists(originalUrl string) bool {
	return store.CheckOriginalUrlExists(originalUrl)
}

func validateShortenUrlParams(params *CreateShortUrlRequestParams) error {
	if params.OriginalUrl == "" {
		return errors.New("Original Url can not be empty")
	}
	if checkOriginalUrlExists(params.OriginalUrl) {
		return errors.New("Original url already exists")
	}
	return nil
}
