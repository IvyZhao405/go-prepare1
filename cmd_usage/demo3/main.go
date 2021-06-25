package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct{
	err error
	output []byte
}

func main() {
	//execute 1 cmd, let it run in a Coroutine, sleep 2; echo hello;

	//at 1 second, we kill this cmd
	var(
		ctx context.Context
		cancelFunc context.CancelFunc
		cmd *exec.Cmd
		resultChan chan *result
		res *result
	)

	//create result array
	resultChan = make(chan*result, 1000)
	//context: chan byte (channel)
	//cancelFunc: close(chan byte)

	ctx, cancelFunc = context.WithCancel(context.TODO())
	go func() {
		var (
			output []byte
			err error
		)
		cmd = exec.CommandContext(ctx,"/bin/bash", "-c", "sleep 2;echo hello;")

		output, err = cmd.CombinedOutput()
		//select {case < ctx.Done(): }
		//kill pid, kill process by id

		//send function result to main coroutine
		resultChan <- &result{
			err: err,
			output: output,
		}
	}()

	//sleep 1 s then continue
	time.Sleep(1 * time.Second)

	//cancel context
	cancelFunc()
	//wait for child coroutine stop, print result in main coroutine
	res = <- resultChan

	//print output
	fmt.Println(res.err, string(res.output))
}
