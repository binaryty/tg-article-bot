package mongo

import (
	"context"
	"github.com/binaryty/tg-bot/internal/lib/er"
	"github.com/binaryty/tg-bot/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type UserStorage struct {
	records *Records
}

type Records struct {
	*mongo.Collection
}

// New a constructor of users storage.
func New(ctx context.Context, connectString string, connectTimeout time.Duration) *UserStorage {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))
	if err != nil {
		log.Fatalf("can't connect to user storage: %v", err)
	}

	records := Records{
		Collection: client.Database("users").Collection("users"),
	}

	return &UserStorage{
		records: &records,
	}
}

// Save new user in user storage.
func (s *UserStorage) Save(ctx context.Context, user *models.User) error {
	_, err := s.records.InsertOne(ctx, user)
	if err != nil {
		return er.Wrap("can't save user", err)
	}

	return nil
}

// Read all users from user storage.
func (s *UserStorage) Read(ctx context.Context) ([]models.User, error) {
	//TODO implement me
	panic("implement me")
}
