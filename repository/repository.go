package repository

import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

var db *mongo.Database

func Init(database *mongo.Database) {
    db = database

    ctx := context.Background()

    if err := NewMessagesRepository(db).InitIndexes(ctx); err != nil {
        log.Printf("Dossier index creation failed: %v\n", err)
    }
}
