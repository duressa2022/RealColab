package usecase

import (
	"context"
	"time"
	"working/super_task/internal/domain"
	"working/super_task/internal/repository"
)

type MessageConversation struct {
	MessageConv *repository.MessageRepository
	Timeout     time.Duration
}

func NewMessageConversation(messageCon *repository.MessageRepository, timeout time.Duration) *MessageConversation {
	return &MessageConversation{
		MessageConv: messageCon,
		Timeout:     timeout,
	}
}

func (mc *MessageConversation) CreateOrRetriveConversation(cxt context.Context, participants *domain.ParticipantRequest, message domain.Message) (*domain.Conversation, error) {
	return mc.MessageConv.CreateOrUpdateConversation(cxt, participants, message)
}
func (mc *MessageConversation) StoreMessage(cxt context.Context, message *domain.Message) (*domain.Message, error) {
	return mc.MessageConv.StoreMessage(cxt, message)
}
