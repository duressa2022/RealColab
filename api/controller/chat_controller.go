package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type ChatController struct {
	ChatUseCase *usecase.ChatUseCase
	Env         *config.Env
}

func NewChatController(env *config.Env, chat *usecase.ChatUseCase) *ChatController {
	return &ChatController{
		Env:         env,
		ChatUseCase: chat,
	}
}
func (cc *ChatController) ConductChatHandler(c *gin.Context) {
	var prompt *domain.ChatRequest
	if err := c.BindJSON(&prompt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "succcess": false, "data": nil})
		return
	}
	Id, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "succcess": false, "data": nil})
		return
	}

	userID, ok := Id.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data type", "succcess": false, "data": nil})
		return
	}

	previousPrompt := cc.ChatUseCase.CreatePrompt(c, userID, prompt.SessionID)
	previousPrompt = fmt.Sprintf("%s \n: %s", previousPrompt, prompt.Prompt)
	cxt := context.Background()
	client, err := genai.NewClient(cxt, option.WithAPIKey(cc.Env.GEMINI_API))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": true, "data": nil})
		return
	}

	defer client.Close()
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(cxt, genai.Text(previousPrompt))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	content := resp.Candidates[0].Content.Parts

	summeryPrompt := fmt.Sprintf(`for this response give the name just by using 10 words: response: %s`, content)

	resp, err = model.GenerateContent(cxt, genai.Text(summeryPrompt))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	summery := resp.Candidates[0].Content.Parts

	_, err = cc.ChatUseCase.StoreChat(c, userID, prompt, content, summery)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	response := map[string]interface{}{
		"message":  "message is given",
		"succcess": true,
		"data":     content,
	}

	c.IndentedJSON(http.StatusOK, response)

}

// handler for working with chat fetching feature
func (cc *ChatController) FetchSessionHandler(c *gin.Context) {
	UserID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id not exist", "success": false, "data": nil})
		return
	}

	userID, okay := UserID.(string)
	if !okay {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "not correct data type", "success": false, "data": nil})
		return
	}

	limit := c.Query("limit")
	request := c.Query("request")

	limitNumber, err := strconv.Atoi(limit)
	if err != nil || limitNumber < 1 {
		limitNumber = 20
	}
	requestNumber, err := strconv.Atoi(request)
	if err != nil || requestNumber < 1 {
		requestNumber = 1
	}

	history, total, err := cc.ChatUseCase.FetchSession(c, userID, int64(requestNumber), int64(limitNumber))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	totalPage := (limitNumber + int(total) - 1) / limitNumber

	response := map[string]interface{}{
		"message":   "okay",
		"success":   true,
		"data":      history,
		"totalPage": totalPage,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for fetching chat for given session
func (cc *ChatController) FetchSessionChat(c *gin.Context) {
	type SessionRequest struct {
		SessionID string `json:"sessionID"`
	}
	UserID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id not found", "success": false, "data": nil})
		return
	}

	userID, ok := UserID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "not correct data form", "success": false, "data": nil})
		return
	}

	var sessionRequest *SessionRequest
	if err := c.BindJSON(&sessionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	chats, err := cc.ChatUseCase.FetchSessionChat(c, userID, sessionRequest.SessionID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	response := map[string]interface{}{
		"message": "okay",
		"success": true,
		"data":    chats,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// handler for deleting session from the system
func (cc *ChatController) DeleteSession(c *gin.Context) {
	type SessionRequest struct {
		SessionID string `json:"sessionID"`
	}
	UserID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id not found", "success": false, "data": nil})
		return
	}

	userID, ok := UserID.(string)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "not correct data form", "success": false, "data": nil})
		return
	}

	var sessionRequest *SessionRequest
	if err := c.BindJSON(&sessionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	err := cc.ChatUseCase.DeleteSession(c, userID, sessionRequest.SessionID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "okay",
		"success": true,
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// method for creating session for the first time
func (cc *ChatController) CreateSession(c *gin.Context) {
	new := c.Query("session")
	if new != "new" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty data", "success": false, "data": nil})
		return
	}

	UserID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id isnot found", "success": false, "data": nil})
		return
	}

	userID, okay := UserID.(string)
	if !okay {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "not correct data type", "success": false, "data": nil})
		return
	}

	session, err := cc.ChatUseCase.CreateSession(c, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "data": nil, "success": false})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    session,
	}
	c.IndentedJSON(http.StatusOK, response)

}
