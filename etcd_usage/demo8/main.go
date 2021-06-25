package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var(
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		putOp clientv3.Op
		getOp clientv3.Op
		opResp clientv3.OpResponse
	)
	//etcd client config
	config = clientv3.Config{
		Endpoints: []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)

	//Create Op: operation
	putOp = clientv3.OpPut("/cron/jobs/job8", "123123")

	//execute op
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Write revision:", opResp.Put().Header.Revision)

	//Create op
	getOp = clientv3.OpGet("/cron/jobs/job8")

	//execute get op
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Data revision:", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("Data value:", string(opResp.Get().Kvs[0].Value))

}
