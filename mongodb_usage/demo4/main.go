package main

import (
	"context"
	"fmt"
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

// FindByJobName jobName filter condition
type FindByJobName struct {
	JobName string `bson:"job_name"` //JobName job10
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
		cond *FindByJobName
		record *LogRecord
		cursor *mongo.Cursor
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

	//4, filter based on jobName, find jobName = job10, find 5
	cond = &FindByJobName{JobName: "job10"} //{"jobName": "job10"}

	//5, search
	//findOptions.SetLimit(2)


	if cursor, err = collection.Find(context.TODO(), cond, options.Find().SetSkip(0), options.Find().SetLimit(2)); err != nil {
		fmt.Println(err)
		return
	}

	//delayed release
	defer cursor.Close(context.TODO())

	//6, iterate collection
	for cursor.Next(context.TODO()) {
		record = &LogRecord{}

		//reflect to bson object
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		//print log
		fmt.Println(*record)

	}
}
