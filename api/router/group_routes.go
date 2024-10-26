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

func initProtectedGroupRoute(env *config.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewGroupRepos(db, domain.GroupCollection, domain.CollectionUser, domain.ConversationCollection, domain.MessageCollection)
	gc := controller.GroupController{
		GroupUseCase: usecase.NewGroupUseCase(tr, timeout),
		Env:          env,
	}
	group.GET("/group", gc.SendMessageHandler)
}
