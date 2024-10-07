package controller

import (
	"net/http"
	"working/super_task/config"
	usecase "working/super_task/internal/usercase"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	HomeUseCase *usecase.HomeUseCase
	Env         *config.Env
}

func NewHomeController(home *usecase.HomeUseCase, env *config.Env) *HomeController {
	return &HomeController{
		HomeUseCase: home,
		Env:         env,
	}
}

func (hc *HomeController) HomeHandler(c *gin.Context) {
	ID, exist := c.Get("id")
	if !exist {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data", "success": false, "data": nil})
		return
	}

	userID, okay := ID.(string)
	if !okay {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error of data type", "success": false, "data": nil})
		return
	}

	homeInformation, err := hc.HomeUseCase.HomeInformation(c, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "success": false, "data": nil})
		return
	}

	response := map[string]interface{}{
		"message": "user home information",
		"success": true,
		"data":    homeInformation,
	}
	c.IndentedJSON(http.StatusOK, response)

}
