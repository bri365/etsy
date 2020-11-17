package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/grpclog"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	endpoints      = []string{"192.168.50.181:2379", "192.168.50.182:2379", "192.168.50.183:2379"}
)

func add(c *clientv3.Client, k, v string) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := c.Put(ctx, k, v)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}

func addBig(c *clientv3.Client, keySize, valueSize int) {
	key := make([]byte, keySize)
	rand.Read(key)
	value := make([]byte, valueSize)
	rand.Read(value)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := c.Put(ctx, string(key[:]), string(value[:]))
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	resp, err := c.Get(ctx, string(key[:]))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(resp.Kvs))
}

func main() {
	clientv3.SetLogger(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))

	kvc, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	// make sure to close the client
	defer kvc.Close()

	fmt.Println(time.Now().Clock())

	// takes about 10 seconds
	for i := 0; i < 10; i++ {
		// add(kvc, fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i))
		addBig(kvc, 65540, 32768)
	}

	fmt.Println(time.Now().Clock())
	// time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := kvc.Get(ctx, "key_99")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%#v\n", kvc)
	fmt.Println(resp.Kvs)
}
