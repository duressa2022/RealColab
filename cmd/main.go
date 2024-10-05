package main

import (
	"fmt"
	"time"

	routes "working/super_task/api/router"
	"working/super_task/config"

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
	routes.SetUpRoute(env, timeout, db, router)
	router.Run(env.ServerAddress)
}
