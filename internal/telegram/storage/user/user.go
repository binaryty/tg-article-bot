package user

import (
	"context"
	"github.com/binaryty/tg-bot/internal/models"
)

type Storage interface {
	Save(context.Context, *models.User) error
	Read(ctx context.Context) ([]models.User, error)
}
