package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TaskCollection="tasks"
)

type Task struct {
	TaskID      primitive.ObjectID `json:"taskID" bson:"_taskID"`
	UserID      primitive.ObjectID `json:"userID" bson:"_userID"`
	Title       string             `json:"title" bson:"title"`
	StartDate   time.Time          `json:"startDate" bson:"startDate"`
	DueDate     time.Time          `json:"dueDate" bson:"dueDate"`
	TimePerDay  int64              `json:"timePerDay" bson:"timePerDay"`
	Description string             `json:"description" bson:"description" `
	Status      string             `json:"status" bson:"status"`
	Type        string             `json:"type" bson:"type"`
	Members     []*UserInformation `json:"members" bson:"members" `
}

type EditTask struct{
	Title       string             `json:"title" bson:"title"`
	TimePerDay  int64              `json:"timePerDay" bson:"timePerDay"`
	Description string             `json:"description" bson:"description" `
	Type        string             `json:"type" bson:"type"`
}

