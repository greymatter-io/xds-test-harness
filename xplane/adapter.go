package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoyendpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	pb "github.com/greymatter-io/xds-test-harness/api/adapter"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

const (
	configPrefix    = "config"
	discoveryPrefix = "discovery"
)

var _ pb.AdapterServer = &xplaneAdapter{}

type xplaneAdapter struct {
	pb.UnimplementedAdapterServer

	kv nats.KeyValue
}

func (xp *xplaneAdapter) SetState(ctx context.Context, req *pb.SetStateRequest) (*pb.SetStateResponse, error) {
	for _, res := range req.Resources {
		err := xp.sendResourceToNATS(res, "", req.Node)
		if err != nil {
			return nil, err
		}
	}

	resp := &pb.SetStateResponse{
		Success: true,
	}

	return resp, nil
}

func (xp *xplaneAdapter) ClearState(ctx context.Context, req *pb.ClearStateRequest) (*pb.ClearStateResponse, error) {
	keys, err := xp.kv.Keys(nats.Context(ctx))
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if strings.Contains(key, req.Node) {
			err = xp.kv.Purge(key)
			if err != nil {
				return nil, err
			}
		}
	}

	resp := &pb.ClearStateResponse{
		Response: "cleared",
	}

	return resp, nil
}

func (xp *xplaneAdapter) UpdateResource(ctx context.Context, req *pb.ResourceRequest) (*pb.UpdateResourceResponse, error) {
	addResp, err := xp.AddResource(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.UpdateResourceResponse{
		Success: addResp.Success,
	}

	return resp, nil
}

func (xp *xplaneAdapter) AddResource(ctx context.Context, req *pb.ResourceRequest) (*pb.AddResourceResponse, error) {
	resp := &pb.AddResourceResponse{
		Success: false,
	}

	any := &anypb.Any{
		TypeUrl: req.TypeUrl,
	}

	err := xp.sendResourceToNATS(any, req.ResourceName, req.Node)
	if err != nil {
		return resp, err
	}
	resp.Success = true

	return resp, nil
}

func (xp *xplaneAdapter) RemoveResource(ctx context.Context, req *pb.ResourceRequest) (*pb.RemoveResourceResponse, error) {
	resp := &pb.RemoveResourceResponse{
		Success: false,
	}

	var objType string

	switch req.TypeUrl {
	case resource.ClusterType:
		objType = "cluster"
	case resource.ListenerType:
		objType = "listener"
	case resource.RouteType:
		objType = "route"
	case resource.EndpointType:
		// TODO(critter): this should use discoverer
		objType = "clusterloadassignment"
	}

	objKey := fmt.Sprintf("%s.%s.%s.%s", configPrefix, req.Node, objType, req.ResourceName)
	if objType == "clusterloadassignment" {
		objKey = fmt.Sprintf("%s.%s.%s", discoveryPrefix, "xds-test-harness", req.ResourceName)
	}

	fmt.Printf("deleting obj in nats at key %s\n", objKey)

	err := xp.kv.Delete(objKey)
	if err != nil {
		return resp, err
	}
	resp.Success = true

	return resp, nil
}

func (xp *xplaneAdapter) sendResourceToNATS(res *anypb.Any, resourceName, node string) error {
	var xdsType types.Resource
	var objType string
	var err error

	objName := resourceName

	switch res.TypeUrl {
	case resource.ClusterType:
		var obj cluster.Cluster
		err = res.UnmarshalTo(&obj)
		xdsType = &obj
		objType = "cluster"
		if obj.Name == "" {
			obj.Name = objName
		}
		objName = obj.Name
		obj.AltStatName = randStringRunes(10)
	case resource.ListenerType:
		var obj listener.Listener
		err = res.UnmarshalTo(&obj)
		xdsType = &obj
		objType = "listener"
		if obj.Name == "" {
			obj.Name = objName
		}
		objName = obj.Name
		obj.StatPrefix = randStringRunes(10)
	case resource.RouteType:
		var obj route.RouteConfiguration
		err = res.UnmarshalTo(&obj)
		xdsType = &obj
		objType = "route"
		if obj.Name == "" {
			obj.Name = objName
		}
		objName = obj.Name
		obj.MaxDirectResponseBodySizeBytes = &wrapperspb.UInt32Value{Value: rand.Uint32()}
	case resource.EndpointType:
		var obj envoyendpoint.ClusterLoadAssignment
		err = res.UnmarshalTo(&obj)
		xdsType = &obj
		objType = "clusterloadassignment"
		if obj.ClusterName == "" {
			obj.ClusterName = objName
		}
		objName = obj.ClusterName
		obj.Policy = &envoyendpoint.ClusterLoadAssignment_Policy{OverprovisioningFactor: &wrapperspb.UInt32Value{Value: rand.Uint32()}}
	default:
		panic(fmt.Sprintf("unknown/unsupported object type: %s", res.TypeUrl))
	}
	if err != nil {
		return err
	}

	marshed, err := protojson.Marshal(xdsType)
	if err != nil {
		return err
	}

	objKey := fmt.Sprintf("%s.%s.%s.%s", configPrefix, node, objType, objName)

	// discovery expects things in a special place in NATS
	// TODO(critter): this should use a publisher
	if res.TypeUrl == resource.EndpointType {
		objKey = fmt.Sprintf("%s.%s.%s", discoveryPrefix, "xds-test-harness", objName)
	}

	fmt.Printf("putting obj %s in nats at key %s\n", string(marshed), objKey)
	_, err = xp.kv.Put(objKey, marshed)
	if err != nil {
		return err
	}

	return nil
}
