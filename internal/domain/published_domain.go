package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PublishedCollection = "published"
	CommentCollection   = "comments"
)

type Published struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	UserID        primitive.ObjectID `json:"userID" bson:"_userID"`
	TaskID        primitive.ObjectID `json:"taskID" bson:"_taskID"`
	Title         string             `json:"title" bson:"title"`
	Description   string             `json:"description" bson:"description"`
	Summery       string             `json:"summery" bson:"summery"`
	VideoUrl      string             `json:"videoUrl" bson:"videoUrl"`
	PictureUrl    string             `json:"pictureUrl" bson:"pictureUrl"`
	Likes         int                `json:"likes" bson:"likes"`
	Dislikes      int                `json:"dislikes" bson:"dislikes"`
	PublishedDate time.Time          `json:"publishedDate" bson:"publishedDate"`
}

type PublishedRequest struct {
	TaskID        primitive.ObjectID `json:"taskID" bson:"_taskID"`
	Title         string             `json:"title" bson:"title"`
	Description   string             `json:"description" bson:"description"`
	Summery       string             `json:"summery" bson:"summery"`
	VideoUrl      string             `json:"videoUrl" bson:"videoUrl"`
	PictureUrl    string             `json:"pictureUrl" bson:"pictureUrl"`
	PublishedDate time.Time          `json:"publishedDate" bson:"publishedDate"`
}

type UpdatePublished struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Summery     string `json:"summery" bson:"summery"`
}

type Comments struct {
	CommentID     primitive.ObjectID `json:"commentID" bson:"_commentID"`
	UserID        primitive.ObjectID `json:"userID" bson:"_userID"`
	PublishedID   primitive.ObjectID `json:"publishedID" bson:"_publishedID"`
	CommenterID   primitive.ObjectID `json:"commenterID" bson:"_commenterID"`
	Comment       string             `json:"comment" bson:"comment"`
	Likes         int                `json:"likes" bson:"likes"`
	Dislikes      int                `json:"dislikes" bson:"dislikes"`
	CommentedDate time.Time          `json:"commentedDate" bson:"commentedDate"`
}

type CommentRequest struct {
	CommenterID primitive.ObjectID `json:"commenterID" bson:"_commenterID"`
	Comment     string             `json:"comment" bson:"comment"`
}
type UpdateComment struct {
	CommentID   primitive.ObjectID `json:"commentID" bson:"_commentID"`
	PublishedID primitive.ObjectID `json:"publishedID" bson:"_publishedID"`
	Comment     string             `json:"comment" bson:"comment"`
}

type CommentInfo struct {
	CommentID   primitive.ObjectID `json:"commentID" bson:"_commentID"`
	PublishedID primitive.ObjectID `json:"publishedID" bson:"_publishedID"`
}
