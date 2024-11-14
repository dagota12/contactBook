package repository

import (
	"context"
	"errors"
	"findApi/bootstrap"
	"findApi/domain"
	"findApi/internal/encryptutil"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var users = make([]*domain.User, 0)
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
		user.Username, _ = encryptutil.DecryptECB(user.Username, []byte(u.SECRET_KEY))
		user.Phone, _ = encryptutil.DecryptECB(user.Phone, []byte(u.SECRET_KEY))
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByUsername retrieves a user by username
func (u *userRepository) GetByUsername(username string) (*domain.User, error) {
	// Encrypt the username for querying
	encUsername, err := encryptutil.EncryptECB(username, []byte(u.SECRET_KEY))
	if err != nil {
		return nil, err
	}

	filter := bson.M{"username": encUsername}
	return u.GetUser(filter)
}

// GetByPhone retrieves a user by phone number
func (u *userRepository) GetByPhone(phone string) (*domain.User, error) {
	// Encrypt the phone for querying
	encPhone, err := encryptutil.EncryptECB(phone, []byte(u.SECRET_KEY))
	if err != nil {
		return nil, err
	}

	filter := bson.M{"phone": encPhone}
	return u.GetUser(filter)
}

// UpdateUser updates a user's details
func (u *userRepository) UpdateUser(filter bson.M, user *domain.User) error {
	// Prepare update data with encryption
	updateData := bson.M{}
	if user.Username != "" {
		encUsername, err := encryptutil.EncryptECB(user.Username, []byte(u.SECRET_KEY))
		if err != nil {
			return err
		}
		updateData["username"] = encUsername
	}

	if user.Phone != "" {
		encPhone, err := encryptutil.EncryptECB(user.Phone, []byte(u.SECRET_KEY))
		if err != nil {
			return err
		}
		updateData["phone"] = encPhone
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
	delRes, err := u.users.DeleteOne(context.TODO(), filter)
	if delRes.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}

// InsertUser adds a new user to the collection
func (u *userRepository) InsertUser(user *domain.User) (*domain.User, error) {
	// Encrypt sensitive fields
	encUsername, err := encryptutil.EncryptECB(user.Username, []byte(u.SECRET_KEY))
	if err != nil {
		return nil, err
	}

	encPhone, err := encryptutil.EncryptECB(user.Phone, []byte(u.SECRET_KEY))
	if err != nil {
		return nil, err
	}

	// Set encrypted values in the user struct
	user.Username = encUsername
	user.Phone = encPhone

	res, err := u.users.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	log.Println(res.InsertedID)
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

// GetUser retrieves a user by a generic filter and decrypts sensitive data
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
	user.Username, _ = encryptutil.DecryptECB(user.Username, []byte(u.SECRET_KEY))
	user.Phone, _ = encryptutil.DecryptECB(user.Phone, []byte(u.SECRET_KEY))

	return &user, nil
}

// NewUserRepository creates a new user repository with collection and secret key
func NewUserRepository(users *mongo.Collection, env *bootstrap.Env) UsersRepo {
	return &userRepository{
		users:      users,
		SECRET_KEY: env.SECRET_KEY,
	}
}
