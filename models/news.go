package models

import "time"

type NewsItem struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	PublishedAt time.Time `json:"published_at"`
	Source      string    `json:"source"`
	Description string    `json:"description"`

	ImageURL string `json:"image_url,omitempty"` // если RSS поддерживает
	Category string `json:"category,omitempty"`  // если будет классификация
}
