package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"sync"
	"time"
)

func main() {
	var (
		config        clientv3.Config
		client        *clientv3.Client
		kv            clientv3.KV
		err           error
		lease         clientv3.Lease
		grantRespon   *clientv3.LeaseGrantResponse
		leaseId       clientv3.LeaseID
		putRespon     *clientv3.PutResponse
		keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
		keepAlive     *clientv3.LeaseKeepAliveResponse
		ctx           context.Context
		cancelFunc    context.CancelFunc
		wg            sync.WaitGroup
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
	if putRespon, err = kv.Put(context.TODO(), "/cron/lease", "1111", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(putRespon)
	ctx, cancelFunc = context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
		lease.Revoke(ctx, leaseId)
	})
	keepAliveChan, err = lease.KeepAlive(ctx, leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}
	//开启协程自动续约
	wg.Add(1)
	go func() {
		for {
			select {
			case keepAlive = <-keepAliveChan:
				if keepAlive == nil {
					fmt.Println("过期了")
					goto END
				} else {
					fmt.Println("续约成功:", keepAlive.ID)
				}
			}
			time.Sleep(1 * time.Second)
		}
	END:
		wg.Done()
	}()
	wg.Wait()
}
