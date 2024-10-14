package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageController struct {
	MessageConUseCase *usecase.MessageConversation
	Env               *config.Env
}

func NewMessageController(env *config.Env, messageCon *usecase.MessageConversation) *MessageController {
	return &MessageController{
		MessageConUseCase: messageCon,
		Env:               env,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[string]*domain.Client)
var mutex sync.Mutex

// handler for working with one to one messaging
func (mc *MessageController) MessageHandler(c *gin.Context) {
	ID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	userID, ok := ID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data type", "success": false, "data": nil})
		return
	}

	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	defer connection.Close()

	client := &domain.Client{
		Connection: connection,
		UserID:     userID,
	}
	mutex.Lock()
	Clients[userID] = client
	mutex.Unlock()

	for {
		messageType, messageByte, err := connection.ReadMessage()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
			return
		}
		mc.messageHelperMethod(c, messageByte, messageType, client)
	}

}

// helper function for working message writing and message storing
func (mc *MessageController) messageHelperMethod(cxt context.Context, messageByte []byte, messageType int, client *domain.Client) error {
	var message *domain.Message
	if err := json.Unmarshal(messageByte, &message); err != nil {
		return err
	}

	message.TimeStamp = primitive.NewDateTimeFromTime(time.Now())
	message.Status = "unread"
	message.MessageType = messageType
	senderID, err := primitive.ObjectIDFromHex(client.UserID)
	if err != nil {
		return err
	}
	message.SenderID = senderID
	conversation, err := mc.MessageConUseCase.CreateOrRetriveConversation(cxt,
		&domain.ParticipantRequest{FirstParticipantID: senderID, SecondParticipantID: message.ReceipentID}, *message)
	if err != nil {
		return err
	}
	message.ConversationId = conversation.ConversationID

	_, err = mc.MessageConUseCase.StoreMessage(cxt, message)
	if err != nil {
		return err
	}

	mutex.Lock()
	receipentClient, ok := Clients[message.ReceipentID.String()]
	mutex.Unlock()
	if ok {
		err := receipentClient.Connection.WriteMessage(websocket.TextMessage, messageByte)
		if err != nil {
			return err
		}
	} else {
		return errors.New("user is not online")
	}
	return nil
}
