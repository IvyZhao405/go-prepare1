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
		lease clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
		kv clientv3.KV
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

	//require lease
	lease = clientv3.NewLease(client)

	//lease for 10s
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	//get lease id
	leaseId = leaseGrantResp.ID

	//ctx, _ := context.WithTimeout(context.TODO(), 5* time.Second)

	//lease 5 seconds, stop lease, then 10 s left. = 15s in total
	//renew lease
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}

	//renew lease response coroutine
	go func(){
		for {
			select {
			case keepResp = <- keepRespChan:
				if keepRespChan == nil {
					fmt.Println("lease expired")
					goto END
				}else { //renew every second, so a response every second
					fmt.Println("renew successful", keepResp.ID)
				}
			}
		}
		END:
	}()

	//get kv client
	kv = clientv3.NewKV(client)
	//put KV using lease
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
	}
	fmt.Println("write success:", putResp.Header.Revision)

	//check if lease expired
	for {
		if getResp, err = kv.Get(context.TODO(),"/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv expired")
			break
		}
		fmt.Println("not expired yet", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}
