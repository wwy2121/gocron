package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

func main() {
	const KEY string = "/cron/abc"
	var (
		config          clientv3.Config
		client          *clientv3.Client
		kv              clientv3.KV
		putRespon       *clientv3.PutResponse
		getRespon       *clientv3.GetResponse
		delRespon       *clientv3.DeleteResponse
		lease           clientv3.Lease
		grantRespon     *clientv3.LeaseGrantResponse
		leaseId         clientv3.LeaseID
		aliveResponChan <-chan *clientv3.LeaseKeepAliveResponse
		aliveRespon     *clientv3.LeaseKeepAliveResponse
		ctx             context.Context
		cancelFunc      context.CancelFunc
		err             error
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
	if putRespon, err = kv.Put(context.TODO(), KEY, "5345345"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(putRespon.Header.Revision)
	if getRespon, err = kv.Get(context.TODO(), KEY); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(getRespon.Kvs[0].Value))
	if delRespon, err = kv.Delete(context.TODO(), KEY); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(delRespon.Deleted)

	lease = clientv3.NewLease(client)
	if grantRespon, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	leaseId = grantRespon.ID
	fmt.Println(leaseId)

	ctx, cancelFunc = context.WithCancel(context.TODO())
	defer cancelFunc()
	defer lease.Revoke(ctx, leaseId)

	if _, err = kv.Put(context.TODO(), "/cron/lease", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	if aliveResponChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			select {
			case aliveRespon = <-aliveResponChan:
				if aliveRespon == nil {
					fmt.Println("过期了")
					goto END
				} else {
					fmt.Println("续租成功")
				}
			}
			time.Sleep(1 * time.Second)
		}
	END:
	}()

	time.Sleep(20 * time.Second)
}
