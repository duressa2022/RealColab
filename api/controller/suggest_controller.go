package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"working/super_task/config"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type SuggestContoller struct {
	SuggestUseCase *usecase.SuggestUseCase
	Env            *config.Env
}

func NewSuggestController(suggest *usecase.SuggestUseCase, env *config.Env) *SuggestContoller {
	return &SuggestContoller{
		SuggestUseCase: suggest,
		Env:            env,
	}
}

// handler for working with the suggestions for the first time
func (sc *SuggestContoller) InitSuggestController(c *gin.Context) {
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

	prompt, err := sc.SuggestUseCase.CreatePrompt(c, userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(sc.Env.GEMINI_API))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	defer client.Close()

	mode := client.GenerativeModel("gemini-1.5-flash")
	resp, err := mode.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	if len(resp.Candidates) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no suggestion", "success": false, "data": nil})
		return
	}

	content := resp.Candidates[0].Content.Parts
	preproccesed := ""
	for index, part := range content {
		preproccesed += fmt.Sprintf("suggest%d: %s\n", index+1, part)
	}

	suggestion := GetSuggestions(preproccesed)
	delete(suggestion, "suggest1")
	response := map[string]interface{}{
		"message": "suggestions",
		"success": true,
		"data":    suggestion,
	}

	c.IndentedJSON(http.StatusOK, response)

}

func GetSuggestions(response string) map[string]string {
	parts := strings.Split(response, "\n")
	result := map[string]string{}

	for _, part := range parts {
		if strings.Contains(part, ":") {

			kv := strings.SplitN(part, ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				result[key] = value
			}
		}
	}
	return result
}
