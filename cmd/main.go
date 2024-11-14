package main

import (
	"context"
	"findApi/api/routes"
	"findApi/bootstrap"
	"findApi/repository/db"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	env := bootstrap.LoadEnv()
	router := gin.Default()
	client := db.NewMongoClient(env)
	defer client.Disconnect(context.TODO())

	db := client.Database(env.DB_NAME)

	// Create indexes for the users collection
	if err := createUserIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	routes.SetupRoutes(router, db, env)
	router.Run(":" + env.PORT)
}

// createUserIndexes creates indexes for the users collection on both phone and username fields
func createUserIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := db.Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	}

	// Create indexes
	_, err := usersCollection.Indexes().CreateMany(ctx, indexes)
	return err
}