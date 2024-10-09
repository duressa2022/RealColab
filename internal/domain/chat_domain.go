package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ChatCollection = "chats"
)

type ChatMessage struct {
	ChatID    primitive.ObjectID `json:"chatID" bson:"_chatID"`
	UserID    primitive.ObjectID `json:"userID" bson:"_userID"`
	Response  string             `json:"response" bson:"response"`
	Prompt    string             `json:"prompt" bson:"prompt"`
	TimeStamp time.Time          `json:"timeStamp" bson:"timeStamp"`
}

type ChatRequest struct {
	Prompt string `json:"prompt" bson:"prompt"`
}
type ChatResponse struct {
	Prompt    string    `json:"prompt" bson:"prompt"`
	Response  string    `json:"response" bson:"response"`
	TimeStamp time.Time `json:"timeSatmp" bson:"timeStamp"`
}

type Message struct {
	Prompt   string `json:"prompt" bson:"prompt"`
	Response string `json:"response" bson:"response"`
}

type ChatInterface interface {
	CreateChat(context.Context, ChatMessage) (ChatMessage, error)
	GetChatByID(context.Context, string, string) (ChatMessage, error)
	GetChats(context.Context, string) ([]*ChatMessage, error)
	DeleteChats(context.Context, string, int64) error
}
