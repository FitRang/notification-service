package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
            Keys:    bson.M{"email": 1},
            Options: options.Index().SetUnique(true),
        },
    }

    _, err := r.Col.Indexes().CreateMany(ctx, models)
    return err
}
