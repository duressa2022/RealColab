package routes

import (
	"time"
	"working/super_task/api/controller"
	"working/super_task/config"
	"working/super_task/internal/domain"
	"working/super_task/internal/repository"
	usecase "working/super_task/internal/usercase"
	"working/super_task/package/mongo"

	"github.com/gin-gonic/gin"
)

func initProtectedChatRoute(env *config.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewChatRepository(db, domain.ChatCollection)
	cc := controller.ChatController{
		ChatUseCase: usecase.NewChatUseCase(tr, timeout),
		Env:         env,
	}
	group.POST("/chat", cc.ConductChatHandler)
}
