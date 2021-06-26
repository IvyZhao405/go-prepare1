package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {

	var (
		ctx context.Context
		cancel context.CancelFunc
		url string
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
	)

	//1. connect to mongodb
	ctx, cancel = context.WithTimeout(context.TODO(), 5*time.Second)
	url = "mongodb://localhost:27017"
	defer cancel()

	if client, err = mongo.Connect(ctx,options.Client().ApplyURI(url)); err != nil {
		fmt.Println(err)
		return
	}
	//2. choose database my_db
	database = client.Database("my_db")
	//3, choose collection my_collection
	collection = database.Collection("my_collection")
	collection = collection
}
