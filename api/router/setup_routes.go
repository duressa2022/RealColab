package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	//"working/super_task/api/middleware"
	middlewares "working/super_task/api/middleware"
	"working/super_task/config"
	"working/super_task/package/mongo"
)

// method for setting the route
func SetUpRoute(env *config.Env, timeout time.Duration, db mongo.Database, router *gin.Engine) {
	publicRoute := router.Group("/auth")
	initPublicUserRoutes(env, timeout, db, publicRoute)
	initProtectedChatRoute(env, timeout, db, publicRoute)
	initProtectedMessagingRoute(env, timeout, db, publicRoute)
	initProtectedGroupRoute(env, timeout, db, publicRoute)

	protectedRoute := router.Group("/", middlewares.JwtAuthMiddleWare(env.AccessTokenSecret))
	initProtectedUserRoutes(env, timeout, db, protectedRoute.Group("user"))
	initProtectedTaskRoute(env, timeout, db, protectedRoute.Group("tasks"))
	initProtectedSuggestRoute(env, timeout, db, protectedRoute.Group("suggest"))
	initProtectedChatRoute(env, timeout, db, protectedRoute.Group("chats"))

}
