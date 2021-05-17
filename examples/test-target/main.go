package main

import (
	"context"
	"flag"
	// "fmt"
	// "net"
	// "log"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	logrus "github.com/sirupsen/logrus"

	// pb "github.com/zachmandeville/tester-prototype/test-target/test-target"
	"github.com/zachmandeville/tester-prototype/examples/test-target/internal/processor"
	"github.com/zachmandeville/tester-prototype/examples/test-target/internal/server"
	"github.com/zachmandeville/tester-prototype/examples/test-target/internal/shim"
	// "google.golang.org/grpc"
)

var (
	nodeID   string
	l        logrus.FieldLogger
	shimPort string
	port     uint
	proc *processor.Processor
)

// type shimServer struct {
// 	pb.UnimplementedShimServer
// }

// func (s *shimServer) GiveCompliment(ctx context.Context, req *pb.ComplimentRequest) (res *pb.ComplimentResponse, err error) {
// 	name := req.Name
// 	compliment := &pb.ComplimentResponse{
// 		Compliment: fmt.Sprintf("Hi, %v, you are GREAT!", name),
// 	}
// 	return compliment, nil
// }

// func (s *shimServer) AddCluster(ctx context.Context, req *pb.ClusterRequest) (res *pb.ClusterResponse, err error) {
// 	cluster := req.Cluster
// 	proc.UpdateSnapshot(cluster)
// 	response := &pb.ClusterResponse{
// 		Message: fmt.Sprintf("A cluster named %v was added\n",cluster),
// 	}
// 	return response, nil
// }


func init() {
	l = logrus.New()
	logrus.SetLevel(logrus.DebugLevel)
	flag.UintVar(&port, "port", 18000, "xDS management server port")
	flag.StringVar(&nodeID, "nodeID", "test-id", "NodeID")
	flag.StringVar(&shimPort, "shimPort", ":17000", "port of xds server shim")
}

func main() {
	flag.Parse()
	cache := cache.NewSnapshotCache(false, cache.IDHash{}, l)
	proc = processor.NewProcessor(
		cache, nodeID, logrus.WithField("context", "processor"))
	go func() {
		// Run the xDS server
		ctx := context.Background()
		srv := serverv3.NewServer(ctx, cache, nil)
		server.RunServer(ctx, srv, port)
	}()

	shim.RunServer(proc, shimPort)
	// run the shim
	// lis, err := net.Listen("tcp", shimPort)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()
	// pb.RegisterShimServer(s, &shimServer{})
	// log.Printf("shim server listening on port %v\n", shimPort)
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("Failed to serve: %v", err)
	// }
}
