package usecase

import (
	"context"
	"time"
	"working/super_task/internal/domain"
	"working/super_task/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublishedUseCase struct {
	PublishedRepos *repository.PublishedRepos
	Timeout        time.Duration
}

func NewPublishedUseCase(timeout time.Duration, repos *repository.PublishedRepos) *PublishedUseCase {
	return &PublishedUseCase{
		PublishedRepos: repos,
		Timeout:        timeout,
	}
}

// method for publishing video
func (pu *PublishedUseCase) PublishVideo(cxt context.Context, UserID string, TaskID string, publishedRequest *domain.PublishedRequest) (*domain.Published, error) {
	var published domain.Published

	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}
	taskID, err := primitive.ObjectIDFromHex(TaskID)
	if err != nil {
		return nil, err
	}
	published.ID = primitive.NewObjectID()
	published.UserID = userID
	published.TaskID = taskID

	published.Title = publishedRequest.Title
	published.Description = publishedRequest.Description
	published.Summery = publishedRequest.Summery

	published.VideoUrl = publishedRequest.VideoUrl
	published.PictureUrl = publishedRequest.PictureUrl
	published.PublishedDate = time.Now()

	return pu.PublishedRepos.CreatePublished(cxt, &published)

}

// method for editing videos meta information
func (pu *PublishedUseCase) EditVideo(cxt context.Context, UserID string, PublishedID string, edit *domain.UpdatePublished) (*domain.Published, error) {
	return pu.PublishedRepos.UpdatePublished(cxt, UserID, PublishedID, edit)
}

// method for updating like for the video
func (pu *PublishedUseCase) LikeVideo(cxt context.Context, UserID string, PublishedID string, value int) (*domain.Published, error) {
	return pu.PublishedRepos.UpdateLikes(cxt, UserID, PublishedID, value)
}

// method for updating dislike for the video
func (pu *PublishedUseCase) DislikeVideo(cxt context.Context, UserID string, Published string, value int) (*domain.Published, error) {
	return pu.PublishedRepos.UpdateDisLikes(cxt, UserID, Published, value)
}

// method for getting all published videos for the user
func (pu *PublishedUseCase) PublishedVideos(cxt context.Context, ID string, page int64, size int64) ([]*domain.Published, error) {
	return pu.PublishedRepos.GetPublisheds(cxt, ID, page, size)
}

// method for deleting video for the user
func (pu *PublishedUseCase) DeleteVideos(cxt context.Context, UserID string, PublishedID string) error {
	return pu.PublishedRepos.DeletePublished(cxt, UserID, PublishedID)
}

// method for creating new comments
func (pu *PublishedUseCase) NewComment(cxt context.Context, UserID string, PublishedID string, commentRequest *domain.CommentRequest) (*domain.Comments, error) {
	var comment domain.Comments
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}

	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return nil, err
	}

	comment.CommentID = primitive.NewObjectID()
	comment.UserID = userID
	comment.PublishedID = publishedID

	comment.CommenterID = commentRequest.CommenterID
	comment.Comment = commentRequest.Comment
	comment.CommentedDate = time.Now()
	return pu.PublishedRepos.CreateComment(cxt, &comment)

}

// method for editing comment by using infos
func (pu *PublishedUseCase) EditComment(cxt context.Context, commentID string, publishedID string, newData *domain.UpdateComment) (*domain.Comments, error) {
	return pu.PublishedRepos.EditComment(cxt, commentID, publishedID, newData)
}

// method for getting all comments
func (pu *PublishedUseCase) GetComments(cxt context.Context, userID string, publishedID string, page int64, size int64) ([]*domain.Comments, error) {
	return pu.PublishedRepos.GetComments(cxt, userID, publishedID, page, size)
}

// method for deleting comments
func (pu *PublishedUseCase) DeleteComment(cxt context.Context, commentID string, publishedID string) error {
	return pu.PublishedRepos.DeleteComment(cxt, commentID, publishedID)
}

// method for updating like for the comment
func (pu *PublishedUseCase) LikeVideoComment(cxt context.Context, UserID string, PublishedID string, commentID string , value int) (*domain.Comments, error) {
	return pu.PublishedRepos.UpdateLikesComment(cxt, UserID, PublishedID, commentID, value)
}

// method for updating dislike for the video
func (pu *PublishedUseCase) DislikeVideoComment(cxt context.Context, UserID string, Published string, commentID string, value int) (*domain.Comments, error) {
	return pu.PublishedRepos.UpdateDisLikesComment(cxt, UserID, Published, commentID, value)
}
