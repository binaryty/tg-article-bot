package article

import (
	"context"
	"github.com/binaryty/tg-bot/internal/models"
)

// ArticlesStorage interface.
type ArticlesStorage interface {
	Save(context.Context, models.Article) error
	Read() ([]models.Article, error)
	ReadByTitle(string) ([]models.Article, error)
	ReadRandom() (*models.Article, error)
}
