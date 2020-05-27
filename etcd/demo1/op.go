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
		putOp       clientv3.Op
		getOp       clientv3.Op
		delOp       clientv3.Op
		opRespon    clientv3.OpResponse
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
		DialTimeout: 10 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)
	putOp = clientv3.OpPut("/cron/opput/", "53535411")
	if opRespon, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(opRespon.Put().Header.Revision)
	getOp = clientv3.OpGet("/cron/opput/")
	if opRespon, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(opRespon.Get().Kvs)
	delOp = clientv3.OpDelete("/cron/opput/")
	if opRespon, err = kv.Do(context.TODO(), delOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(opRespon.Del().Deleted)
	//申请租约
	lease = clientv3.NewLease(client)
	if grantRespon, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	//租约id
	leaseId = grantRespon.ID
	ctx, cancelFunc = context.WithCancel(context.TODO())

	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	//开启事务
	txn = kv.Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock", "locked", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock"))

	//提交事务
	if txnRespon, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	//是否抢到了锁
	if !txnRespon.Succeeded {
		fmt.Println("锁已经被抢走了")
		return
	}
	//处理业务
	fmt.Println("hedadasd")
	time.Sleep(5 * time.Second)
}
