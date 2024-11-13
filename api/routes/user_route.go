package routes

import (
	"findApi/api/controller"
	"findApi/bootstrap"
	"findApi/repository"
	"findApi/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRoute(r *gin.Engine, db *mongo.Database, env *bootstrap.Env) {
	repo := repository.NewUserRepository(db.Collection("users"),env)
	usecase := usecase.NewUsersUseCase(repo)
	controller := &controller.UserController{UserUsecase: usecase, Env: env}
	r.POST("/users", controller.CreateUser)          // Create a new user
	r.GET("/users/username/:username", controller.GetUserByUsername) // Get user by username
	r.GET("/users/phone/:phone", controller.GetUserByPhone)         // Get user by phone
	r.PUT("/users", controller.UpdateUser)        // Update user by username or phone
	r.DELETE("/users", controller.DeleteUser)     // Delete user by username or phone
	r.GET("/users", controller.FindAllUsers)      // Get all users
}