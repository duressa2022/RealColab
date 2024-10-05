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

// method for updating user information
func (ur *UserRepository) UpdateMain(cxt context.Context,userInfo *domain.UserUpdateMainInfo,UserId string)(*domain.UserUpdateMainInfo,error){
	userCollection:=ur.database.Collection(ur.collection)
	updatingInfo:=bson.M{
		"firstName":userInfo.FirstName,
		"lastName":userInfo.LastName,
		"dateOfBirth":userInfo.DateOfBirth,
		"profileUrl":userInfo.ProfileUrl,
	}

	userID,err:=primitive.ObjectIDFromHex(UserId)
	if err!=nil{
		return nil,err
	}

	updatedResult,err:=userCollection.UpdateOne(cxt,bson.D{{Key: "_userID",Value: userID}},bson.D{{Key: "$set",Value: updatingInfo}})
	if err!=nil{
		return nil,err
	}

	if updatedResult.MatchedCount==0{
		return nil,errors.New("zero document is matached")
	}
	if updatedResult.ModifiedCount==0{
		return nil,errors.New("zero document is modified")
	}

	var response *domain.UserUpdateMainInfo
	err=userCollection.FindOne(cxt,bson.D{{Key: "_userID",Value: userID}}).Decode(&response)
	if err!=nil{
		return nil,err
	}
	return response,nil

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
