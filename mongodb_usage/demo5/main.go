package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//start time less than
//{"$lt": timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

//{"timePoint.startTime":{"$lt":timestamp}}
type DeleteCond struct {
	beforeCond TimeBeforeCond`bson:"timePoint.startTime"`
}

func main() {
	var (
		ctx context.Context
		cancel context.CancelFunc
		url string
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		delCond *DeleteCond
		delResult *mongo.DeleteResult
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
	database = client.Database("cron")
	//3, choose collection my_collection
	collection = database.Collection("log")

	//4, delete all logs that's earlier than current time($lt less than)
	//delete({"timePoint.startTime":{"$lt:current time}})
	delCond = &DeleteCond{beforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	//execute delete
	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Deleted lines:", delResult.DeletedCount)
}
