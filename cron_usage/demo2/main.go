package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr *cronexpr.Expression
	nextTime time.Time //expr:Next(now)
}
func main() {
	//need one dispatcher coroutine to check all Cron jobs, whoever expired execute it.

	var (
		cronJob *CronJob
		expr *cronexpr.Expression
		now time.Time
		scheduleTable map[string]*CronJob //key: job name, value: CronJob struct
	)

	scheduleTable = make(map[string]*CronJob)

	//current time
	now = time.Now()

	//1, define 2 Cron jobs
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr: expr,
		nextTime: expr.Next(now),
	}
	//add  job to schedule table
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr: expr,
		nextTime: expr.Next(now),
	}
	//add  job to schedule table
	scheduleTable["job2"] = cronJob

	//start dispatcher coroutine
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)
		//check cronjob
		for {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				//check if expired
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//start a coroutine to execute job
					go func(jobName string) {
						fmt.Println("Executing", jobName)
					}(jobName)

					//calculate next time
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "Next executing time:", cronJob.nextTime)
				}
			}
			//sleep 100 ms
			select {
			case <-time.NewTimer(100 * time.Millisecond).C:
			}
			//time.Sleep(100 * time.Millisecond)
		}
	}()
	time.Sleep(100*time.Second)
}
