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
	UserID             primitive.ObjectID       `json:"userID" bson:"_userID"`
	FirstName          string                   `json:"firstName" bson:"firstName"`
	LastName           string                   `json:"lastName" bson:"lastName"`
	Email              string                   `json:"email" bson:"email"`
	PhoneNumber        string                   `json:"phoneNumber" bson:"phoneNumber"`
	DateOfBirth        string                   `json:"dateOfBirth" bson:"dateOfBirth"`
	ProfileUrl         string                   `json:"profileUrl" bson:"profileUrl"`
	Rating             []map[string]interface{} `json:"rating" bson:"rating"`
	TaskInformation    TaskInformation          `json:"taskInformation" bson:"taskInformation"`
	TwostepVer         bool                     `json:"twoStepVer" bson:"twoStepVer"`
	NotificationChoice NotificationPreference   `json:"notificationChoice" bson:"notificationChoice"`
}
type TaskInformation struct {
	Completed  int `json:"completed" bson:"completed"`
	OnProgress int `json:"onProgress" bson:"onProgress"`
	AllTasks   int `json:"allTasks" bson:"allTasks"`
	Expired    int `json:"expired" bson:"expired"`
	Archived   int `json:"archived" bson:"archived"`
}

type UserJwtInformation struct {
	UserID primitive.ObjectID `json:"userID" bson:"_userID"`
	Email  string             `json:"email" bson:"email"`
}
type UserUpdateMainInfo struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
	ProfileUrl  string `json:"profileUrl" bson:"profileUrl"`
}
type UserSecurityInfo struct {
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	TwostepVer  bool   `json:"twoStepVer" bson:"twoStepVer"`
}
type NotificationPreference struct {
	ChooseEmail    bool `json:"chooseEmail" bson:"chooseEmail"`
	ChooseSuperApp bool `json:"chooseSuperApp" bson:"chooseSuper"`
}

type UserInterface interface {
	InsertUser(context.Context, interface{}) (interface{}, error)
}
