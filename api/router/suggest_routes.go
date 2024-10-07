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

func initProtectedSuggestRoute(env *config.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewTaskRepository(domain.TaskCollection, db)
	sc := controller.SuggestContoller{
		SuggestUseCase: usecase.NewSuggestUseCase(tr, timeout),
		Env:            env,
	}
	group.GET("/suggest", sc.InitSuggestController)
}
