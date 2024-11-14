package controller

import (
	"findApi/bootstrap"
	"findApi/domain"
	"findApi/internal/encryptutil"
	"findApi/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type UserController struct {
	UserUsecase usecase.UsersUseCase
	Env         *bootstrap.Env
}

// CreateUser handles the creation of a new user
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user domain.User
	// Parse the request body to get the user details
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	phoneExist,_ := c.UserUsecase.GetUserByPhone(user.Phone)
	usernameExist,_ := c.UserUsecase.GetUserByUsername(user.Username)
	if phoneExist != nil && usernameExist != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone or username already exists"})
		return
	}

	// Call use case to insert the user
	createdUser, err := c.UserUsecase.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user" + err.Error()})
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

// UpdateUser handles updating a user by username, phone, or both
func (c *UserController) UpdateUser(ctx *gin.Context) {
	var req domain.UpdateReq
	// Parse the request body to get the update details
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if req.NewUsername == "" && req.NewPhone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "New username or phone is required"})
		return
	}

	// Initialize a filter for the MongoDB query
	filter := bson.M{}
	if req.Username != "" {
		encUsername, err := encryptutil.EncryptECB(req.Username, []byte(c.Env.SECRET_KEY))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt username"})
			return
		}
		filter["username"] = encUsername
	}
	if req.Phone != "" {
		encPhone, err := encryptutil.EncryptECB(req.Phone, []byte(c.Env.SECRET_KEY))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt phone"})
			return
		}
		filter["phone"] = encPhone
	}

	// Validate that at least one identifier (phone or username) is provided
	if len(filter) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username or Phone is required"})
		return
	}

	updateUser := domain.User{
		Username: req.NewUsername,
		Phone:    req.NewPhone,
	}

	// Call the use case to update the user
	err := c.UserUsecase.UpdateUser(filter, &updateUser)
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
		encUsername, err := encryptutil.EncryptECB(user.Username,[]byte(c.Env.SECRET_KEY))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt username"})
			return
		}
		filter = bson.M{"username": encUsername}
	} else if user.Phone != "" {
		// Encrypt the phone with the same salt/IV used in the database
		encPhone, err := encryptutil.EncryptECB(user.Phone, []byte(c.Env.SECRET_KEY))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt phone"})
			return
		}
		filter = bson.M{"phone": encPhone}
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
