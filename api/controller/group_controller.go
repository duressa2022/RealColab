package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupController struct {
	GroupUseCase *usecase.GroupUseCase
	Env          *config.Env
}

func NewGroupController(env *config.Env, group *usecase.GroupUseCase) *GroupController {
	return &GroupController{
		GroupUseCase: group,
		Env:          env,
	}
}

var group_upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var GroupClient = make(map[string]map[string]*domain.GroupClient)
var group_mutex sync.Mutex

func (gc *GroupController) SendMessage(c *gin.Context) {
	UserID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	userID, ok := UserID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of type", "success": false, "data": nil})
		return
	}

	groupID := c.Param("groupID")

	connection, err := group_upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil, "success": false})
		return
	}

	defer connection.Close()
	client := &domain.GroupClient{
		Connection: connection,
		UserID:     userID,
		GroupID:    groupID,
	}

	group_mutex.Lock()
	if _, ok := GroupClient[groupID]; !ok {
		GroupClient[groupID] = make(map[string]*domain.GroupClient)
	}
	GroupClient[groupID][userID] = client
	group_mutex.Unlock()

	for {
		messageType, messageByte, err := connection.ReadMessage()
		if err != nil {
			client.Connection.WriteMessage(websocket.TextMessage, []byte(`"error":"error of disconnection"`))
			mutex.Lock()
			delete(Clients, userID)
			mutex.Unlock()
			return
		}

		err = gc.messageHelperMethod(c, messageByte, messageType, client)
		if err != nil {
			client.Connection.WriteMessage(websocket.TextMessage, []byte(`{"error":"message transmission is halted"}`))
		}

	}

}

// helper function for working message writing and message storing
func (gc *GroupController) messageHelperMethod(cxt context.Context, messageByte []byte, messageType int, client *domain.GroupClient) error {
	var message *domain.GroupMessage
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

	groupID, err := primitive.ObjectIDFromHex(client.GroupID)
	if err != nil {
		return err
	}
	message.GroupID = groupID
	conversation, err := gc.GroupUseCase.CreateOrRetriveConversation(cxt, *message)
	if err != nil {
		return err
	}
	message.ConversationId = conversation.ConversationID

	_, err = gc.GroupUseCase.StoreMessage(cxt, message)
	if err != nil {
		return err
	}

	mutex.Lock()
	mutex.Lock()
	for _, recipient := range GroupClient[client.GroupID] {
		if recipient.UserID != client.UserID {
			err := recipient.Connection.WriteJSON(message)
			if err != nil {
				delete(GroupClient[client.GroupID], recipient.UserID)
			}
		}
	}
	mutex.Unlock()
	mutex.Unlock()
	return nil
}

func (gc *GroupController) CreateGroupHandler(c *gin.Context) {
	var group *domain.GroupRequest
	if err := c.BindJSON(&group); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

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

	createdGroup, err := gc.GroupUseCase.CreateGroup(c, group, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    createdGroup,
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) UpdateGroupHandler(c *gin.Context) {
	var group *domain.GroupRequest
	if err := c.BindJSON(&group); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

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

	updated, err := gc.GroupUseCase.UpdateGroup(c, group, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    updated,
	}
	c.IndentedJSON(http.StatusOK, response)

}

func (gc *GroupController) DeleteGroupHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "data is not provided", "succcess": false, "data": nil})
		return
	}

	err := gc.GroupUseCase.DeleteGroup(c, groupID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) GetGroupInformationHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no data is given", "success": false, "data": nil})
		return
	}

	group, err := gc.GroupUseCase.GetGroupInformation(c, groupID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message":  "ok",
		"succcess": true,
		"data":     group,
	}
	c.IndentedJSON(http.StatusOK, response)

}

func (gc *GroupController) AddMemberHandler(c *gin.Context) {
	var user *domain.UserInGroup
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil, "success": false})
		return
	}

	groupID := c.Param("groupID")

	err := gc.GroupUseCase.AddMember(c, groupID, user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) GetAllMembersHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data is provided", "data": nil, "success": false})
		return
	}

	members, err := gc.GroupUseCase.GetAllMembers(c, groupID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "ok",
		"data":    members,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) DeleteMemberHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data is provided", "data": nil, "success": false})
		return
	}

	userID := c.Param("userID")
	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data", "success": false, "data": nil})
	}

	err := gc.GroupUseCase.DeleteMember(c, userID, groupID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "ok",
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) PromoteHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data is provided", "data": nil, "success": false})
		return
	}

	var user *domain.UserInGroup
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil, "success": false})
		return
	}
	err := gc.GroupUseCase.PromoteUser(c, user, groupID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "ok",
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) DemoteHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data is provided", "data": nil, "success": false})
		return
	}

	userID := c.Param("userID")
	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data", "success": false, "data": nil})
	}

	err := gc.GroupUseCase.DemoteUser(c, userID, groupID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "ok",
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) BlockedUserHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data is provided", "data": nil, "success": false})
		return
	}

	blocked := map[string]string{
		"_userID": "userID",
		"reason":  "reason",
	}
	if err := c.BindJSON(&blocked); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil, "success": false})
		return
	}

	userID := c.Param("userID")
	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data", "success": false, "data": nil})
	}

	err := gc.GroupUseCase.BlockUser(c, userID, groupID, blocked["reason"])
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "ok",
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) UnBlockUserHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data is provided", "data": nil, "success": false})
		return
	}

	userID := c.Param("userID")
	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data", "success": false, "data": nil})
	}

	err := gc.GroupUseCase.UnblockUser(c, userID, groupID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "ok",
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func (gc *GroupController) GetMessages(c *gin.Context) {
	groupID := c.Param("groupID")
	if groupID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data", "success": false, "data": nil})
		return
	}

	size := c.Query("size")
	page := c.Query("page")

	sizeNumber, err := strconv.Atoi(size)
	if err != nil || sizeNumber < 1 {
		sizeNumber = 20
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	messages, totalMessages, err := gc.GroupUseCase.GetMessages(c, groupID, int64(sizeNumber), int64(pageNumber))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	totalPage := (sizeNumber + int(totalMessages)) / sizeNumber

	response := map[string]interface{}{
		"message":    "ok",
		"success":    true,
		"data":       messages,
		"totalpages": totalPage,
	}
	c.IndentedJSON(http.StatusOK, response)
}
