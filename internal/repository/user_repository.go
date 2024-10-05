package repository

import (
	"context"
	"errors"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(collection string, database mongo.Database) *UserRepository {
	return &UserRepository{
		database:   database,
		collection: collection,
	}
}

// method for getting notification choice
func (ur *UserRepository) GetNotificationChoice(cxt context.Context, Id string) (*domain.NotificationPreference, error) {
	userCollection := ur.database.Collection(ur.collection)
	UserId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return nil, err
	}

	var notChoice *domain.NotificationPreference
	err = userCollection.FindOne(cxt, bson.D{{Key: "_userID", Value: UserId}}).Decode(&notChoice)
	if err != nil {
		return nil, err
	}
	return notChoice, nil
}

// method for updating notfication choice
func (ur *UserRepository) UpdateNotificationChoice(cxt context.Context, change *domain.NotificationPreference, Id string) (*domain.NotificationPreference, error) {
	userCollection := ur.database.Collection(ur.collection)
	userid, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return nil, err
	}

	updating := bson.M{
		"chooseEmail":    change.ChooseEmail,
		"chooseSuperApp": change.ChooseSuperApp,
	}

	updatedResult,err:=userCollection.UpdateOne(cxt,bson.D{{Key: "_userID",Value: userid}},bson.D{{Key: "$set",Value: updating}})
	if err!=nil{
		return nil,err
	}

	if updatedResult.ModifiedCount==0{
		return nil,errors.New("no modified data")
	}
	if updatedResult.MatchedCount==0{
		return nil,errors.New("no matched docs")
	}

	return change,err
}

// method for updating user password
func (ur *UserRepository) UpdatePassword(cxt context.Context, Pass *domain.ChangePassword, Id string) error {
	userCollection := ur.database.Collection(ur.collection)
	UserId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return err
	}

	updatingInfo := bson.M{
		"password": Pass.NewPassword,
	}

	updateResult, err := userCollection.UpdateOne(cxt, bson.D{{Key: "_userID", Value: UserId}}, bson.D{{Key: "$set", Value: updatingInfo}})
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return errors.New("no matched docs found")
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no modified docs found")
	}
	return nil
}

// method for getting user security information
func (ur *UserRepository) GetSecurityInfo(cxt context.Context, Id string) (*domain.UserSecurityInfo, error) {
	userCollection := ur.database.Collection(ur.collection)
	userId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return nil, err
	}

	var securityInfo *domain.UserSecurityInfo
	err = userCollection.FindOne(cxt, bson.D{{Key: "_userID", Value: userId}}).Decode(&securityInfo)
	if err != nil {
		return nil, err
	}

	return securityInfo, err
}

// method for updating user information
func (ur *UserRepository) UpdateMain(cxt context.Context, userInfo *domain.UserUpdateMainInfo, UserId string) (*domain.UserUpdateMainInfo, error) {
	userCollection := ur.database.Collection(ur.collection)
	updatingInfo := bson.M{
		"firstName":   userInfo.FirstName,
		"lastName":    userInfo.LastName,
		"dateOfBirth": userInfo.DateOfBirth,
		"profileUrl":  userInfo.ProfileUrl,
	}

	userID, err := primitive.ObjectIDFromHex(UserId)
	if err != nil {
		return nil, err
	}

	updatedResult, err := userCollection.UpdateOne(cxt, bson.D{{Key: "_userID", Value: userID}}, bson.D{{Key: "$set", Value: updatingInfo}})
	if err != nil {
		return nil, err
	}

	if updatedResult.MatchedCount == 0 {
		return nil, errors.New("zero document is matached")
	}
	if updatedResult.ModifiedCount == 0 {
		return nil, errors.New("zero document is modified")
	}

	var response *domain.UserUpdateMainInfo
	err = userCollection.FindOne(cxt, bson.D{{Key: "_userID", Value: userID}}).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil

}

// method for inserting user information into database
func (ur *UserRepository) InsertUser(cxt context.Context, user *domain.UserInformation) (*domain.UserResponse, error) {
	userCollection := ur.database.Collection(ur.collection)
	_, err := userCollection.InsertOne(cxt, user)
	if err != nil {
		return nil, err
	}

	var Response *domain.UserResponse
	err = userCollection.FindOne(cxt, bson.D{{Key: "_userID", Value: user.UserID}}).Decode(&Response)
	if err != nil {
		return nil, err
	}

	return Response, nil
}

// method for user information by using email
func (ur *UserRepository) GetUserByEmail(cxt context.Context, email string) (*domain.UserResponse, error) {
	userCollection := ur.database.Collection(ur.collection)
	var userReponse *domain.UserResponse
	err := userCollection.FindOne(cxt, bson.D{{Key: "email", Value: email}}).Decode(&userReponse)
	if err != nil {
		return nil, err
	}
	return userReponse, nil

}

// method for getting user by using email
func (ur *UserRepository) GetUserByEmailLogin(cxt context.Context, email string) (*domain.UserInformation, error) {
	userCollection := ur.database.Collection(ur.collection)
	var userReponse *domain.UserInformation
	err := userCollection.FindOne(cxt, bson.D{{Key: "email", Value: email}}).Decode(&userReponse)
	if err != nil {
		return nil, err
	}
	return userReponse, nil
}

// method for getting user by using phone
func (ur *UserRepository) GetUserByPhone(cxt context.Context, phone string) (*domain.UserInformation, error) {
	userCollection := ur.database.Collection(ur.collection)
	var userReponse *domain.UserInformation
	err := userCollection.FindOne(cxt, bson.D{{Key: "phoneNumber", Value: phone}}).Decode(&userReponse)
	if err != nil {
		return nil, err
	}
	return userReponse, nil
}
