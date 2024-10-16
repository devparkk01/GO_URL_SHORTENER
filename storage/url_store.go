package storage

import (
	"database/sql"
	"errors"
	"os"
	"sync"

	"URL_SHORTENER/models"
)

const SQLITE = "sqlite3"

type URLStore struct {
	db    *sql.DB    // Sqlite database
	mutex sync.Mutex // Mutex for thread safety
}

type URLOperations interface {
	InsertUrl(url *models.Url) error
	CheckShortUrlExists(shortUrl string) bool
	CheckOriginalUrlExists(originalUrl string) bool
	GetOriginalUrl(shortUrl string) (*models.Url, error)
	DeleteShortUrl(shortUrl string) error
	UpdateShortUrl(updatedShortUrl string, shortUrl string, created_at string) error
}

func NewURLStore(dbPath string) (*URLStore, error) {
	dbPath = os.Getenv("DB_PATH")
	if dbPath == "" {
		return nil, errors.New("DB_PATH environment variable not set")
	}
	db, err := sql.Open(SQLITE, dbPath)
	if err != nil {
		return nil, err
	}
	// Create the table if not already created
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS "urls" (
				original_url TEXT PRIMARY KEY NOT NULL,
				short_url TEXT NOT NULL,
				created_at TEXT NOT NULL
			);
		`)
	if err != nil {
		return nil, errors.New("Failed to create table: " + err.Error())
	}
	// Create index on the short_url
	_, err = db.Exec(`
				CREATE index IF NOT EXISTS idx_short ON urls (short_url);
		`)
	if err != nil {
		return nil, errors.New("Failed to create index: " + err.Error())
	}
	return &URLStore{
		db: db,
	}, nil
}

func (s *URLStore) InsertUrl(url *models.Url) error {
	// Lock the mutex before performing insert operation
	s.mutex.Lock()
	defer s.mutex.Unlock()
	insertUrlQuery := `INSERT INTO urls (original_url, short_url, created_at) VALUES (?, ?, ?)`
	_, err := s.db.Exec(insertUrlQuery, url.OriginalUrl, url.ShortUrl, url.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *URLStore) CheckShortUrlExists(shortUrl string) bool {
	checkShortUrlQuery := `SELECT short_url FROM urls WHERE short_url = ?`
	err := s.db.QueryRow(checkShortUrlQuery, shortUrl).Scan(&shortUrl)
	if err == nil {
		return true
	}
	return false
}
func (s *URLStore) CheckOriginalUrlExists(originalUrl string) bool {
	checkOriginalUrlQuery := `SELECT original_url from urls WHERE original_url= ?`
	err := s.db.QueryRow(checkOriginalUrlQuery, originalUrl).Scan(&originalUrl)
	if err == nil {
		return true
	}
	return false
}

func (s *URLStore) GetOriginalUrl(shortUrl string) (*models.Url, error) {
	getOriginalUrlQuery := `SELECT original_url, short_url, created_at FROM urls WHERE short_url = ?`
	var url models.Url
	err := s.db.QueryRow(getOriginalUrlQuery, shortUrl).Scan(&url.OriginalUrl, &url.ShortUrl, &url.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (s *URLStore) DeleteShortUrl(shortUrl string) error {
	// Lock the mutex before performing delete operation
	s.mutex.Lock()
	defer s.mutex.Unlock()
	deleteUrlQuery := `DELETE FROM urls WHERE short_url = ?`
	_, err := s.db.Exec(deleteUrlQuery, shortUrl)
	if err != nil {
		return err
	}
	return nil
}

func (s *URLStore) UpdateShortUrl(updatedShortUrl string, shortUrl string, created_at string) error {
	// Lock the mutex before performing update operation
	s.mutex.Lock()
	defer s.mutex.Unlock()
	updateUrlQuery := `UPDATE urls SET short_url = ?, created_at = ? WHERE short_url = ?`
	_, err := s.db.Exec(updateUrlQuery, updatedShortUrl, created_at, shortUrl)
	if err != nil {
		return err
	}
	return nil
}

func (s *URLStore) Close() {
	_ = s.db.Close()
}
