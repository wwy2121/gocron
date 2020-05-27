package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

func main() {
	var (
		config        clientv3.Config
		client        *clientv3.Client
		kv            clientv3.KV
		ctx           context.Context
		putRepon      *clientv3.PutResponse
		getRespon     *clientv3.GetResponse
		delRespon     *clientv3.DeleteResponse
		grantRespon   *clientv3.LeaseGrantResponse
		lease         clientv3.Lease
		leaseId       clientv3.LeaseID
		keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
		keepAlive     *clientv3.LeaseKeepAliveResponse
		watch         clientv3.WatchChan
		err           error
	)
	//config
	config = clientv3.Config{
		Endpoints:   []string{"192.168.99.100:2379"},
		DialTimeout: 5 * time.Second,
	}
	//client
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	//new
	kv = clientv3.NewKV(client)
	//timeout
	ctx, _ = context.WithTimeout(context.TODO(), 20*time.Second)

	go func() {
		watch = client.Watch(ctx, "/cron/lease")
		for v := range watch {
			for _, val := range v.Events {
				fmt.Println("动态变化值", val.Type, string(val.Kv.Key), val.PrevKv)
			}
		}
	}()

	//put
	if putRepon, err = kv.Put(ctx, "/cron/c", "345545"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(putRepon.Header)
	//get
	if getRespon, err = kv.Get(ctx, "/cron", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(getRespon.Kvs)
	for k, v := range getRespon.Kvs {
		fmt.Println(k, string(v.Key), string(v.Value))
	}
	//del
	if delRespon, err = kv.Delete(ctx, "/cron/c"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(delRespon.Deleted)
	//租约
	lease = clientv3.NewLease(client)
	if grantRespon, err = lease.Grant(ctx, 10); err != nil {
		fmt.Println(err)
		return
	}
	leaseId = grantRespon.ID
	if putRepon, err = kv.Put(ctx, "/cron/lease", "4535354", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	if keepAliveChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(keepAliveChan)
	time.Sleep(30 * time.Second)

	//续租
	go func() {
		for {
			select {
			case keepAlive = <-keepAliveChan:
				if keepAlive == nil {
					fmt.Println("过期了")
					goto END
				} else {
					fmt.Println("续租成功", keepAlive.ID)
				}
			}
			time.Sleep(1 * time.Second)
		}
	END:
	}()
	fmt.Println(putRepon.Header)
	time.Sleep(100 * time.Second)

}
