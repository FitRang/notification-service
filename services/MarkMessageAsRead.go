package services

import (
	"context"

	"github.com/Foxtrot-14/FitRang/notification-service/apperror"
	"github.com/Foxtrot-14/FitRang/notification-service/db"
	"github.com/Foxtrot-14/FitRang/notification-service/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m *MessageService) MarkMessageAsRead(
	ctx context.Context,
	messageID string,
) (*model.Message, error) {
	emailID, err := getEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	objID, err := bson.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, apperror.New(
			apperror.BadInput,
			"Invalid message ID",
		)
	}

	filter := bson.M{
		"_id":      objID,
		"receiver": emailID,
	}

	update := bson.M{
		"$set": bson.M{
			"isRead": true,
		},
	}

	var updated db.Message
	err = m.Repo.Col.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().
			SetReturnDocument(options.After),
	).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.New(
				apperror.NotFound,
				"Message not found",
			)
		}

		return nil, apperror.Wrap(
			apperror.Internal,
			"Failed to mark message as read",
			err,
		)
	}

	return db.ToGraphQLMessage(&updated), nil
}
