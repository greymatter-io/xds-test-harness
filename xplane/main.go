package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/greymatter-io/xds-test-harness/api/adapter"
	"github.com/nats-io/nats.go"

	"google.golang.org/grpc"
)

var bucketName = "config_tests"

func main() {
	conn, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		panic(err)
	}

	js, err := conn.JetStream()
	if err != nil {
		panic(err)
	}

	cfg := &nats.KeyValueConfig{
		Bucket: bucketName,
	}

	kv, err := js.CreateKeyValue(cfg)
	if err != nil {
		panic(err)
	}

	xpa := &xplaneAdapter{
		kv: kv,
	}
	port := 17000

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAdapterServer(s, xpa)
	log.Printf("Testsuite Adapter listening on %v\n", port)
	if err := s.Serve(lis); err != nil {
		log.Printf("Adapter failed to serve: %v", err)
	}

	log.Println("deleting bucket:", bucketName)
	err = js.DeleteKeyValue(bucketName)
	if err != nil {
		panic(err)
	}
}
