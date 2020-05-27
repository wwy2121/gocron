package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		ctx     context.Context
		err     error
	)
	config = clientv3.Config{
		Endpoints:   []string{"192.168.99.100:2379"},
		DialTimeout: 5 * time.Second,
	}
	client, err = clientv3.New(config)
	if err != nil {
		fmt.Print(err)
		return
	}
	kv = clientv3.NewKV(client)
	ctx, _ = context.WithTimeout(context.TODO(), 5*time.Second)

	//put
	putResp, err = kv.Put(ctx, "/cron/a", "123", clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(putResp.Header.Revision)

	//get
	getResp, err = kv.Get(ctx, "/cron/a", clientv3.WithPrefix())
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(getResp.Kvs)
	client.Close()

}
