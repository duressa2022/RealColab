package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ChatCollection    = "chats"
	SessionCollection = "sessions"
)

type ChatMessage struct {
	ChatID      primitive.ObjectID `json:"chatID" bson:"_chatID"`
	UserID      primitive.ObjectID `json:"userID" bson:"_userID"`
	SessionID   primitive.ObjectID `json:"sessionID" bson:"_sessionID"`
	ChatSummery interface{}        `json:"chatSummery" bson:"chatSummery"`
	Response    interface{}        `json:"response" bson:"response"`
	Prompt      string             `json:"prompt" bson:"prompt"`
	TimeStamp   time.Time          `json:"timeStamp" bson:"timeStamp"`
}

type SessionHistory struct {
	UserID      primitive.ObjectID `json:"userID" bson:"_userID"`
	SessionID   primitive.ObjectID `json:"sessionID" bson:"_sessionID"`
	ChatSummery interface{}        `json:"chatSummery" bson:"chatSummery"`
	Activated   bool               `json:"activated" bson:"activated"`
	TimeStamp   time.Time          `json:"timeStamp" bson:"timeStamp"`
}

type ChatRequest struct {
	Prompt    string `json:"prompt" bson:"prompt"`
	SessionID string `json:"sessionID" bson:"_sessionID"`
}
type ChatResponse struct {
	Prompt      string      `json:"prompt" bson:"prompt"`
	Response    interface{} `json:"response" bson:"response"`
	ChatSummery interface{} `json:"chatSummery" bson:"chatSummery"`
	TimeStamp   time.Time   `json:"timeSatmp" bson:"timeStamp"`
}

type RequestResponse struct {
	Prompt   string      `json:"prompt" bson:"prompt"`
	Response interface{} `json:"response" bson:"response"`
}

type ChatInterface interface {
	CreateChat(context.Context, ChatMessage) (ChatMessage, error)
	GetChatByID(context.Context, string, string) (ChatMessage, error)
	GetChats(context.Context, string) ([]*ChatMessage, error)
	DeleteChats(context.Context, string, int64) error
}
