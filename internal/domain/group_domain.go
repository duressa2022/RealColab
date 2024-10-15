package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	GroupID          primitive.ObjectID `json:"groupID" bson:"_groupID"`
	GroupName        string             `json:"groupName" bson:"groupName"`
	GroupAdmin       []UserResponse     `json:"groupAdmin"  bson:"groupAdmin"`
	GroupMembers     []UserResponse     `json:"groupMember" bson:"groupMember"`
	GroupType        string             `json:"groupType" bson:"groupType"`
	PictureUrl       string             `json:"pictureUrl" bson:"pictureUrl"`
	GroupInformation string             `json:"groupInformation" bson:"groupInformation"`
	SharedFile       []string           `json:"sharedFile" bson:"sharedFile"`
	Messages         []Message          `json:"messages" bson:"messages"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
}

type GroupRequest struct {
	GroupName        string `json:"groupName" bson:"groupName"`
	GroupType        string `json:"groupType" bson:"groupType"`
	PictureUrl       string `json:"pictureUrl" bson:"pictureUrl"`
	GroupInformation string `json:"groupInformation" bson:"groupInformation"`
}
