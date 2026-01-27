package services

import (
	"context"

	"github.com/Foxtrot-14/FitRang/notification-service/apperror"
	"github.com/Foxtrot-14/FitRang/notification-service/db"
	"github.com/Foxtrot-14/FitRang/notification-service/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MessageService) GetUnreadMessages(ctx context.Context) ([]*model.Message, error) {
	emailID, err := getEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"receiver": emailID,
		"isRead":   false,
	}
	cursor, err := m.Repo.Col.Find(ctx, filter)
	if err != nil {
		return nil, apperror.Wrap(
			apperror.Internal,
			"Failed to fetch messages",
			err,
		)
	}
	defer cursor.Close(ctx)

	var dbMessages []*db.Message
	if err := cursor.All(ctx, &dbMessages); err != nil {
		return nil, apperror.Wrap(
			apperror.Internal,
			"Failed to decode messages",
			err,
		)
	}

	if len(dbMessages) == 0 {
		return nil, apperror.New(
			apperror.NotFound,
			"No Messages Found",
		)
	}

	result := make([]*model.Message, 0, len(dbMessages))
	for _, msg := range dbMessages {
		result = append(result, db.ToGraphQLMessage(msg))
	}

	return result, nil
}
