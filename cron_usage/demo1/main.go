package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main(){
	var (
		expr *cronexpr.Expression
		err error
		now time.Time
		nextTime time.Time
	)
	//linux crontab
	//minute (0-59), hour (0-23), day (1-31), month (1-12), week(0-7)
	//cronexpr supports seconds and year (2018- 2099)

	//execute every 5 minute
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}
	//current time
	now = time.Now()

	//next execution time
	nextTime = expr.Next(now)
	fmt.Println(now, nextTime)

	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("crontab working", nextTime)
	})

	time.Sleep(5 * time.Second)
}
