package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MessageRepository struct {
	Col *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) *MessageRepository {
	return &MessageRepository{
		Col: db.Collection("messages"),
	}
}

func (r *MessageRepository) InitIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys: bson.M{"receiver": 1},
		},
		{
			Keys: bson.M{"receiver": 1, "isRead": 1},
		},
	}
	_, err := r.Col.Indexes().CreateMany(ctx, models)
	return err
}
