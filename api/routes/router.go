package routes

import (
	"findApi/bootstrap"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


func SetupRoutes(router *gin.Engine,db *mongo.Database,env *bootstrap.Env) {

	NewUserRoute(router,db,env)

}