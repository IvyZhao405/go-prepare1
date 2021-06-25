package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {

	//use lease to expire lock if node goes down
	//op
	//txn transaction: if else then

	var(
		config clientv3.Config
		client *clientv3.Client
		err error
		lease clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		ctx context.Context
		cancelFunc context.CancelFunc
		kv clientv3.KV
		txn clientv3.Txn
		txnResp *clientv3.TxnResponse
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

	//1. put lock (create lease, auto renew lease, compete for a key with lease)

	//require lease
	lease = clientv3.NewLease(client)

	//register lease for 5s
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	//get lease id
	leaseId = leaseGrantResp.ID

	ctx, cancelFunc = context.WithCancel(context.TODO())

	//renew lease
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	//make sure when func() exit, auto lease renew would stop
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

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

	//if key doesn't exist, then put it, else race for lock failed
	kv = clientv3.NewKV(client)
	//create transaction
	txn = kv.Txn(context.TODO())

	//2. process transaction
	//if key does not exist
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=",0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9")) //then race for lock failed

	//submit transaction
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return // no problem
	}

	//check if we have the lock
	if !txnResp.Succeeded {
		fmt.Println("Lock is occupied:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	//process task if we have the lock
	fmt.Println("Working on task")
	time.Sleep(5 * time.Second)

	//3. release lock (cancel lease auto renew, release lease)
	//defer will cancel lease and revoke lease, related KV will be deleted
}
