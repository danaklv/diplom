package repositories

import (
	"database/sql"
	"dl/models"
)

type NewsRepository struct {
	DB *sql.DB
}

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{DB: db}
}

// ------------------------ SAVE NEWS ------------------------

func (r *NewsRepository) SaveNews(items []models.NewsItem) error {
	if len(items) == 0 {
		return nil
	}

	stmt, err := r.DB.Prepare(`
        INSERT INTO news (title, link, published_at, source, description)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (link) DO NOTHING
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		if _, err := stmt.Exec(
			item.Title,
			item.Link,
			item.PublishedAt,
			item.Source,
			item.Description,
		); err != nil {
			// ошибка не возвращаем, чтобы не останавливать весь парсинг
			// но собираем первую и возвращаем после цикла
			return err
		}
	}

	return nil
}

// ------------------------ GET ALL NEWS ------------------------

func (r *NewsRepository) GetAllNews() ([]models.NewsItem, error) {
	rows, err := r.DB.Query(`
        SELECT id, title, link, published_at, source, description
        FROM news
        ORDER BY published_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []models.NewsItem

	for rows.Next() {
		var n models.NewsItem
		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Link,
			&n.PublishedAt,
			&n.Source,
			&n.Description,
		); err != nil {
			return nil, err
		}
		news = append(news, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}
