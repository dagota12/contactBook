package repository

import (
	"context"
	"errors"
	"findApi/bootstrap"
	"findApi/domain"
	"findApi/internal/encryptutil"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo interface {
	InsertUser(user *domain.User) (*domain.User, error)
	GetUser(filter bson.M) (*domain.User, error)
	GetByPhone(phone string) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	UpdateUser(filter bson.M, user *domain.User) error
	DeleteUser(filter bson.M) error
	FindAll() ([]*domain.User, error)
}

type userRepository struct {
	users      *mongo.Collection
	SECRET_KEY string
}

// FindAll retrieves all users from the collection
func (u *userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := u.users.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}

		// Decrypt user data before returning
		user.Username = encryptutil.Decrypt(user.Username, u.SECRET_KEY)
		user.Phone = encryptutil.Decrypt(user.Phone, u.SECRET_KEY)
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByUsername retrieves a user by username
func (u *userRepository) GetByUsername(username string) (*domain.User, error) {
	filter := bson.M{"username": encryptutil.Encrypt(username, u.SECRET_KEY)}
	return u.GetUser(filter)
}

// GetByPhone retrieves a user by phone number
func (u *userRepository) GetByPhone(phone string) (*domain.User, error) {
	filter := bson.M{"phone": encryptutil.Encrypt(phone, u.SECRET_KEY)}
	return u.GetUser(filter)
}

// UpdateUser updates a user's details
func (u *userRepository) UpdateUser(filter bson.M, user *domain.User) error {
	// Prepare update data with encryption
	updateData := bson.M{}
	if user.Username != "" {
		updateData["username"] = encryptutil.Encrypt(user.Username, u.SECRET_KEY)
	}
	if user.Phone != "" {
		updateData["phone"] = encryptutil.Encrypt(user.Phone, u.SECRET_KEY)
	}

	// Perform the update
	updateRes, err := u.users.UpdateOne(context.TODO(), filter, bson.M{"$set": updateData})
	if updateRes.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}

// DeleteUser deletes a user by filter
func (u *userRepository) DeleteUser(filter bson.M) error {
	_, err := u.users.DeleteOne(context.TODO(), filter)
	return err
}

// InsertUser adds a new user to the collection
func (u *userRepository) InsertUser(user *domain.User) (*domain.User, error) {
	// Encrypt sensitive fields
	user.Username = encryptutil.Encrypt(user.Username, u.SECRET_KEY)
	user.Phone = encryptutil.Encrypt(user.Phone, u.SECRET_KEY)

	_, err := u.users.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by a generic filter
func (u *userRepository) GetUser(filter bson.M) (*domain.User, error) {
	var user domain.User
	err := u.users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No user found
		}
		return nil, err
	}

	// Decrypt sensitive data
	user.Username = encryptutil.Decrypt(user.Username, u.SECRET_KEY)
	user.Phone = encryptutil.Decrypt(user.Phone, u.SECRET_KEY)

	return &user, nil
}

// NewUserRepository creates a new user repository with collection and secret key
func NewUserRepository(users *mongo.Collection, env *bootstrap.Env) UsersRepo {
	return &userRepository{
		users:      users,
		SECRET_KEY: env.SECRET_KEY,
	}
}
