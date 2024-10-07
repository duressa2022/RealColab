package controller

import (
	"net/http"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"
	tokens "working/super_task/package/token"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUseCase *usecase.UserUseCase
	Env         *config.Env
}

func NewUserController(env *config.Env, userUsecase *usecase.UserUseCase) *UserController {
	return &UserController{
		UserUseCase: userUsecase,
		Env:         env,
	}
}

// method for updating notification choice
func (uc *UserController) UpdateNotificationChoice(c *gin.Context) {
	Id, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "not found id of error", "success": false, "data": nil})
		return
	}

	userid, okay := Id.(string)
	if !okay {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	var change *domain.NotificationPreference
	if err := c.BindJSON(&change); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found", "data": nil, "success": true})
		return
	}

	changed, err := uc.UserUseCase.UpdateNotificationChoice(c, change, userid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found", "data": nil, "success": true})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "updated",
		"data":    changed,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// method for getting notification choice
func (uc *UserController) GetNotificationChoice(c *gin.Context) {
	Id, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "not found id of error", "success": false, "data": nil})
		return
	}

	userid, okay := Id.(string)
	if !okay {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	notification, err := uc.UserUseCase.GetNotificationChoice(c, userid)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": true, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "data",
		"data":    notification,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// method for changing password
func (uc *UserController) ChangePassword(c *gin.Context) {
	Id, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "not found id of error", "success": false, "data": nil})
		return
	}

	var changeInfo *domain.ChangePassword
	if err := c.BindJSON(&changeInfo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	userid, okay := Id.(string)
	if !okay {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	err := uc.UserUseCase.ChangePassword(c, changeInfo, userid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}
	response := map[string]interface{}{
		"success": true,
		"message": "changed",
		"data":    nil,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// get user security information
func (uc *UserController) GetSecurityInfo(c *gin.Context) {
	Id, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found id of error", "success": false, "data": nil})
		return
	}

	userid, okay := Id.(string)
	if !okay {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "not found id of error", "success": false, "data": nil})
		return
	}

	userSecurity, err := uc.UserUseCase.GetSecurityInfo(c, userid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "message got",
		"data":    userSecurity,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// update profile handler
func (uc *UserController) UpdateMainInfo(c *gin.Context) {
	var updateUser *domain.UserUpdateMainInfo
	if err := c.BindJSON(&updateUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}

	Id, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found email", "success": false, "data": nil})
		return
	}

	userid, okay := Id.(string)
	if !okay {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "data type erro", "success": false, "data": nil})
		return
	}

	updatedUser, err := uc.UserUseCase.UpdateMainInfo(c, updateUser, userid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "updated",
		"data":    updatedUser,
	}
	c.IndentedJSON(http.StatusOK, response)

}

// registration handler
func (uc *UserController) RegisterUser(c *gin.Context) {
	var userInformation *domain.UserRegistrationRequest
	if err := c.BindJSON(&userInformation); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid information", "data": nil})
		return
	}

	createdUser, err := uc.UserUseCase.RegisterUser(c, userInformation)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}
	response := map[string]interface{}{
		"success": true,
		"message": "registerd successfully",
		"data":    createdUser,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handler for working with login
func (uc *UserController) Login(c *gin.Context) {
	var login *domain.LoginRequest
	if err := c.BindJSON(&login); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"success": false, "message": "error of data", "data": nil})
		return
	}

	userInformation, err := uc.UserUseCase.LoginUser(c, login)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}

	accessToken, err := tokens.CreateAccessToken(userInformation, uc.Env.AccessTokenSecret, uc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}

	refreshToken, err := tokens.CreateRefreshToken(userInformation, uc.Env.RefreshTokenSecret, uc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}
	c.SetCookie("accessToken", accessToken, uc.Env.AccessTokenExpiryHour, "", "", true, true)
	c.SetCookie("refreshToken", refreshToken, uc.Env.RefreshTokenExpiryHour, "", "", true, true)

	response := map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"success":      true,
		"message":      "login in Succcesful",
		"data":         map[string]interface{}{},
	}
	c.IndentedJSON(http.StatusOK, response)
}
