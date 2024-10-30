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

func initProtectedPublishedRoute(env *config.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	pr := repository.NewPublishedRepos(db, domain.PublishedCollection, domain.CommentCollection)
	pc := controller.PublishedController{
		PublishedUseCase: usecase.NewPublishedUseCase(timeout, pr),
		Env:              env,
	}
	group.POST("/publish", pc.PublishedVideosHandler)
	group.PUT("/edit:publishID", pc.EditVideoHandler)
	group.PUT("/like/:publishID", pc.LikeVideoHandler)
	group.PUT("/dislike/:publishID", pc.DisLikeVideoHandler)
	group.GET("/videos", pc.GetPublishedsHandler)
	group.DELETE("/delete/:publishID", pc.DeleteVideosHandler)
	group.POST("/comment", pc.CreateCommentHandler)
	group.PUT("/comment/edit", pc.EditCommentHandler)
	group.PUT("/like/:publishID", pc.LikeVideoHandler)
	group.PUT("/dislike/:publishID", pc.DisLikeVideoHandler)
	group.DELETE("/comment/delete/:publishID", pc.DeleteCommentHandler)
	group.GET("/comment/:publishedID", pc.CreateCommentHandler)
	group.PUT("/comment/like/:publishID", pc.LikeVideoHandler)
	group.PUT("/comment/dislike/:publishID", pc.DisLikeVideoHandler)

}
