package repository

import (
	"context"
	"errors"
	"time"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupRepos struct {
	database          mongo.Database
	collection        string
	userCollection    string
	conCollection     string
	messageCollection string
}

func NewGroupRepos(db mongo.Database, collection string, userCollection string, conv string, message string) *GroupRepos {
	return &GroupRepos{
		database:          db,
		collection:        collection,
		userCollection:    userCollection,
		conCollection:     conv,
		messageCollection: message,
	}
}

// method for storing message on the database
func (gr *GroupRepos) StoreMessage(cxt context.Context, message *domain.GroupMessage) (*domain.GroupMessage, error) {
	messageCollection := gr.database.Collection(gr.messageCollection)
	_, err := messageCollection.InsertOne(cxt, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// method for creating new conversation or retiving existing one
func (gr *GroupRepos) CreateOrUpdateConversation(cxt context.Context, message domain.GroupMessage) (*domain.GroupConversation, error) {
	convCollection := gr.database.Collection(gr.conCollection)

	var convModel *domain.GroupConversation
	err := convCollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: message.GroupID}}).Decode(&convModel)
	if err != nil {
		newConvMode := &domain.GroupConversation{
			ConversationID: primitive.NewObjectID(),
			LastMessage:    message,
			CreatedAt:      primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt:      primitive.NewDateTimeFromTime(time.Now()),
		}
		convModel = newConvMode
	} else {
		convModel.LastMessage = message
		convModel.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	}
	return convModel, nil

}

// method for storing data on the database
func (gr *GroupRepos) StoreGroup(cxt context.Context, group *domain.Group) (*domain.Group, error) {
	groupCollection := gr.database.Collection(gr.collection)
	_, err := groupCollection.InsertOne(cxt, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// method for updating group information
func (gr *GroupRepos) UpdateGroup(cxt context.Context, group *domain.GroupRequest, ID string) (*domain.Group, error) {
	groupcollection := gr.database.Collection(gr.collection)
	updating := bson.M{
		"groupName":        group.GroupName,
		"groupType":        group.GroupType,
		"groupInformation": group.GroupInformation,
		"pictureUrl":       group.PictureUrl,
	}

	groupID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	_, err = groupcollection.UpdateOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}, updating)
	if err != nil {
		return nil, err
	}

	var updatedGroup *domain.Group
	err = groupcollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(&updatedGroup)
	if err != nil {
		return nil, err
	}
	return updatedGroup, nil

}

// method for deleting the group by using id
func (gr *GroupRepos) DeleteGroup(cxt context.Context, ID string) error {
	groupCollection := gr.database.Collection(gr.collection)
	groupID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	_, err = groupCollection.DeleteOne(cxt, bson.D{{Key: "groupID", Value: groupID}})
	return err

}

// method for adding member into the the group
func (gr *GroupRepos) AddMember(cxt context.Context, Id string, user *domain.UserInGroup) error {
	groupCollection := gr.database.Collection(gr.collection)
	groupID, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return err
	}

	var group *domain.Group
	err = groupCollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(&group)
	if err != nil {
		return err
	}

	updating := bson.M{
		"groupMembers": append(group.GroupMembers, user),
	}

	group.GroupMembers = append(group.GroupMembers, user)
	_, err = groupCollection.UpdateOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}, updating)
	if err != nil {
		return err
	}
	return err
}

// method for getting all members
func (gr *GroupRepos) GetAllMembers(cxt context.Context, ID string) ([]map[string]interface{}, error) {
	groupCollection := gr.database.Collection(gr.collection)
	userCollection := gr.database.Collection(gr.userCollection)
	groupID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	var group *domain.Group
	err = groupCollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(&group)
	if err != nil {
		return nil, err
	}

	var members []*domain.UserInGroup
	members = append(members, group.GroupMembers...)

	var membersInfos []map[string]interface{}
	for _, member := range members {
		var user *domain.UserResponse
		err := userCollection.FindOne(cxt, bson.D{{Key: "userID", Value: member.UserID}}).Decode(&user)
		if err != nil {
			return nil, err
		}
		userinfo := map[string]interface{}{
			"firstName":  user.FirstName,
			"lastName":   user.LastName,
			"profileUrl": user.ProfileUrl,
			"userID":     user.UserID,
		}
		membersInfos = append(membersInfos, userinfo)

	}
	return membersInfos, nil
}

// method for deleting a member from the group
func (gr *GroupRepos) DeleteMember(cxt context.Context, UserID string, GroupID string) error {
	groupCollection := gr.database.Collection(gr.collection)

	userId, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return err
	}

	groupID, err := primitive.ObjectIDFromHex(GroupID)
	if err != nil {
		return err
	}

	var group domain.Group
	err = groupCollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(group)
	if err != nil {
		return err
	}

	userIndex := -1
	for index, member := range group.GroupMembers {
		if member.UserID == userId {
			userIndex = index
			break
		}
	}

	if userIndex == -1 {
		return errors.New("user not found within the database")
	}

	var newMembers []*domain.UserInGroup
	newMembers = append(newMembers, group.GroupMembers[:userIndex]...)
	newMembers = append(newMembers, group.GroupMembers[userIndex+1:]...)

	updatingMember := bson.M{
		"groupMembers": newMembers,
	}
	updateResult, err := groupCollection.UpdateOne(cxt, bson.D{{Key: "_groupID", Value: group}}, updatingMember)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no modified docs")
	}
	if updateResult.MatchedCount == 0 {
		return errors.New("no new matched docs")
	}
	return nil

}

// method for getting all group information
func (gr *GroupRepos) GetGroupInformation(cxt context.Context, ID string) (*domain.Group, error) {
	groupCollection := gr.database.Collection(gr.collection)
	groupID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	var group *domain.Group
	err = groupCollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return group, nil

}

// method for adding admin into the goup
func (gr *GroupRepos) AddAdmin(cxt context.Context, Id string, user *domain.UserInGroup) error {
	groupCollection := gr.database.Collection(gr.collection)
	groupID, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return err
	}

	var group *domain.Group
	err = groupCollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(&group)
	if err != nil {
		return err
	}

	updating := bson.M{
		"groupAdmin": append(group.GroupAdmin, user),
	}

	group.GroupMembers = append(group.GroupMembers, user)
	_, err = groupCollection.UpdateOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}, updating)
	if err != nil {
		return err
	}
	return err
}

// method for removing admin from the group
func (gr *GroupRepos) DeleteAdmin(cxt context.Context, UserID string, GroupID string) error {
	groupCollection := gr.database.Collection(gr.collection)

	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return err
	}

	groupID, err := primitive.ObjectIDFromHex(GroupID)
	if err != nil {
		return err
	}

	var group domain.Group
	err = groupCollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(&group)
	if err != nil {
		return err
	}

	adminIndex := -1
	for index, admin := range group.GroupAdmin {
		if admin.UserID == userID {
			adminIndex = index
			break
		}
	}

	if adminIndex == -1 {
		return errors.New("no admin by given id")
	}

	var admins []*domain.UserInGroup
	admins = append(admins, group.GroupAdmin[:adminIndex]...)
	admins = append(admins, group.GroupAdmin[adminIndex+1:]...)

	updatingAdmin := bson.M{
		"groupAdmins": admins,
	}

	updateResult, err := groupCollection.UpdateOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}, updatingAdmin)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount == 0 {
		return errors.New("no matched docs")
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no modified docs")
	}

	return nil

}

// method for blocking user from the group
func (gr *GroupRepos) BlockUser(cxt context.Context, UserID string, GroupID string, reason string) error {
	groupcollection := gr.database.Collection(gr.collection)

	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return err
	}

	groupID, err := primitive.ObjectIDFromHex(GroupID)
	if err != nil {
		return err
	}

	var group *domain.Group
	err = groupcollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(&group)
	if err != nil {
		return err
	}

	blockedUsers := group.BlockedUser
	blockedUsers[userID] = reason

	updatingUsers := bson.M{
		"blockedUsers": blockedUsers,
	}

	updateResult, err := groupcollection.UpdateOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}, updatingUsers)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return errors.New("no matched docs are found")
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no modfied docs are found")

	}

	return nil
}

// method for unblocking user from the group
func (gr *GroupRepos) UnBlockUser(cxt context.Context, UserID string, GroupID string) error {
	groupCollection := gr.database.Collection(gr.collection)

	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return err
	}

	groupID, err := primitive.ObjectIDFromHex(GroupID)
	if err != nil {
		return err
	}

	var group domain.Group
	err = groupCollection.FindOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}).Decode(&group)
	if err != nil {
		return err
	}

	delete(group.BlockedUser, userID)
	updatingResult := bson.M{
		"blockedUser": group.BlockedUser,
	}
	updateResult, err := groupCollection.UpdateOne(cxt, bson.D{{Key: "_groupID", Value: groupID}}, updatingResult)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no modified docs is found")
	}
	if updateResult.MatchedCount == 0 {
		return errors.New("no matched docs are found")
	}
	return nil
}

// method for handling user leaving the group
func (gr *GroupRepos) GetMessages(cxt context.Context, GroupID string, size int64, page int64) ([]*domain.Message, int64, error) {
	groupCollection := gr.database.Collection(gr.collection)

	groupID, err := primitive.ObjectIDFromHex(GroupID)
	if err != nil {
		return nil, 0, err
	}
	skip := (size - 1) / page
	opts := options.Find().SetSkip(skip).SetLimit(size).SetSort(bson.D{{Key: "timeStamp", Value: -1}})

	var groupMessages []*domain.Message
	cursor, err := groupCollection.Find(cxt, bson.D{{Key: "_groupID", Value: groupID}}, opts)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(cxt)

	for cursor.Next(cxt) {
		var message *domain.Message
		err := cursor.Decode(&cursor)
		if err != nil {
			return nil, 0, err
		}
		groupMessages = append(groupMessages, message)
	}

	totalMessages, err := groupCollection.CountDocuments(cxt, bson.D{{Key: "_groupID", Value: groupID}})
	if err != nil {
		return nil, 0, err
	}

	return groupMessages, totalMessages, nil

}
