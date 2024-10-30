package repository

import (
	"context"
	"errors"
	"working/super_task/internal/domain"
	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PublishedRepos struct {
	database            mongo.Database
	publishedCollection string
	commentsCollection  string
}

func NewPublishedRepos(db mongo.Database, published string, comments string) *PublishedRepos {
	return &PublishedRepos{
		database:            db,
		publishedCollection: published,
		commentsCollection:  comments,
	}
}

// method for adding new published into the database
func (pr *PublishedRepos) CreatePublished(cxt context.Context, published *domain.Published) (*domain.Published, error) {
	publishedCollection := pr.database.Collection(pr.publishedCollection)
	_, err := publishedCollection.InsertOne(cxt, published)
	if err != nil {
		return nil, err
	}
	return published, nil
}

// method for updating published by new data
func (pr *PublishedRepos) UpdatePublished(cxt context.Context, UserID string, PublishedID string, newData *domain.UpdatePublished) (*domain.Published, error) {
	publishedCollection := pr.database.Collection(pr.publishedCollection)
	updateInformation := bson.M{
		"title":       newData.Title,
		"description": newData.Description,
		"summery":     newData.Summery,
	}

	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return nil, err
	}

	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_userID", Value: userID}}
	updateResult, err := publishedCollection.UpdateOne(cxt, filter, bson.D{{Key: "$set", Value: updateInformation}})
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("no matched docs")
	}
	if updateResult.ModifiedCount == 0 {
		return nil, errors.New("no modified docs")
	}

	var published *domain.Published
	err = publishedCollection.FindOne(cxt, bson.D{{Key: "_publishedID", Value: published}}).Decode(&published)
	if err != nil {
		return nil, err
	}
	return published, nil

}

// method for updating like for the published
func (pr *PublishedRepos) UpdateLikes(cxt context.Context, UserID string, PublishedID string, value int) (*domain.Published, error) {
	publishedCollection := pr.database.Collection(pr.publishedCollection)
	updateInformation := bson.D{{Key: "$inc", Value: bson.D{{Key: "likes", Value: value}}}}

	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return nil, err
	}
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_userID", Value: userID}}
	updateResult, err := publishedCollection.UpdateOne(cxt, filter, updateInformation)
	if err != nil {
		return nil, err
	}
	if updateResult.ModifiedCount == 0 {
		return nil, errors.New("no docs are modified")
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("no docs are matched")
	}

	var published *domain.Published
	err = publishedCollection.FindOne(cxt, bson.D{{Key: "_publishedID", Value: publishedID}}).Decode(&published)
	if err != nil {
		return nil, err
	}
	return published, nil
}

// method for updating dislike for the published
func (pr *PublishedRepos) UpdateDisLikes(cxt context.Context, UserID string, PublishedID string, value int) (*domain.Published, error) {
	publishedCollection := pr.database.Collection(pr.publishedCollection)
	updateInformation := bson.D{{Key: "$inc", Value: bson.D{{Key: "disLikes", Value: value}}}}

	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return nil, err
	}
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_userID", Value: userID}}
	updateResult, err := publishedCollection.UpdateOne(cxt, filter, updateInformation)
	if err != nil {
		return nil, err
	}
	if updateResult.ModifiedCount == 0 {
		return nil, errors.New("no docs are modified")
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("no docs are matched")
	}

	var published *domain.Published
	err = publishedCollection.FindOne(cxt, bson.D{{Key: "_publishedID", Value: publishedID}}).Decode(&published)
	if err != nil {
		return nil, err
	}
	return published, nil
}

// method for getting all published videos for the user
func (pr *PublishedRepos) GetPublisheds(cxt context.Context, UserID string, page int64, size int64) ([]*domain.Published, error) {
	publishedCollection := pr.database.Collection(pr.publishedCollection)
	skips := (page - 1) * size
	opts := options.Find().SetSkip(skips).SetLimit(size).SetSort(bson.D{{Key: "timeStamp", Value: -1}})

	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}

	cursor, err := publishedCollection.Find(cxt, bson.D{{Key: "_userID", Value: userID}}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(cxt)

	var publisheds []*domain.Published
	for cursor.Next(cxt) {
		var published *domain.Published
		err := cursor.Decode(&published)
		if err != nil {
			return nil, err
		}
		publisheds = append(publisheds, published)
	}
	return publisheds, nil

}

// method deleting pusblished from database
func (pr *PublishedRepos) DeletePublished(cxt context.Context, UserID string, PublishedID string) error {
	publishedCollection := pr.database.Collection(pr.publishedCollection)

	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return err
	}
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_userID", Value: userID}}
	_, err = publishedCollection.DeleteOne(cxt, filter)
	return err

}

// method for creating new comments in the database
func (pr *PublishedRepos) CreateComment(cxt context.Context, comment *domain.Comments) (*domain.Comments, error) {
	commentCollection := pr.database.Collection(pr.commentsCollection)
	_, err := commentCollection.InsertOne(cxt, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// method for editing posted comments
func (pr *PublishedRepos) EditComment(cxt context.Context, CommentID string, PublishedID string, newData *domain.UpdateComment) (*domain.Comments, error) {
	commentCollection := pr.database.Collection(pr.commentsCollection)
	commentID, err := primitive.ObjectIDFromHex(CommentID)
	if err != nil {
		return nil, err
	}
	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return nil, err
	}

	updateInformation := bson.M{
		"comment": newData.Comment,
	}

	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_commentID", Value: commentID}}

	updatedResult, err := commentCollection.UpdateOne(cxt, filter, bson.D{{Key: "$set", Value: updateInformation}})
	if err != nil {
		return nil, err
	}
	if updatedResult.MatchedCount == 0 {
		return nil, errors.New("no matched count is found")
	}
	if updatedResult.ModifiedCount == 0 {
		return nil, errors.New("no modified count is found")
	}

	var comment *domain.Comments
	err = commentCollection.FindOne(cxt, filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return comment, nil

}

// method for getting all comments
func (pr *PublishedRepos) GetComments(cxt context.Context, UserID string, PublishedID string, page int64, limit int64) ([]*domain.Comments, error) {
	commentCollection := pr.database.Collection(pr.commentsCollection)
	skips := (page - 1) * limit
	opts := options.Find().SetLimit(limit).SetSkip(skips)

	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return nil, err
	}
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_userID", Value: userID}}

	cursor, err := commentCollection.Find(cxt, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(cxt)

	var comments []*domain.Comments
	for cursor.Next(cxt) {
		var comment *domain.Comments
		err := cursor.Decode(&comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, err

}

// method for deleting comment by commentID
func (pr *PublishedRepos) DeleteComment(cxt context.Context, CommentID string, PublishedID string) error {
	commentCollection := pr.database.Collection(pr.commentsCollection)

	commentID, err := primitive.ObjectIDFromHex(CommentID)
	if err != nil {
		return err
	}
	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_commentID", Value: commentID}}
	_, err = commentCollection.DeleteOne(cxt, filter)
	return err
}

// method for updating like for comments
func (pr *PublishedRepos) UpdateLikesComment(cxt context.Context, UserID string, PublishedID string, value int) (*domain.Comments, error) {
	commentCollection := pr.database.Collection(pr.commentsCollection)
	updateInformation := bson.D{{Key: "$inc", Value: bson.D{{Key: "likes", Value: value}}}}

	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return nil, err
	}
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_userID", Value: userID}}
	updateResult, err := commentCollection.UpdateOne(cxt, filter, updateInformation)
	if err != nil {
		return nil, err
	}
	if updateResult.ModifiedCount == 0 {
		return nil, errors.New("no docs are modified")
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("no docs are matched")
	}

	var comment *domain.Comments
	err = commentCollection.FindOne(cxt, bson.D{{Key: "_publishedID", Value: publishedID}}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// method for updating dislike for commens
func (pr *PublishedRepos) UpdateDisLikesComment(cxt context.Context, UserID string, PublishedID string, value int) (*domain.Comments, error) {
	commentCollection := pr.database.Collection(pr.commentsCollection)
	updateInformation := bson.D{{Key: "$inc", Value: bson.D{{Key: "disLikes", Value: value}}}}

	publishedID, err := primitive.ObjectIDFromHex(PublishedID)
	if err != nil {
		return nil, err
	}
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_publishedID", Value: publishedID}, {Key: "_userID", Value: userID}}
	updateResult, err := commentCollection.UpdateOne(cxt, filter, updateInformation)
	if err != nil {
		return nil, err
	}
	if updateResult.ModifiedCount == 0 {
		return nil, errors.New("no docs are modified")
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("no docs are matched")
	}

	var comment *domain.Comments
	err = commentCollection.FindOne(cxt, filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
