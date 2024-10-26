package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	userID := c.Query("username")

	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	defer func() {
		mutex.Lock()
		delete(Clients, userID)
		mutex.Unlock()
		connection.Close()
	}()

	client := &domain.Client{
		Connection: connection,
		UserID:     userID,
	}
	mutex.Lock()
	Clients[userID] = client
	mutex.Unlock()

	for {
		_, messageByte, err := connection.ReadMessage()
		if err != nil {
			client.Connection.WriteMessage(websocket.TextMessage, []byte(`"error":"error of disconnection"`))
			mutex.Lock()
			delete(Clients, userID)
			mutex.Unlock()
			return
		}

		var message *domain.OneToOneMessage
		if err := json.Unmarshal(messageByte, &message); err != nil {
			client.Connection.WriteMessage(websocket.TextMessage, []byte(`"error":"error of disconnection"`))
			return
		}
		message.SenderID = userID
		SendMessage(Clients, message)

	}
}

// method for sending the message for specific users
func SendMessage(connections map[string]*domain.Client, message *domain.OneToOneMessage) error {
	messageByte, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client, ok := connections[message.ReceipentID]
	if !ok {
		return errors.New("error of data type")
	}

	err = client.Connection.WriteMessage(websocket.TextMessage, messageByte)
	if err != nil {
		return err
	}

	return nil
}

// handler for adding new contacts
func (mc *MessageController) AddContactHandler(c *gin.Context) {
	var contact *domain.Contact
	if err := c.BindJSON(&contact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	UserID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error od id ", "success": false, "data": nil})
		return
	}

	userID, ok := UserID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of type", "success": false, "data": nil})
		return
	}

	contactInfo, err := mc.MessageConUseCase.AddContact(c, userID, contact)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "okay",
		"success": true,
		"data":    contactInfo,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for searching user of the system by using username
func (mc *MessageController) SearchUserHandler(c *gin.Context) {
	type SearchValue struct {
		SearchTerm string `json:"searchTerm"`
	}
	var searchTerm *SearchValue
	if err := c.BindJSON(&searchTerm); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "Data": nil})
		return
	}

	searchResult, err := mc.MessageConUseCase.SearchUser(c, searchTerm.SearchTerm)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "okay",
		"success": true,
		"data":    searchResult,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// handler for fetching message history
func (mc *MessageController) FetchMessageHistoryHandler(c *gin.Context) {
	type Contact struct {
		ContactID string `json:"contactID"`
	}
	var contact *Contact
	if err := c.BindJSON(&contact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	UserID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error od id ", "success": false, "data": nil})
		return
	}
	userID, ok := UserID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of type", "success": false, "data": nil})
		return
	}

	request := c.Query("request")
	limit := c.Query("limit")

	requestNumber, err := strconv.Atoi(request)
	if err != nil || requestNumber < 1 {
		requestNumber = 1
	}
	limitNumber, err := strconv.Atoi(limit)
	if err != nil || limitNumber < 1 {
		limitNumber = 50
	}

	messages, err := mc.MessageConUseCase.FetchMessagesHistory(c, userID, contact.ContactID, int64(requestNumber), int64(limitNumber))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "okay",
		"success": true,
		"data":    messages,
	}
	c.IndentedJSON(http.StatusOK, response)
}
