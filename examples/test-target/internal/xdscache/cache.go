package xdscache

import (
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/golang/protobuf/ptypes"

	pb "github.com/zachmandeville/tester-prototype/api/adapter"
)

type XDSCache struct {
	Listeners map[string]Listener
	Routes    map[string]Route
	Clusters  map[string]*cluster.Cluster
	Endpoints map[string]Endpoint
}

func (xds *XDSCache) ClusterContents() []types.Resource {
	var r []types.Resource
	for _, c := range xds.Clusters {
		r = append(r, c)
	}
	return r
}

func (xds *XDSCache) AddCluster(c *pb.Clusters_Cluster) {
	seconds := time.Duration(c.ConnectTimeout["seconds"])
	xds.Clusters[c.Name] = &cluster.Cluster{
		Name:           c.Name,
		ConnectTimeout: ptypes.DurationProto(seconds * time.Second),
	}
}
