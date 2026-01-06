package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Sender    string             `bson:"sender"`
	Receiver  string             `bson:"receiver"`
	Message   string             `bson:"message"`
	IsRead    bool               `bson:"isRead"`
	CreatedAt string             `bson:"createdAt"`
}

type UserIdentity struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type BMessage struct {
	Sender   UserIdentity `json:"sender"`
	Receiver UserIdentity `json:"receiver"`
	Message  string       `json:"message"`
}
