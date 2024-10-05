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

// method for init public route for user
func initPublicUserRoutes(env *config.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(domain.CollectionUser, db)
	uc := controller.UserController{
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Env:         env,
	}
	group.POST("/register", uc.RegisterUser)
	group.POST("/login", uc.Login)
}

// method for init protected for user
//func initProtectedUserRoutes(env *config.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
//	ur := repository.NewUserRepository(domain.CollectionUser, db)
//	uc := controllers.UserController{
//		UserUseCase: usecase.NewUserUseCase(timeout, ur),
//		Env:         env,
//	}

//}
