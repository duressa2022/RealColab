package repository

import (
	"context"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	database   mongo.Database
	collection string
}

func NewChatRepository(database mongo.Database, collection string) *ChatRepository {
	return &ChatRepository{
		database:   database,
		collection: collection,
	}
}

// method for creating or posting chat on the database
func (cr *ChatRepository) CreateChat(cxt context.Context, message *domain.ChatMessage) (*domain.ChatMessage, error) {
	chatCollection := cr.database.Collection(cr.collection)
	_, err := chatCollection.InsertOne(cxt, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// method for getting chat by using userid and chat id
func (cr *ChatRepository) GetChatByID(cxt context.Context, UserID string, ChatID string) (*domain.ChatMessage, error) {
	chatCollection := cr.database.Collection(cr.collection)

	userid, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}
	chatid, err := primitive.ObjectIDFromHex(ChatID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_userID", Value: userid}, {Key: "_chatID", Value: chatid}}
	var response *domain.ChatMessage
	err = chatCollection.FindOne(cxt, filter).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil

}

// method for getting chats by only by using userid
func (cr *ChatRepository) GetChats(cxt context.Context, ID string) ([]*domain.ChatMessage, error) {
	var response []*domain.ChatMessage
	chatCollection := cr.database.Collection(cr.collection)

	userid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	cursor, err := chatCollection.Find(cxt, bson.D{{Key: "_userID", Value: userid}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(cxt) {
		var chat *domain.ChatMessage
		err := cursor.Decode(&chat)
		if err != nil {
			return nil, err
		}
		response = append(response, chat)
	}
	return response, nil
}

// method for deleting message from the collection if limit is reached
func (cr *ChatRepository) DeleteChats(cxt context.Context, ID string, limit int64) error {
	chatCollection := cr.database.Collection(cr.collection)
	userid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	currentChats, err := chatCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: userid}})
	if err != nil {
		return err
	}

	if currentChats <= limit {
		return nil
	}

	opts := options.Find().SetSort(bson.D{{Key: "timeStamp", Value: 1}}).SetLimit(currentChats - limit)
	cursor, err := chatCollection.Find(cxt, bson.D{{Key: "_userID", Value: userid}}, opts)
	if err != nil {
		return err
	}

	var deletedChats []*domain.ChatMessage
	for cursor.Next(cxt) {
		var deletedChat *domain.ChatMessage
		err := cursor.Decode(&deletedChat)
		if err != nil {
			return err
		}
		deletedChats = append(deletedChats, deletedChat)
	}

	for _, chat := range deletedChats {
		_, err := chatCollection.DeleteOne(cxt, bson.D{{Key: "_chatID", Value: chat.ChatID}})
		if err != nil {
			return err
		}
	}
	return nil
}
