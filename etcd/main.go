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
		err    error
		client *clientv3.Client
	)
	config = clientv3.Config{
		Endpoints:   []string{"192.168.99.100:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	putResp, err := client.Put(context.TODO(), "/job/v5", "push the box", clientv3.WithPrevKV()) //withPrevKV()是为了获取操作前已经有的key-value
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(putResp.Header.Revision)
	delResp, err := client.Delete(context.TODO(), "/job", clientv3.WithPrefix())
	fmt.Printf("%v", delResp.Header.Revision)
	defer client.Close()

}
