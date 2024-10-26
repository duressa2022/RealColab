package repository

import (
	"context"
	"errors"
	"time"
	"working/super_task/internal/domain"

	"working/super_task/package/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	database               mongo.Database
	messageCollection      string
	conversationCollection string
	userCollection         string
}

func NewMessageRepos(db mongo.Database, collectionMessage string, collectionConversation string, userCollection string) *MessageRepository {
	return &MessageRepository{
		database:               db,
		messageCollection:      collectionMessage,
		conversationCollection: collectionConversation,
		userCollection:         userCollection,
	}
}

// method for creating new conversation or retiving existing one
func (mr *MessageRepository) CreateOrUpdateConversation(cxt context.Context, participant *domain.ParticipantRequest, message domain.Message) (*domain.Conversation, error) {
	convCollection := mr.database.Collection(mr.conversationCollection)
	filter := bson.D{
		{Key: "participants", Value: participant.FirstParticipantID},
		{Key: "participants", Value: participant.SecondParticipantID},
	}

	var convModel *domain.Conversation
	err := convCollection.FindOne(cxt, filter).Decode(&convModel)
	if err != nil {
		newConvMode := &domain.Conversation{
			ConversationID: primitive.NewObjectID(),
			IsGroup:        false,
			Participants: []domain.Participant{
				{ParticipantID: participant.FirstParticipantID},
				{ParticipantID: participant.SecondParticipantID},
			},
			LastMessage: message,
			CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		}
		convModel = newConvMode
	} else {
		convModel.LastMessage = message
		convModel.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	}
	return convModel, nil

}

// method for fetching message history from the database
func (mr *MessageRepository) FetchMessages(cxt context.Context, ConvID string, page int64, size int64) ([]*domain.Message, error) {
	messageCollection := mr.database.Collection(mr.messageCollection)
	skip := (page - 1) / size
	opts := options.Find().SetSkip(skip).SetLimit(size).SetSort(bson.D{{Key: "timeStamp", Value: -1}})
	convID, err := primitive.ObjectIDFromHex(ConvID)
	if err != nil {
		return nil, err
	}

	var messages []*domain.Message
	cursor, err := messageCollection.Find(cxt, bson.D{{Key: "_conversationID", Value: convID}}, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(cxt) {
		var message *domain.Message
		err := cursor.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// method for storing message on the database
func (cr *MessageRepository) StoreMessage(cxt context.Context, message *domain.Message) (*domain.Message, error) {
	messageCollection := cr.database.Collection(cr.messageCollection)
	_, err := messageCollection.InsertOne(cxt, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// method for getting all unread message
func (cr *MessageRepository) UnReadMessage(cxt context.Context, ID string) ([]*domain.Message, error) {
	messageCollection := cr.database.Collection(cr.messageCollection)
	conversationID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	var messages []*domain.Message
	cursor, err := messageCollection.Find(cxt, bson.D{
		{Key: "_conversationID", Value: conversationID},
		{Key: "status", Value: "unread"}})

	if err != nil {
		return nil, err
	}

	for cursor.Next(cxt) {
		var message *domain.Message
		err := cursor.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil

}

// method for creating/adding new for the user
func (cr *MessageRepository) AddContact(cxt context.Context, UserID string, contact *domain.Contact) (map[string]interface{}, error) {
	userCollection := cr.database.Collection(cr.userCollection)
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}

	var user *domain.UserResponse
	err = userCollection.FindOne(cxt, bson.D{{Key: "_userID", Value: userID}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	updateUser := bson.M{
		"contacts": append(user.Contacts, contact),
	}

	updated, err := userCollection.UpdateOne(cxt, bson.D{{Key: "_userID", Value: userID}}, bson.D{{Key: "$set", Value: updateUser}})
	if err != nil {
		return nil, err
	}
	if updated.ModifiedCount == 0 {
		return nil, errors.New("no modified docs")
	}
	if updated.MatchedCount == 0 {
		return nil, errors.New("no matched docs")
	}
	response, err := cr.GetContactInfo(cxt, contact.ContactID.String())
	return response, err

}

// method for getting contact information by using id
func (cr *MessageRepository) GetContactInfo(cxt context.Context, ContactID string) (map[string]interface{}, error) {
	userCollection := cr.database.Collection(cr.userCollection)
	contactID, err := primitive.ObjectIDFromHex(ContactID)
	if err != nil {
		return nil, err
	}

	var user *domain.UserResponse
	err = userCollection.FindOne(cxt, bson.D{{Key: "_userID", Value: contactID}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"firstName":  user.FirstName,
		"lastName":   user.LastName,
		"userName":   user.UserName,
		"profileUrl": user.ProfileUrl,
		"status":     user.Status,
	}
	return response, nil

}

// method for searching user/contacts by using username
func (cr *MessageRepository) SearchUser(cxt context.Context, searchTerm string) ([]map[string]interface{}, error) {
	userCollection := cr.database.Collection(cr.userCollection)
	filter := bson.D{{Key: "username", Value: bson.D{{Key: "$regex", Value: searchTerm}, {Key: "$options", Value: "i"}}}}

	cursor, err := userCollection.Find(cxt, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(cxt)

	var responses []map[string]interface{}
	for cursor.Next(cxt) {
		var user *domain.UserResponse
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		response := map[string]interface{}{
			"firstName":  user.FirstName,
			"lastName":   user.LastName,
			"username":   user.UserName,
			"profileUrl": user.ProfileUrl,
			"status":     user.Status,
		}
		responses = append(responses, response)

	}
	return responses, nil
}

// method for fetching message history
func (cr *MessageRepository) FetchMessageHistory(cxt context.Context, UserID string, ContactID string, requestNumber int64, limit int64) ([]*domain.Message, error) {
	messageCollection := cr.database.Collection(cr.messageCollection)
	userID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return nil, err
	}
	contactID, err := primitive.ObjectIDFromHex(ContactID)
	if err != nil {
		return nil, err
	}

	skip := (requestNumber - 1) * limit
	filter := bson.D{{Key: "_userID", Value: userID}, {Key: "_receipentID", Value: contactID}}
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "timeStamp", Value: 1}})

	cursor, err := messageCollection.Find(cxt, filter, opts)
	if err != nil {
		return nil, err
	}

	var messages []*domain.Message
	for cursor.Next(cxt) {
		var message *domain.Message
		err := cursor.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// method for deleting message from the database
func (cr *MessageRepository) DeletMessage(cxt context.Context, MessageID string) error {
	messageCollection := cr.database.Collection(cr.messageCollection)
	messageID, err := primitive.ObjectIDFromHex(MessageID)
	if err != nil {
		return err
	}
	_, err = messageCollection.DeleteOne(cxt, bson.D{{Key: "_messageID", Value: messageID}})
	return err
}

// method updating/editing  messages
func (cr *MessageRepository) EditMessage(cxt context.Context, MessageID string, message *domain.EditMessage) (*domain.Message, error) {
	messageCollection := cr.database.Collection(cr.messageCollection)
	editMessage := bson.M{
		"messageContent": message.MessageContent,
		"mediaUrl":       message.MediaUrl,
		"status":         message.Status,
	}

	messageID, err := primitive.ObjectIDFromHex(MessageID)
	if err != nil {
		return nil, err
	}

	updated, err := messageCollection.UpdateOne(cxt, bson.D{{Key: "_message", Value: messageID}}, bson.D{{Key: "$set", Value: editMessage}})
	if err != nil {
		return nil, err
	}
	if updated.MatchedCount == 0 {
		return nil, errors.New("no matched docs")
	}
	if updated.ModifiedCount == 0 {
		return nil, errors.New("no modified docs")
	}

	var UpdatedMessage *domain.Message
	err = messageCollection.FindOne(cxt, bson.D{{Key: "_messageID", Value: messageID}}).Decode(&UpdatedMessage)
	if err != nil {
		return nil, err
	}
	return UpdatedMessage, nil

}

// method for deleting entire chat/messaging history
func (cr *MessageRepository) DeleteMessageHistory(cxt context.Context, conversationID string) error {
	messageCollection := cr.database.Collection(cr.messageCollection)
	conversationColletion := cr.database.Collection(cr.conversationCollection)

	conveID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return err
	}

	deleted, err := messageCollection.DeleteMany(cxt, bson.D{{Key: "_conversationID", Value: conveID}})
	if deleted.DeletedCount == 0 {
		return errors.New("no item is deleted")
	}
	if err != nil {
		return err
	}

	_, err = conversationColletion.DeleteOne(cxt, bson.D{{Key: "_conversationID", Value: conveID}})
	if err != nil {
		return err
	}
	return nil
}
