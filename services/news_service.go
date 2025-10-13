package services

import (
	"log"
	"time"

	"dl/models"
	"dl/repositories"

	"github.com/mmcdole/gofeed"
)

type NewsService struct {
	Repo *repositories.NewsRepository
}

func NewNewsService(repo *repositories.NewsRepository) *NewsService {
	return &NewsService{Repo: repo}
}

func (s *NewsService) UpdateNews() error {
	parser := gofeed.NewParser()

	feeds := []string{
		"https://ecogosfond.kz/feed/",
	}

	var allNews []models.NewsItem

	for _, url := range feeds {
		feed, err := parser.ParseURL(url)
		if err != nil {
			log.Println("⚠️ RSS parse error:", err)
			continue
		}

		for _, item := range feed.Items {
			published := time.Now()
			if item.PublishedParsed != nil {
				published = *item.PublishedParsed
			}

			allNews = append(allNews, models.NewsItem{
				Title:       item.Title,
				Link:        item.Link,
				PublishedAt: published,
				Source:      feed.Title,
				Description: item.Description,
			})
		}
	}

	return s.Repo.SaveNews(allNews)
}

func (s *NewsService) GetAllNews() ([]models.NewsItem, error) {
	return s.Repo.GetAllNews()
}
