package controller

import (
	"net/http"
	"strconv"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
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
