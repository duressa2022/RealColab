package usecase

import (
	"context"
	"time"
	"working/super_task/internal/domain"
	"working/super_task/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupUseCase struct {
	GroupRepos *repository.GroupRepos
	Timeout    time.Duration
}

func NewGroupUseCase(group *repository.GroupRepos, time time.Duration) *GroupUseCase {
	return &GroupUseCase{
		GroupRepos: group,
		Timeout:    time,
	}
}

// method for creating group based on the information
func (gr *GroupUseCase) CreateGroup(cxt context.Context, group *domain.GroupRequest, ID string) (*domain.Group, error) {
	creatorID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	var createdGroup domain.Group
	createdGroup.CreatorID = creatorID
	createdGroup.GroupName = group.GroupName
	createdGroup.GroupInformation = group.GroupInformation
	createdGroup.GroupType = group.GroupType
	createdGroup.PictureUrl = group.PictureUrl
	createdGroup.CreatedAt = time.Now()
	return gr.GroupRepos.StoreGroup(cxt, &createdGroup)
}

// method for updating group information
func (gr *GroupUseCase) UpdateGroup(cxt context.Context, group *domain.GroupRequest, ID string) (*domain.Group, error) {
	return gr.GroupRepos.UpdateGroup(cxt, group, ID)
}

// method for deleting group by using groupid
func (gr *GroupUseCase) DeleteGroup(cxt context.Context, ID string) error {
	return gr.GroupRepos.DeleteGroup(cxt, ID)
}

// method for getting group information by using id
func (gr *GroupUseCase) GetGroupInformation(cxt context.Context, Id string) (*domain.Group, error) {
	return gr.GroupRepos.GetGroupInformation(cxt, Id)
}

// method getting adding member into the group
func (gr *GroupUseCase) AddMember(cxt context.Context, ID string, user *domain.UserInGroup) error {
	return gr.GroupRepos.AddMember(cxt, ID, user)
}

// method for deleting member from the group
func (gr *GroupUseCase) DeleteMember(cxt context.Context, UserID string, GroupID string) error {
	return gr.GroupRepos.DeleteMember(cxt, UserID, GroupID)
}

// method for promoting user into the admin
func (gr *GroupUseCase) PromoteUser(cxt context.Context, user *domain.UserInGroup, GroupID string) error {
	return gr.GroupRepos.AddAdmin(cxt, GroupID, user)
}

// method for demoting admin from the group
func (gr *GroupUseCase) DemoteUser(cxt context.Context, UserId string, GroupID string) error {
	return gr.GroupRepos.DeleteAdmin(cxt, UserId, GroupID)
}

// method for blocking user from the group
func (gr *GroupUseCase) BlockUser(cxt context.Context, UserId string, GroupID string, reason string) error {
	return gr.GroupRepos.BlockUser(cxt, UserId, GroupID, reason)
}

// method for unblocking user from the group
func (gr *GroupUseCase) UnblockUser(cxt context.Context, UserID string, GroupID string) error {
	return gr.GroupRepos.UnBlockUser(cxt, UserID, GroupID)
}

// method for getting group message
func (gr *GroupUseCase) GetMessages(cxt context.Context, groupID string, size int64, page int64) ([]*domain.Message, int64, error) {
	return gr.GroupRepos.GetMessages(cxt, groupID, size, page)
}

// method for getting all members
func (gr *GroupUseCase) GetAllMembers(cxt context.Context, groupID string) ([]map[string]interface{}, error) {
	return gr.GroupRepos.GetAllMembers(cxt, groupID)
}

// method for creating or retriving group conversation
func (gr *GroupUseCase) CreateOrRetriveConversation(cxt context.Context, message domain.GroupMessage) (*domain.GroupConversation, error) {
	return gr.GroupRepos.CreateOrUpdateConversation(cxt, message)
}
func (gr *GroupUseCase) StoreMessage(cxt context.Context, message *domain.GroupMessage) (*domain.GroupMessage, error) {
	return gr.GroupRepos.StoreMessage(cxt, message)
}
