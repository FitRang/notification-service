package db

import (
	"github.com/Foxtrot-14/FitRang/notification-service/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Message struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Sender    string        `bson:"sender"`
	Receiver  string        `bson:"receiver"`
	Message   string        `bson:"message"`
	IsRead    bool          `bson:"isRead"`
	CreatedAt string        `bson:"createdAt"`
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

func ToGraphQLMessage(msg *Message) *model.Message {
	return &model.Message{
		ID:        msg.ID.Hex(),
		Sender:    msg.Sender,
		Receiver:  msg.Receiver,
		Message:   msg.Message,
		IsRead:    msg.IsRead,
		CreatedAt: msg.CreatedAt,
	}
}
