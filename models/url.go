package models

type Url struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}
