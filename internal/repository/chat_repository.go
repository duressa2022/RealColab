package repository

import (
	"context"
	"errors"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	database          mongo.Database
	chatCollection    string
	sessionCollection string
}

const (
	sessionlimit = 100
	chatLimit    = 100
)

func NewChatRepository(database mongo.Database, chat string, session string) *ChatRepository {
	return &ChatRepository{
		database:          database,
		chatCollection:    chat,
		sessionCollection: session,
	}
}

// method for creating or posting chat on the database
func (cr *ChatRepository) CreateChat(cxt context.Context, message *domain.ChatMessage) (*domain.ChatMessage, error) {
	chatCollection := cr.database.Collection(cr.chatCollection)
	_, err := chatCollection.InsertOne(cxt, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// method for creating new sessions for chatting
func (cr *ChatRepository) CreateNewSession(cxt context.Context, session *domain.SessionHistory) (*domain.SessionHistory, error) {
	sessionCollection := cr.database.Collection(cr.sessionCollection)
	currentSession, err := sessionCollection.CountDocuments(cxt, bson.D{{Key: "_userID", Value: session.UserID}})
	if err != nil {
		return nil, err
	}

	if currentSession >= sessionlimit {
		return nil, errors.New("memory is full")
	}

	_, err = sessionCollection.InsertOne(cxt, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// method for getting all sessions that is created by user
func (cr *ChatRepository) GetAllSessions(cxt context.Context, UserID string, requestNumber int64, limit int64) ([]*domain.SessionHistory, int64, error) {
	sessionCollection := cr.database.Collection(cr.sessionCollection)

	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, 0, err
	}

	sessionSkip := (requestNumber - 1) * limit
	opts := options.Find().SetSkip(sessionSkip).SetLimit(limit).SetSort(bson.D{{Key: "timeStamp", Value: -1}})
	filter := bson.D{{Key: "_userID", Value: userID}}

	cursor, err := sessionCollection.Find(cxt, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(cxt)

	var sessions []*domain.SessionHistory
	for cursor.Next(cxt) {
		var session *domain.SessionHistory
		err := cursor.Decode(&session)
		if err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, session)
	}

	totalSessions, err := sessionCollection.CountDocuments(cxt, filter)
	if err != nil {
		return nil, 0, err
	}
	return sessions, totalSessions, nil

}

// method for deleting session by using ID
func (cr *ChatRepository) DeleteSession(cxt context.Context, UserID string, SessionID string) error {
	chatCollection := cr.database.Collection(cr.chatCollection)
	sessionCollection := cr.database.Collection(cr.sessionCollection)

	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return err
	}
	sessionID, err := primitive.ObjectIDFromHex(SessionID)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_userID", Value: userID}, {Key: "_sessionID", Value: sessionID}}

	_, err = sessionCollection.DeleteOne(cxt, filter)
	if err != nil {
		return err
	}

	_, err = chatCollection.DeleteMany(cxt, filter)
	if err != nil {
		return err
	}
	return nil

}

// method for getting chats for given session and user
func (cr *ChatRepository) GetChatsForSession(cxt context.Context, UserID string, SessionID string) ([]*domain.ChatResponse, error) {
	var response []*domain.ChatResponse
	chatCollection := cr.database.Collection(cr.chatCollection)

	userid, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}

	sessionID, err := primitive.ObjectIDFromHex(SessionID)
	if err != nil {
		return nil, err
	}

	cursor, err := chatCollection.Find(cxt, bson.D{{Key: "_userID", Value: userid}, {Key: "_sessionID", Value: sessionID}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(cxt) {
		var chat *domain.ChatResponse
		err := cursor.Decode(&chat)
		if err != nil {
			return nil, err
		}
		response = append(response, chat)
	}

	return response, nil
}

// method for storing messages/ or chats
func (cr *ChatRepository) StoreChatForSession(cxt context.Context, chat *domain.ChatMessage) (*domain.ChatMessage, error) {
	chatCollection := cr.database.Collection(cr.chatCollection)
	sessionCollection := cr.database.Collection(cr.sessionCollection)

	userID := chat.UserID
	sessionID := chat.SessionID
	filter := bson.D{{Key: "_userID", Value: userID}, {Key: "_sessionID", Value: sessionID}}

	var session *domain.SessionHistory
	err := sessionCollection.FindOne(cxt, filter).Decode(&session)
	if err != nil {
		return nil, err
	}

	if !session.Activated {
		updateSession := bson.M{
			"chatSummery": chat.ChatSummery,
			"activated":   true,
		}

		_, err := sessionCollection.UpdateOne(cxt, filter, bson.D{{Key: "$set", Value: updateSession}})
		if err != nil {
			return nil, err
		}
	}

	_, err = chatCollection.InsertOne(cxt, chat)
	if err != nil {
		return nil, err
	}
	return chat, nil

}
