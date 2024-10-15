package repository

import (
	"context"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupRepos struct {
	database   mongo.Database
	collection string
}

func NewGroupRepos(db mongo.Database, collection string) *GroupRepos {
	return &GroupRepos{
		database:   db,
		collection: collection,
	}
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
func (gr *GroupRepos) AddMember(cxt context.Context, Id string, user domain.UserResponse) error {
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
func (gr *GroupRepos) GetAllMembers(cxt context.Context, ID string) ([]*domain.UserResponse, error) {
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

	var members []*domain.UserResponse
	for _, member := range group.GroupMembers {
		members = append(members, &member)
	}
	return members, nil
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

// method fpr adding admin into the goup
func (gr *GroupRepos) AddAdmin(cxt context.Context, Id string, user domain.UserResponse) error {
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
