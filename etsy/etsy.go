package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/grpclog"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	endpoints      = []string{"192.168.50.181:2379"}
)

func main() {
	clientv3.SetLogger(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))

	kvc, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer kvc.Close() // make sure to close the client

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = kvc.Put(ctx, "key_1", "value_1")
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	resp, err := kvc.Get(ctx, "key_1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Header)
	fmt.Println(resp.Kvs)
}
