package controller

import (
	"net/http"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
)

type PublishedController struct {
	PublishedUseCase *usecase.PublishedUseCase
	Env              *config.Env
}

func NewPublishedController(published *usecase.PublishedUseCase, env *config.Env) *PublishedController {
	return &PublishedController{
		PublishedUseCase: published,
		Env:              env,
	}
}

// handler for working publishing video
func (pc *PublishedController) PublishedVideosHandler(c *gin.Context) {
	var published *domain.PublishedRequest
	if err := c.BindJSON(&published); err != nil {
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

	created, err := pc.PublishedUseCase.PublishVideo(c, userID, published.TaskID.String(), published)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    created,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for editing published videos
func (pc *PublishedController) EditVideoHandler(c *gin.Context) {
	var newData *domain.UpdatePublished
	if err := c.BindJSON(&newData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
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

	publishedID := c.Param("publishedID")
	if publishedID != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error empty data", "success": false, "data": nil})
		return
	}

	updated, err := pc.PublishedUseCase.EditVideo(c, userID, publishedID, newData)
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

// handler for working with like/video version
func (pc *PublishedController) LikeVideoHandler(c *gin.Context) {
	type LikeValue struct {
		Value bool `json:"value"`
	}

	var Value *LikeValue
	if err := c.BindJSON(&Value); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}
	var value int
	if Value.Value {
		value = 1
	} else {
		value = -1
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

	publishedID := c.Param("publishedID")
	if publishedID != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error empty data", "success": false, "data": nil})
		return
	}

	updated, err := pc.PublishedUseCase.LikeVideo(c, userID, publishedID, value)
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

// handler for working with dislike/video version
func (pc *PublishedController) DisLikeVideoHandler(c *gin.Context) {
	type DisLikeValue struct {
		Value bool `json:"value"`
	}

	var Value *DisLikeValue
	if err := c.BindJSON(&Value); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}
	var value int
	if Value.Value {
		value = 1
	} else {
		value = -1
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

	publishedID := c.Param("publishedID")
	if publishedID != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error empty data", "success": false, "data": nil})
		return
	}

	updated, err := pc.PublishedUseCase.DislikeVideo(c, userID, publishedID, value)
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
