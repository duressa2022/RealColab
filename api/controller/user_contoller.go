package controller

import (
	"net/http"
	"working/super_task/config"
	"working/super_task/internal/domain"
	usecase "working/super_task/internal/usercase"

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

// registration handler
func (uc *UserController) RegisterUser(c *gin.Context) {
	var userInformation *domain.UserRegistrationRequest
	if err := c.BindJSON(&userInformation); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid information", "data": nil})
		return
	}

	createdUser, err := uc.UserUseCase.RegisterUser(c, userInformation)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"success": false, "message": err, "data": nil})
		return
	}
	response := map[string]interface{}{
		"success": true,
		"message": "registerd successfully",
		"data":    createdUser,
	}
	c.IndentedJSON(http.StatusOK, response)
}
