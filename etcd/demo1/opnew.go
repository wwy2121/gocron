package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

func main() {
	var (
		config      clientv3.Config
		client      *clientv3.Client
		kv          clientv3.KV
		lease       clientv3.Lease
		grantRespon *clientv3.LeaseGrantResponse
		leaseId     clientv3.LeaseID
		ctx         context.Context
		cancelFunc  context.CancelFunc
		txn         clientv3.Txn
		txnRespon   *clientv3.TxnResponse
		err         error
	)
	config = clientv3.Config{
		Endpoints:   []string{"192.168.99.100:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	if grantRespon, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	leaseId = grantRespon.ID
	ctx, cancelFunc = context.WithCancel(context.TODO())
	defer cancelFunc()
	defer lease.Revoke(ctx, leaseId)
	//开启事务
	txn = kv.Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/mmm"), "=", 0)).
		Then(clientv3.OpPut("/cron/mmm", "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/mmm"))
	if txnRespon, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}
	if !txnRespon.Succeeded {
		fmt.Println("锁已经被抢了")
		return
	}
	fmt.Println("处理业务逻辑")
	time.Sleep(10 * time.Second)
}
