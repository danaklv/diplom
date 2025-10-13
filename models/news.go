package models

import "time"


type NewsItem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	PublishedAt time.Time `json:"published_at"`
	Source      string    `json:"source"`
	Description string    `json:"description"`
}