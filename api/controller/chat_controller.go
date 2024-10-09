package controller

import (
	"context"
	"fmt"
	"net/http"
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

	previousPrompt := cc.ChatUseCase.CreatePrompt(c, userID)
	previousPrompt = fmt.Sprintf("%s \n prompt: %s", previousPrompt, prompt.Prompt)

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
	
	response := map[string]interface{}{
		"message":  "message is given",
		"succcess": true,
		"data":     content,
	}

	c.IndentedJSON(http.StatusOK, response)

}
