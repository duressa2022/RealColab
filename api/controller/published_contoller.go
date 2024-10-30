package controller

import (
	"net/http"
	"strconv"
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

// handler for working/getting published videos by using id
func (pc *PublishedController) GetPublishedsHandler(c *gin.Context) {
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

	limit := c.Query("limit")
	page := c.Query("page")

	limitNumber, err := strconv.Atoi(limit)
	if err != nil || limitNumber < 1 {
		limitNumber = 50
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	published, err := pc.PublishedUseCase.PublishedVideos(c, userID, int64(pageNumber), int64(limitNumber))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    published,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for deleting video for the user
func (pc *PublishedController) DeleteVideosHandler(c *gin.Context) {
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
	if publishedID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
	}

	err := pc.PublishedUseCase.DeleteVideos(c, userID, publishedID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": true, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// handler for creating new comments
func (pc *PublishedController) CreateCommentHandler(c *gin.Context) {
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
	if publishedID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
	}

	var comment *domain.CommentRequest
	if err := c.BindJSON(&comment); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": true, "data": nil})
		return
	}

	created, err := pc.PublishedUseCase.NewComment(c, userID, publishedID, comment)
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

// handler for editing the comments
func (pc *PublishedController) EditCommentHandler(c *gin.Context) {
	var editRequest *domain.UpdateComment
	if err := c.BindJSON(&editRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": true, "data": nil})
		return
	}

	updated, err := pc.PublishedUseCase.EditComment(c, editRequest.CommentID.String(), editRequest.PublishedID.String(), editRequest)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": true, "data": nil})
		return
	}
	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    updated,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// handler for working/getting comment for given videos
func (pc *PublishedController) GetCommentsHandler(c *gin.Context) {
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

	limit := c.Query("limit")
	page := c.Query("page")

	limitNumber, err := strconv.Atoi(limit)
	if err != nil || limitNumber < 1 {
		limitNumber = 50
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	publishedID := c.Param("publishedID")

	comments, err := pc.PublishedUseCase.GetComments(c, userID, publishedID, int64(pageNumber), int64(limitNumber))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    comments,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for deleting the comment
func (pc *PublishedController) DeleteCommentHandler(c *gin.Context) {
	var deleteRequest *domain.UpdateComment
	if err := c.BindJSON(&deleteRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": true, "data": nil})
		return
	}

	err := pc.PublishedUseCase.DeleteComment(c, deleteRequest.CommentID.String(), deleteRequest.PublishedID.String())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": true, "data": nil})
		return
	}
	response := map[string]interface{}{
		"message": "ok",
		"success": true,
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)

}
// handler for working with like comments version
func (pc *PublishedController) LikeCommentHandler(c *gin.Context) {
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

	var commentInfo *domain.CommentInfo
	if err:=c.BindJSON(&commentInfo);err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":err.Error(),"success":false,"data":nil})
	}

	updated, err := pc.PublishedUseCase.LikeVideoComment(c,userID,commentInfo.PublishedID.String(),commentInfo.CommentID.String(),value)
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

// handler for working with disliking video version
func (pc *PublishedController) DisLikeCommentHandler(c *gin.Context) {
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

	var commentInfo *domain.CommentInfo
	if err:=c.BindJSON(&commentInfo);err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":err.Error(),"success":false,"data":nil})
	}

	updated, err := pc.PublishedUseCase.DislikeVideoComment(c,userID,commentInfo.PublishedID.String(),commentInfo.CommentID.String(),value)
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
