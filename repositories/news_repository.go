package repositories

import (
	"database/sql"
	"dl/models"
	"log"
)

type NewsRepository struct {
	DB *sql.DB
}

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{DB: db}
}

func (r *NewsRepository) SaveNews(items []models.NewsItem) error {
	for _, item := range items {
		_, err := r.DB.Exec(`
			INSERT INTO news (title, link, published_at, source, description)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (link) DO NOTHING
		`, item.Title, item.Link, item.PublishedAt, item.Source, item.Description)
		if err != nil {
			log.Println("❌ Error inserting news:", err)
		}
	}
	return nil
}

func (r *NewsRepository) GetAllNews() ([]models.NewsItem, error) {
	rows, err := r.DB.Query(`
		SELECT id, title, link, published_at, source, description
		FROM news ORDER BY published_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []models.NewsItem
	for rows.Next() {
		var n models.NewsItem
		err := rows.Scan(&n.ID, &n.Title, &n.Link, &n.PublishedAt, &n.Source, &n.Description)
		if err != nil {
			return nil, err
		}
		news = append(news, n)
	}
	return news, nil
}
