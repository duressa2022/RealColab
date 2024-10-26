package main

import (
	"fmt"
	"time"
	routes "working/super_task/api/router"
	"working/super_task/config"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	app, err := config.App()
	if err != nil {
		fmt.Println(err)
		return
	}

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	router := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},        // allow frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // allowed methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // allow credentials
	}

	router.Use(cors.New(config))
	routes.SetUpRoute(env, timeout, db, router)
	router.Run(env.ServerAddress)
}
