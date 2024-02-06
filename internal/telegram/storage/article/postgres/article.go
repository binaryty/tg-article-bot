package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/binaryty/tg-bot/internal/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

// New create a new article storage.
func New() *Storage {
	dataSource := "user=postgres password=postgres dbname=articles sslmode=disable"
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil
	}

	return &Storage{
		db: db,
	}
}

// Save an article in storage.
func (s *Storage) Save(ctx context.Context, article models.Article) error {
	query := `INSERT INTO articles(title, url, thumb_url, published_at) VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING;`

	if _, err := s.db.Exec(query,
		article.Title,
		article.Link,
		article.ThumbUrl,
		article.PublishedAt); err != nil {
		return err
	}

	return nil
}

// Read articles from storage.
func (s *Storage) Read() ([]models.Article, error) {
	query := `SELECT id, title, url, thumb_url, published_at FROM articles`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	articles := make([]models.Article, 0)

	for rows.Next() {
		var art models.Article
		if err := rows.Scan(&art.ID, &art.Title, &art.Link, &art.PublishedAt); err != nil {
			return nil, err
		}

		articles = append(articles, art)
	}

	return articles, nil
}

// ReadByTitle read an articles by title from storage.
func (s *Storage) ReadByTitle(title string) ([]models.Article, error) {
	query := `SELECT id, title, url, thumb_url, published_at FROM articles WHERE LOWER(title) LIKE '%' || $1 || '%'	ORDER BY published_at DESC;`

	rows, err := s.db.Query(query, title)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	articles := make([]models.Article, 0)

	for rows.Next() {
		var art models.Article
		if err := rows.Scan(&art.ID, &art.Title, &art.Link, &art.ThumbUrl, &art.PublishedAt); err != nil {
			return nil, err
		}

		articles = append(articles, art)
	}

	return articles, nil
}

// ReadRandom read a random article from storage.
func (s *Storage) ReadRandom() (*models.Article, error) {
	query := `SELECT id, title, url, thumb_url, published_at FROM articles ORDER BY RANDOM() LIMIT 1`

	var art models.Article

	err := s.db.QueryRow(query).Scan(&art.ID, &art.Title, &art.Link, &art.ThumbUrl, &art.PublishedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no saved articles: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random article: %v", err)
	}

	return &art, nil
}
