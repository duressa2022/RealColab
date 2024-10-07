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

func initProtectedTaskRoute(env *config.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewTaskRepository(domain.TaskCollection, db)
	tc := controller.TaskController{
		TaskUseCase: usecase.NewTaskUseCase(timeout, tr),
		Env:         env,
	}
	group.POST("/addTask", tc.AddTaskHandler)
	group.GET("/", tc.GetTaskByTypeHandler)
	group.PUT("/archive/:taskID", tc.ArchiveTaskHandler)
	group.PUT("/edit/:taskID", tc.EditTaskHandler)
	group.GET("/search", tc.SearchTaskHandler)

}
