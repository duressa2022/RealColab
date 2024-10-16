package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	GroupID          primitive.ObjectID            `json:"groupID" bson:"_groupID"`
	CreatorID        primitive.ObjectID            `json:"creatorID" bson:"_creatorID"`
	GroupName        string                        `json:"groupName" bson:"groupName"`
	GroupAdmin       []*UserInGroup                `json:"groupAdmin"  bson:"groupAdmin"`
	GroupMembers     []*UserInGroup                `json:"groupMember" bson:"groupMember"`
	GroupType        string                        `json:"groupType" bson:"groupType"`
	PictureUrl       string                        `json:"pictureUrl" bson:"pictureUrl"`
	GroupInformation string                        `json:"groupInformation" bson:"groupInformation"`
	SharedFile       []string                      `json:"sharedFile" bson:"sharedFile"`
	Messages         []*Message                    `json:"messages" bson:"messages"`
	BlockedUser      map[primitive.ObjectID]string `json:"blockedUser" bson:"blockedUser"`
	CreatedAt        time.Time                     `json:"createdAt" bson:"createdAt"`
}

type GroupRequest struct {
	GroupName        string `json:"groupName" bson:"groupName"`
	GroupType        string `json:"groupType" bson:"groupType"`
	PictureUrl       string `json:"pictureUrl" bson:"pictureUrl"`
	GroupInformation string `json:"groupInformation" bson:"groupInformation"`
}
