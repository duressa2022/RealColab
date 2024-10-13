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
func (cu *ChatUseCase) StoreChat(cxt context.Context, message *domain.RequestResponse, Id string) (*domain.ChatMessage, error) {
	userID, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return nil, err
	}
	chatID := primitive.NewObjectID()

	var chatMessage domain.ChatMessage
	chatMessage.ChatID = chatID
	chatMessage.UserID = userID
	chatMessage.Prompt = message.Prompt
	chatMessage.Response = message.Response
	chatMessage.TimeStamp = time.Now()

	return cu.ChatRepository.CreateChat(cxt, &chatMessage)

}

// method for making prompt of the database
func (cu *ChatUseCase) CreatePrompt(cxt context.Context, ID string) string {
	messages, err := cu.ChatRepository.GetChats(cxt, ID)
	prompt := `Assume you are a digital assistance that is works on advicing user:
	            :you are going to suggest different kinds of tasks based on user request
				:you are going to advice the how to do different kinds works or tasks
				:you are going make plan for different `
	if err != nil {
		return prompt
	}
	for _, message := range messages {
		prompt += fmt.Sprintf("prompt: %s \n response: %s \n", message.Prompt, message.Response)

	}
	return prompt
}
