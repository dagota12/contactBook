package controller

import (
	"findApi/bootstrap"
	"findApi/domain"
	"findApi/internal/encryptutil"
	"findApi/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type UserController struct {
	UserUsecase usecase.UsersUseCase
	Env *bootstrap.Env
}


// CreateUser handles the creation of a new user
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user domain.User
	// Parse the request body to get the user details
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call use case to insert the user
	createdUser, err := c.UserUsecase.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return the created user with a 201 Created status
	ctx.JSON(http.StatusCreated, createdUser)
}

// GetUserByUsername handles fetching a user by username
func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username") // Get the username from the URL path

	// Call use case to fetch user by username
	user, err := c.UserUsecase.GetUserByUsername(username)
	if err != nil || user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the user details with a 200 OK status
	ctx.JSON(http.StatusOK, user)
}

// GetUserByPhone handles fetching a user by phone number
func (c *UserController) GetUserByPhone(ctx *gin.Context) {
	phone := ctx.Param("phone") // Get the phone number from the URL path

	// Call use case to fetch user by phone number
	user, err := c.UserUsecase.GetUserByPhone(phone)
	if err != nil || user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the user details with a 200 OK status
	ctx.JSON(http.StatusOK, user)
}

// UpdateUser handles updating a user by username or phone
func (c *UserController) UpdateUser(ctx *gin.Context) {
	var user domain.User
	// Parse the request body to get the updated user details
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Determine the filter based on the provided username or phone
	var filter bson.M
	if user.Username != "" {
		filter = bson.M{"username": encryptutil.Encrypt(user.Username, c.Env.SECRET_KEY)}
	} else if user.Phone != "" {
		filter = bson.M{"phone": encryptutil.Encrypt(user.Phone, c.Env.SECRET_KEY)}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username or Phone is required"})
		return
	}

	// Call the use case to update the user
	log.Println(user)
	log.Println(filter)
	err := c.UserUsecase.UpdateUser(filter, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Return a success message with a 200 OK status
	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser handles deleting a user by username or phone
func (c *UserController) DeleteUser(ctx *gin.Context) {
	var user domain.User
	// Parse the request body to get the user details
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Determine the filter based on the provided username or phone
	var filter bson.M
	if user.Username != "" {
		filter = bson.M{"username": user.Username}
	} else if user.Phone != "" {
		filter = bson.M{"phone": user.Phone}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username or Phone is required"})
		return
	}

	// Call the use case to delete the user
	err := c.UserUsecase.DeleteUser(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Return a success message with a 200 OK status
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// FindAllUsers handles fetching all users
func (c *UserController) FindAllUsers(ctx *gin.Context) {
	users, err := c.UserUsecase.FindAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	// Return the list of users with a 200 OK status
	ctx.JSON(http.StatusOK, users)
}