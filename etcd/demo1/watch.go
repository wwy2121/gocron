package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		//kv clientv3.KV
		watch clientv3.WatchChan
		err   error
	)
	config = clientv3.Config{
		Endpoints:   []string{"192.168.99.100:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			watch = client.Watch(context.TODO(), "/cron/test")
			for val := range watch {
				for k, v := range val.Events {
					fmt.Println(k, string(v.Kv.Key), string(v.Kv.Value))
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(100 * time.Second)

}
