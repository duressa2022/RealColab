package usecase

import (
	"context"
	"fmt"
	"time"
	"working/super_task/internal/domain"
	"working/super_task/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUseCase struct {
	ChatRepository *repository.ChatRepository
	Timeout        time.Duration
}

func NewChatUseCase(chat *repository.ChatRepository, timeout time.Duration) *ChatUseCase {
	return &ChatUseCase{
		ChatRepository: chat,
		Timeout:        timeout,
	}
}

// method for storing chat into the database
func (cu *ChatUseCase) StoreChat(cxt context.Context, UserID string, prompt *domain.ChatRequest, content interface{}, summery interface{}) (*domain.ChatMessage, error) {
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}
	chatID := primitive.NewObjectID()

	sessionID, err := primitive.ObjectIDFromHex(prompt.SessionID)
	if err != nil {
		return nil, err
	}

	var chatMessage domain.ChatMessage
	chatMessage.ChatID = chatID
	chatMessage.UserID = userID
	chatMessage.SessionID = sessionID
	chatMessage.Prompt = prompt.Prompt
	chatMessage.Response = content
	chatMessage.ChatSummery = summery
	chatMessage.TimeStamp = time.Now()

	return cu.ChatRepository.StoreChatForSession(cxt, &chatMessage)

}

// method for making prompt of the database
func (cu *ChatUseCase) CreatePrompt(cxt context.Context, UserID string, SessionID string) string {
	messages, err := cu.ChatRepository.GetChatsForSession(cxt, UserID, SessionID)

	prompt := `Assume you are a digital assistance that is works on advicing user:
	            :you are going to suggest different kinds of tasks based on user request
				:you are going to advice the how to do different kinds works or tasks
				:you are going make plan for different tasks based on the user history.
				:make your response in html format with proper htmltags(ul,li,p,br)
				:the chat history is given to you as prompt and response format
				:you are only to answer information related about techinical tasks
				:create tasks/ or projects and secheduling events 
				:first what types of the pojects that the user wants to work
				:then give for user suggestions,plans and brainstorming,
				:use this chat history to answer questions and use it as memory `
	if err != nil {
		return prompt
	}
	for _, message := range messages {
		prompt += fmt.Sprintf("prompt: %s \n response: %s \n", message.Prompt, message.Response)

	}
	return prompt

}

// method for fetching sessions for user
func (cu *ChatUseCase) FetchSession(cxt context.Context, ID string, requestNumber int64, limit int64) ([]*domain.SessionHistory, int64, error) {
	return cu.ChatRepository.GetAllSessions(cxt, ID, requestNumber, limit)
}

// method for fetching chat message for given session
func (cu *ChatUseCase) FetchSessionChat(cxt context.Context, userID string, sessionID string) ([]*domain.ChatResponse, error) {
	return cu.ChatRepository.GetChatsForSession(cxt, userID, sessionID)
}

// method for deleting session from the system
func (cu *ChatUseCase) DeleteSession(cxt context.Context, userID string, sessionID string) error {
	return cu.ChatRepository.DeleteSession(cxt, userID, sessionID)
}

// method fro creating new session
func (cu *ChatUseCase) CreateSession(cxt context.Context, UserID string) (*domain.SessionHistory, error) {
	var session domain.SessionHistory
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}
	session.UserID = userID
	session.SessionID = primitive.NewObjectID()
	session.ChatSummery = "NewChat"
	session.Activated = false
	session.TimeStamp = time.Now()
	return cu.ChatRepository.CreateNewSession(cxt, &session)
}
