package repository

import (
	"context"
	"time"
	"working/super_task/internal/domain"

	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	database               mongo.Database
	messageCollection      string
	conversationCollection string
}

func NewMessageRepos(db mongo.Database, collectionMessage string, collectionConversation string) *MessageRepository {
	return &MessageRepository{
		database:               db,
		messageCollection:      collectionMessage,
		conversationCollection: collectionConversation,
	}
}

// method for creating new conversation or retiving existing one
func (mr *MessageRepository) CreateOrUpdateConversation(cxt context.Context, participant *domain.ParticipantRequest, message domain.Message) (*domain.Conversation, error) {
	convCollection := mr.database.Collection(mr.conversationCollection)
	filter := bson.D{
		{Key: "participants", Value: participant.FirstParticipantID},
		{Key: "participants", Value: participant.SecondParticipantID},
	}

	var convModel *domain.Conversation
	err := convCollection.FindOne(cxt, filter).Decode(&convModel)
	if err != nil {
		newConvMode := &domain.Conversation{
			ConversationID: primitive.NewObjectID(),
			IsGroup:        false,
			Participants: []domain.Participant{
				{ParticipantID: participant.FirstParticipantID},
				{ParticipantID: participant.SecondParticipantID},
			},
			LastMessage: message,
			CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		}
		convModel = newConvMode
	} else {
		convModel.LastMessage = message
		convModel.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	}
	return convModel, nil

}

// method for fetching message history from the database
func (mr *MessageRepository) FetchMessages(cxt context.Context, ConvID string, page int64, size int64) ([]*domain.Message, error) {
	messageCollection := mr.database.Collection(mr.messageCollection)
	skip := (page - 1) / size
	opts := options.Find().SetSkip(skip).SetLimit(size).SetSort(bson.D{{Key: "timeStamp", Value: -1}})
	convID, err := primitive.ObjectIDFromHex(ConvID)
	if err != nil {
		return nil, err
	}

	var messages []*domain.Message
	cursor, err := messageCollection.Find(cxt, bson.D{{Key: "_conversationID", Value: convID}}, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(cxt) {
		var message *domain.Message
		err := cursor.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// method for storing message on the database
func (cr *MessageRepository) StoreMessage(cxt context.Context, message *domain.Message) (*domain.Message, error) {
	messageCollection := cr.database.Collection(cr.messageCollection)
	_, err := messageCollection.InsertOne(cxt, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// method for getting all unread message
func (cr *MessageRepository) UnReadMessage(cxt context.Context, ID string) ([]*domain.Message, error) {
	messageCollection := cr.database.Collection(cr.messageCollection)
	conversationID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	var messages []*domain.Message
	cursor, err := messageCollection.Find(cxt, bson.D{
		{Key: "_conversationID", Value: conversationID},
		{Key: "status", Value: "unread"}})

	if err != nil {
		return nil, err
	}

	for cursor.Next(cxt) {
		var message *domain.Message
		err := cursor.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil

}

// message for working w
