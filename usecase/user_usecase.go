package usecase

import (
	"findApi/domain"
	"findApi/repository"

	"go.mongodb.org/mongo-driver/bson"
)

// UsersUseCase defines the interface for use case operations for managing users.
type UsersUseCase interface {
	// CreateUser adds a new user using either the username or phone number
	CreateUser(user *domain.User) (*domain.User, error)

	// GetUserByUsername retrieves a user by their username
	GetUserByUsername(username string) (*domain.User, error)

	// GetUserByPhone retrieves a user by their phone number
	GetUserByPhone(phone string) (*domain.User, error)

	// UpdateUser updates a user by username or phone
	UpdateUser(filter bson.M, user *domain.User) error

	// DeleteUser deletes a user by username or phone
	DeleteUser(filter bson.M) error

	// FindAllUsers retrieves all users
	FindAllUsers() ([]*domain.User, error)
}











type usersUseCase struct {
	repo repository.UsersRepo
}

// NewUsersUseCase creates a new instance of UsersUseCase with the given repository
func NewUsersUseCase(repo repository.UsersRepo) UsersUseCase {
	return &usersUseCase{
		repo: repo,
	}
}

// CreateUser adds a new user using either the username or phone number
func (u *usersUseCase) CreateUser(user *domain.User) (*domain.User, error) {
	// Validation or additional business logic can be added here
	return u.repo.InsertUser(user)
}

// GetUserByUsername retrieves a user by their username
func (u *usersUseCase) GetUserByUsername(username string) (*domain.User, error) {
	// Get user by username
	return u.repo.GetByUsername(username)
}

// GetUserByPhone retrieves a user by their phone number
func (u *usersUseCase) GetUserByPhone(phone string) (*domain.User, error) {
	// Get user by phone number
	return u.repo.GetByPhone(phone)
}

// UpdateUser updates a user by username or phone
func (u *usersUseCase) UpdateUser(filter bson.M, user *domain.User) error {
	// Business logic for updating the user can be added here (e.g., validating fields)
	return u.repo.UpdateUser(filter, user)
}

// DeleteUser deletes a user by username or phone
func (u *usersUseCase) DeleteUser(filter bson.M) error {
	// Business logic for deleting a user can be added here
	return u.repo.DeleteUser(filter)
}

// FindAllUsers retrieves all users
func (u *usersUseCase) FindAllUsers() ([]*domain.User, error) {
	// Get all users from the repository
	return u.repo.FindAll()
}
