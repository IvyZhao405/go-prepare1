package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var(
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		deleteResp *clientv3.DeleteResponse
		kvpair *mvccpb.KeyValue
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

	//delete kv
	if deleteResp, err = kv.Delete(context.TODO(), "/cron/jobs/job2", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	}
	//the kv value before delete
	if len(deleteResp.PrevKvs) != 0 {
		for _, kvpair = range deleteResp.PrevKvs{
			fmt.Println("Deleted:", string(kvpair.Key), string(kvpair.Value))
		}
	}
}
