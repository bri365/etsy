package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	endpoints      = []string{"192.168.50.181:2379", "192.168.50.182:2379", "192.168.50.183:2379"}

	serverAddr = "192.168.50.181:2379"
	hostName   = "etcd1.rvidiot.io"
	keyFile    = "client-key.pem"
	certFile   = "client.pem"
	caFile     = "ca.pem"
)

func main() {
	clientCert, err := tls.LoadX509KeyPair("client.pem", "client-key.pem")
	if err != nil {
		log.Fatalf("Failed to load client certificate and key. %s.", err)
	}

	caCert, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		log.Fatalf("Failed to load CA certificate. %s.", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to add CA certificate to pool. %s.", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
	}

	cred := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(cred), grpc.WithBlock()}
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(cred)}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to dial %v", err)
	}

	defer conn.Close()
	client := NewKVClient(conn)

	// Put
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	pq := PutRequest{
		Key:   []byte("abc"),
		Value: []byte("def"),
	}
	pr, err := client.Put(ctx, &pq)
	cancel()
	if err != nil {
		log.Fatalf("Put failed %v", err)
	}
	fmt.Printf("Put %d\n", pr.Header.Revision)

	// Put
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	pq = PutRequest{
		Key:   []byte("abcd"),
		Value: []byte("efgh"),
	}
	pr, err = client.Put(ctx, &pq)
	cancel()
	if err != nil {
		log.Fatalf("Put failed %v", err)
	}
	fmt.Printf("Put %d\n", pr.Header.Revision)

	// Range
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	rq := RangeRequest{
		Key:      []byte("abc"),
		RangeEnd: []byte("b"),
	}
	rr, err := client.Range(ctx, &rq)
	cancel()
	if err != nil {
		log.Fatalf("Range failed %v", err)
	}
	fmt.Printf("Range %v\n", rr.Kvs)
}
