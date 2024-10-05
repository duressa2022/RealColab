package repository

import (
	"context"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
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
