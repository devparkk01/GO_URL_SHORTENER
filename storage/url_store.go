package storage

import (
	"URL_SHORTENER/models"
	"database/sql"
)

type URLStore struct {
	db *sql.DB
}

func NewURLStore(dbPath string) (*URLStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &URLStore{
		db: db,
	}, nil
}

func (s *URLStore) InsertUrl(url *models.Url) error {
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
	deleteUrlQuery := `DELETE FROM urls WHERE short_url = ?`
	_, err := s.db.Exec(deleteUrlQuery, shortUrl)
	if err != nil {
		return err
	}
	return nil
}

func (s *URLStore) UpdateShortUrl(updatedShortUrl string, shortUrl string, created_at string) error {
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
