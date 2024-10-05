package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser string = "users"
)

type UserInformation struct {
	UserID          primitive.ObjectID       `json:"userID" bson:"_userID"`
	FirstName       string                   `json:"firstName" bson:"firstName"`
	LastName        string                   `json:"lastName" bson:"lastName"`
	Email           string                   `json:"email" bson:"email"`
	PhoneNumber     string                   `json:"phoneNumber" bson:"phoneNumber"`
	Password        string                   `json:"password" bson:"password"`
	DateOfBirth     string                   `json:"dateOfBirth" bson:"dateOfBirth"`
	ProfileUrl      string                   `json:"profileUrl" bson:"profileUrl"`
	Rating          []map[string]interface{} `json:"rating" bson:"rating"`
	TaskInformation []map[string]int         `json:"taskInformation" bson:"taskInformation"`
	TwostepVer      bool                     `json:"twoStepVer" bson:"twoStepVer"`
	CreatedAt       time.Time                `json:"created" bson:"created"`
}

type UserRegistrationRequest struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
	Password    string `json:"password" bson:"password"`
}

type UserResponse struct {
	UserID          primitive.ObjectID       `json:"userID" bson:"_userID"`
	FirstName       string                   `json:"firstName" bson:"firstName"`
	LastName        string                   `json:"lastName" bson:"lastName"`
	Email           string                   `json:"email" bson:"email"`
	PhoneNumber     string                   `json:"phoneNumber" bson:"phoneNumber"`
	DateOfBirth     string                   `json:"dateOfBirth" bson:"dateOfBirth"`
	ProfileUrl      string                   `json:"profileUrl" bson:"profileUrl"`
	Rating          []map[string]interface{} `json:"rating" bson:"rating"`
	TaskInformation []map[string]int         `json:"taskInformation" bson:"taskInformation"`
	TwostepVer      bool                     `json:"twoStepVer" bson:"twoStepVer"`
}

type UserJwtInformation struct {
	UserID primitive.ObjectID `json:"userID" bson:"_userID"`
	Email  string             `json:"email" bson:"email"`
}

type UserInterface interface {
	InsertUser(context.Context, interface{}) (interface{}, error)
}
