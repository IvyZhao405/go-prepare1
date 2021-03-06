package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName string `bson:"job_name"` //job name
	Command string `bson:"command"`//command
	Err string `bson:"err"`//error
	Content string `bson:"content"`// script output
	TimePoint TimePoint `bson:"timePoint"`// executing time
}


func main (){

	var (
		ctx context.Context
		cancel context.CancelFunc
		url string
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		record *LogRecord
		result *mongo.InsertOneResult
		docId primitive.ObjectID
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

	//4, insert record
	record = &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println(err)
		return
	}
	//_id:  default create a unique ID, objectID: 12byte binary
	docId = result.InsertedID.(primitive.ObjectID)
	fmt.Print("Self Increasing ID:", docId.Hex())

}
