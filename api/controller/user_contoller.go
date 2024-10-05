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

// handler for working with login
func (uc *UserController) Login(c *gin.Context){
	var login *domain.LoginRequest
	if err:=c.BindJSON(&login);err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"success":false,"message":"error of data","data":nil})
		return 
	}

	userInformation,err:=uc.UserUseCase.LoginUser(c,login)
	if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError,gin.H{"success":false,"message":err.Error(),"data":nil})
		return 
	}

	accessToken,err:=tokens.CreateAccessToken(userInformation,uc.Env.AccessTokenSecret,uc.Env.AccessTokenExpiryHour)
	if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError,gin.H{"success":false,"message":err.Error(),"data":nil})
		return 
	}

	refreshToken,err:=tokens.CreateRefreshToken(userInformation,uc.Env.RefreshTokenSecret,uc.Env.RefreshTokenExpiryHour)
	if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError,gin.H{"success":false,"message":err.Error(),"data":nil})
		return 
	}
	c.SetCookie("accessToken",accessToken,uc.Env.AccessTokenExpiryHour,"","",true,true)
	c.SetCookie("refreshToken",refreshToken,uc.Env.RefreshTokenExpiryHour,"","",true,true)

	response:=map[string]interface{}{
		"accessToken":accessToken,
		"refreshToken":refreshToken,
		"success":true,
		"message":"login in Succcesful",
		"data":map[string]interface{}{},
	}
	c.IndentedJSON(http.StatusOK,response)
}
