package db

import (
	"context"
	"findApi/bootstrap"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(env *bootstrap.Env) *mongo.Client{

	options := options.Client().ApplyURI(env.MONGO_URI)
	client,err := mongo.Connect(context.TODO(),options)


	if err != nil{
		log.Fatal("Faild to connect to db:" + err.Error())
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil{
		log.Fatal("Faild to ping db:" + err.Error())
	}

	log.Println("Connected to DB")

	return client


}

