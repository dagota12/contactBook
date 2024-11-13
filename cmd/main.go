package main

import (
	"findApi/api/routes"
	"findApi/bootstrap"
	"findApi/repository/db"

	"github.com/gin-gonic/gin"
)

func main(){
	env := bootstrap.LoadEnv()
	router := gin.Default()
	client := db.NewMongoClient(env)
	db := client.Database(env.DB_NAME)

	routes.SetupRoutes(router,db,env)
	router.Run(":"+env.PORT)
}