package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TaskCollection = "tasks"
)

type Task struct {
	TaskID      primitive.ObjectID   `json:"taskID" bson:"_taskID"`
	UserID      primitive.ObjectID   `json:"userID" bson:"_userID"`
	Title       string               `json:"title" bson:"title"`
	StartDate   string               `json:"startDate" bson:"startDate"`
	DueDate     string               `json:"dueDate" bson:"dueDate"`
	TimePerDay  string               `json:"timePerDay" bson:"timePerDay"`
	Description string               `json:"description" bson:"description" `
	Status      string               `json:"status" bson:"status"`
	Type        string               `json:"type" bson:"type"`
	Members     []*MemberInformation `json:"members" bson:"members" `
}

type TaskRequest struct {
	Title       string `json:"title" bson:"title"`
	StartDate   string `json:"startDate" bson:"startDate"`
	DueDate     string `json:"dueDate" bson:"dueDate"`
	TimePerDay  string `json:"timePerDay" bson:"timePerDay"`
	Description string `json:"description" bson:"description" `
	Status      string `json:"status" bson:"status"`
	Type        string `json:"type" bson:"type"`
}

type PrivateTask struct {
	TaskID      primitive.ObjectID `json:"taskID" bson:"_taskID"`
	UserID      primitive.ObjectID `json:"userID" bson:"_userID"`
	Title       string             `json:"title" bson:"title"`
	StartDate   string             `json:"startDate" bson:"startDate"`
	DueDate     string             `json:"dueDate" bson:"dueDate"`
	TimePerDay  string             `json:"timePerDay" bson:"timePerDay"`
	Description string             `json:"description" bson:"description" `
	Status      string             `json:"status" bson:"status"`
}

type SharedTask struct {
	TaskID      primitive.ObjectID   `json:"taskID" bson:"_taskID"`
	UserID      primitive.ObjectID   `json:"userID" bson:"_userID"`
	Title       string               `json:"title" bson:"title"`
	StartDate   string               `json:"startDate" bson:"startDate"`
	DueDate     string               `json:"dueDate" bson:"dueDate"`
	TimePerDay  string               `json:"timePerDay" bson:"timePerDay"`
	Description string               `json:"description" bson:"description" `
	Status      string               `json:"status" bson:"status"`
	Members     []primitive.ObjectID `json:"members" bson:"members" `
}

type EditTask struct {
	Title       string `json:"title" bson:"title"`
	TimePerDay  string `json:"timePerDay" bson:"timePerDay"`
	Description string `json:"description" bson:"description" `
	Type        string `json:"type" bson:"type"`
}
type MemberInformation struct {
	MemberID   primitive.ObjectID `json:"userID" bson:"_userID"`
	MemberRole string             `json:"memberRole" bson:"memberRole"`
}

type SearchTerm struct {
	SearchTerm string `json:"searchTerm" bson:"searchTerm"`
}
