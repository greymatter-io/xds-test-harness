package runner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	pb "github.com/greymatter-io/xds-test-harness/api/adapter"
	parser "github.com/greymatter-io/xds-test-harness/pkg/parser"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/anypb"
)

func (r *Runner) LoadSteps(ctx *godog.ScenarioContext) {
	// setting state
	ctx.Step(`^a target setup with service "([^"]*)", resources "([^"]*)", and starting version "([^"]*)"$`, r.TargetSetupWithServiceResourcesAndVersion)
	ctx.Step(`^a target setup with multiple services "([^"]*)", each with resources "([^"]*)", and starting version "([^"]*)"$`, r.TargetSetupWithServiceResourcesAndVersion)
	// client subscriptions
	ctx.Step(`^the Client does a wildcard subscription to "([^"]*)"$`, r.ClientDoesAWildcardSubscriptionToService)
	ctx.Step(`^the Client subscribes to resources "([^"]*)" for "([^"]*)"$`, r.ClientSubscribesToASubsetOfResourcesForService)
	ctx.Step(`^the Client updates subscription to a resource\("([^"]*)"\) of "([^"]*)" with version "([^"]*)"$`, r.ClientUpdatesSubscriptionToAResourceForServiceWithVersion)
	ctx.Step(`^the Client unsubscribes from all resources for "([^"]*)"$`, r.ClientUnsubscribesFromAllResourcesForService)
	ctx.Step(`^the Client unsubscribes from resource "([^"]*)" for service "([^"]*)"$`, r.ClientUnsubscribesFromResourceForService)
	// receiving resources
	ctx.Step(`^the Client receives the resources "([^"]*)" and version "([^"]*)" for "([^"]*)"$`, r.ClientReceivesResourcesAndVersionForService)
	ctx.Step(`^the Client receives only the resource "([^"]*)" and version "([^"]*)" for the service "([^"]*)"$`, r.ClientReceivesOnlyTheResourceAndVersionForTheService)
	ctx.Step(`^the Client does not receive any message from "([^"]*)"$`, r.ClientDoesNotReceiveAnyMessageFromService)
	ctx.Step(`^the Client receives notice that resource "([^"]*)" was removed for service "([^"]*)"$`, r.ClientReceivesNoticeThatResourceWasRemovedForService)
	ctx.Step(`^the client does not receive resource "([^"]*)" of service "([^"]*)" at version "([^"]*)"$`, r.ClientDoesNotReceiveResourceOfServiceAtVersion)
	// resources are added or updated
	ctx.Step(`^the resource "([^"]*)" is added to the "([^"]*)" with version "([^"]*)"$`, r.ResourceIsAddedToServiceWithVersion)
	ctx.Step(`^a resource "([^"]*)" is added to the "([^"]*)" with version "([^"]*)"$`, r.ResourceIsAddedToServiceWithVersion)
	ctx.Step(`^the resources "([^"]*)" are added to the "([^"]*)" with version "([^"]*)"$`, r.ResourceIsAddedToServiceWithVersion)
	ctx.Step(`^the resource "([^"]*)" of service "([^"]*)" is updated to version "([^"]*)"$`, r.ResourceOfServiceIsUpdatedToVersion)
	ctx.Step(`^the resource "([^"]*)" is removed from the "([^"]*)"$`, r.ResourceIsRemovedFromTheService)
	// misc. client server validation
	ctx.Step(`^the service never responds more than necessary$`, r.TheServiceNeverRespondsMoreThanNecessary)
	ctx.Step(`^the resources "([^"]*)" and version "([^"]*)" for "([^"]*)" came in a single response$`, r.ResourcesAndVersionForServiceCameInASingleResponse)
	ctx.Step(`^for service "([^"]*)", no resource other than "([^"]*)" has same version or nonce$`, r.NoOtherResourceHasSameVersionOrNonce)
	ctx.Step(`^for service "([^"]*)", no resource other than "([^"]*)" has same nonce$`, r.NoOtherResourceHasSameNonce)
	ctx.Step(`^the Client sends an ACK to which the "([^"]*)" does not respond$`, r.TheServiceNeverRespondsMoreThanNecessary)
}

///////////////////////////////////////////////////////////////////////////////////
//# Setting State
///////////////////////////////////////////////////////////////////////////////////

// Creates a snapshot to be sent, via the adapter, to the target implementation,
// setting the state for the rest of the steps.
func (r *Runner) TargetSetupWithServiceResourcesAndVersion(services, resources, version string) error {
	resourceNames := strings.Split(resources, ",")
	serviceNames := strings.Split(services, ",")
	anyResources := []*anypb.Any{}

	for _, service := range serviceNames {
		typeUrl, err := parser.ServiceToTypeURL(service)
		if err != nil {
			return err
		}
		for _, name := range resourceNames {
			var any *anypb.Any
			var err error
			switch typeUrl {
			case parser.TypeUrlCDS:
				c := &cluster.Cluster{Name: name}
				any, err = anypb.New(c)
			case parser.TypeUrlLDS:
				l := &listener.Listener{Name: name}
				any, err = anypb.New(l)
			case parser.TypeUrlEDS:
				e := &endpoint.ClusterLoadAssignment{ClusterName: name}
				any, err = anypb.New(e)
			case parser.TypeUrlRDS:
				r := &route.RouteConfiguration{Name: name}
				any, err = anypb.New(r)
			}
			if err != nil {
				return err
			}
			anyResources = append(anyResources, any)
		}

	}
	stateRequest := pb.SetStateRequest{
		Node:      r.NodeID,
		Version:   version,
		Resources: anyResources,
	}

	c := pb.NewAdapterClient(r.Adapter.Conn)

	_, err := c.SetState(context.Background(), &stateRequest)
	if err != nil {
		return fmt.Errorf("cannot set target with given state: %v", err)
	}

	// r.Cache.StartState = snapshot
	return nil
}

///////////////////////////////////////////////////////////////////////////////////
//# Client subscriptions
//////////////////////////////////////////////////////////////////////////////////

// Wrapper to start stream, without resources, for given service
func (r *Runner) ClientDoesAWildcardSubscriptionToService(service string) error {
	resources := []string{}
	err := r.ClientSubscribesToServiceForResources(service, resources)
	return err
}

func (r *Runner) ClientSubscribesToASubsetOfResourcesForService(subset, service string) error {
	resources := strings.Split(subset, ",")
	err := r.ClientSubscribesToServiceForResources(service, resources)
	return err
}

// Takes service and creates a runner.Service with a fresh xDS stream
// for the given service. This is the heart of a test, as it sets up
// the request/response loops that verify the service is working properly.
func (r *Runner) ClientSubscribesToServiceForResources(srv string, resources []string) error {
	typeUrl, err := parser.ServiceToTypeURL(srv)
	if err != nil {
		return err
	}

	r.Validate.Resources[typeUrl] = make(map[string]ValidateResource)
	// initiate a map for delta tests, in case we get any removed resource notifications
	r.Validate.RemovedResources[typeUrl] = make(map[string]ValidateResource)
	for _, resource := range resources {
		r.Validate.Resources[typeUrl][resource] = ValidateResource{}
	}

	// check if we are updating existing stream or starting a new one.
	if (!r.Incremental && r.Service.Sotw != nil) ||
		(r.Incremental && r.Service.Delta != nil) {
		request := r.newRequest(resources, typeUrl)
		r.Service.Channels.Req <- request
		log.Debug().
			Msgf("Sent new subscribing request: %v\n", request)
		return nil
	} else {
		var builder serviceBuilder
		if r.Aggregated {
			builder = getBuilder("ADS")
		} else {
			builder = getBuilder(srv)
		}
		builder.openChannels()
		if r.Incremental {
			err := builder.setDeltaStream(r.Target.Conn)
			if err != nil {
				return err
			}
		} else {
			err := builder.setSotwStream(r.Target.Conn)
			if err != nil {
				return err
			}
		}
		r.Service = builder.getService(srv)
		request := r.newRequest(resources, typeUrl)
		r.SubscribeRequest = request
		log.Debug().
			Msgf("Sending first subscribing request: %v\n", request.String())
		go r.Stream(r.Service)
		go r.Ack(r.Service)
		return nil
	}
}

func (r *Runner) ClientUpdatesSubscriptionToAResourceForServiceWithVersion(resource, service, version string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		err := fmt.Errorf("cannot determine typeURL for given service: %v", service)
		return err
	}

	current := r.Validate.Resources[typeUrl][resource]

	request := &discovery.DiscoveryRequest{
		VersionInfo:   current.Version,
		ResourceNames: []string{resource},
		TypeUrl:       typeUrl,
		ResponseNonce: current.Nonce,
	}
	any, _ := anypb.New(request)

	r.Validate.Resources[typeUrl] = make(map[string]ValidateResource)
	r.Validate.Resources[typeUrl][resource] = ValidateResource{
		Version: current.Version,
		Nonce:   current.Nonce,
	}
	r.SubscribeRequest = any

	log.Debug().Msgf("Sending Request To Update Subscription: %v", request)
	r.Service.Channels.Req <- any
	return nil
}

func (r *Runner) ClientUnsubscribesFromAllResourcesForService(service string) error {
	typeURL, err := parser.ServiceToTypeURL(service)
	if err != nil {
		err := fmt.Errorf("cannot determine typeURL for given service: %v", service)
		return err
	}

	// we just need a nonce to tell the server we are up to dote and this is a new
	// subscription request. Simple way to grab one from the list of 4.
	var lastNonce string
	for _, v := range r.Validate.Resources[typeURL] {
		lastNonce = v.Nonce
	}
	request := &discovery.DiscoveryRequest{
		ResourceNames: []string{""},
		TypeUrl:       typeURL,
		ResponseNonce: lastNonce,
	}
	r.Validate.Resources[typeURL] = make(map[string]ValidateResource)
	any, _ := anypb.New(request)
	r.SubscribeRequest = any
	log.Debug().
		Msgf("Sending unsubscribe request: %v", request.String())
	r.Service.Channels.Req <- any
	return nil
}

// A delta specific test, as delta can explicitly unsubscribe, whereas sotw can only update their subscription
// set up a delta discovery request unsubscribing for given resource, and pass it along the channel.
func (r *Runner) ClientUnsubscribesFromResourceForService(resource, service string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		err := fmt.Errorf("cannot determine typeURL for given service: %v", service)
		return err
	}

	request := &discovery.DeltaDiscoveryRequest{
		TypeUrl:                  typeUrl,
		ResourceNamesUnsubscribe: []string{resource},
	}
	any, _ := anypb.New(request)

	delete(r.Validate.Resources[typeUrl], resource)
	r.SubscribeRequest = any

	log.Debug().Msgf("Sending Unsubscribe Request: %v", request)
	r.Service.Channels.Req <- any
	return nil
}

///////////////////////////////////////////////////////////////////////////////////
//# Receiving resources
///////////////////////////////////////////////////////////////////////////////////

// Loop through the service's response cache until we get the expected response
// or we reach the deadline for the service.
func (r *Runner) ClientReceivesResourcesAndVersionForService(resources, version, service string) error {
	expectedResources := strings.Split(resources, ",")
	stream := r.Service

	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		err := fmt.Errorf("cannot determine typeURL for given service: %v", service)
		return err
	}
	done := time.After(800 * time.Millisecond)
	for {
		select {
		case err := <-stream.Channels.Err:
			// if there isn't an error in a response,
			// the error will be passed down from the stream when
			// it reaches its context deadline.
			return fmt.Errorf("could not find expected response within grace period of 10 seconds. %v", err)
		case <-done:
			actualResources := r.Validate.Resources[typeUrl]
			log.Debug().Msgf("Current resources: %v", r.Validate.Resources)
			// TODO(critter): fork this repo and remove this stuff
			for _, resource := range expectedResources {
				_, ok := actualResources[resource]
				if !ok {
					fmt.Println(actualResources)
					return fmt.Errorf("could not find resource from responses. Expected: %v, Actual: %v", resource, actualResources)
				}
				//if actual.Version != version {
				//return fmt.Errorf("found resource, but not correct version. Expected: %v, Actual: %v", version, actual.Version)
				//}
			}
			return nil
		}
	}
}

// Loop again, but this time continuing if resource in cache has more than one entry.
// The test is itended for when you update a subscription to now only care about a single resource.
// The response you reeceive should only have a single entry in its resources, otherwise we fail.
// Won't work for LDS/CDS where it is conformant to pass along more than you need.
func (r *Runner) ClientReceivesOnlyTheResourceAndVersionForTheService(resource, version, service string) error {
	done := time.After(800 * time.Millisecond)
	for {
		select {
		case err := <-r.Service.Channels.Err:
			return fmt.Errorf("could not find expected response within grace period of 10 seconds or encountered error: %v", err)
		case <-done:
			typeUrl, err := parser.ServiceToTypeURL(service)
			if err != nil {
				return fmt.Errorf("issue converting service to typeUrl, was it written correctly?")
			}
			resources := r.Validate.Resources[typeUrl]
			for name := range resources {
				if name != resource && name != "" {
					return fmt.Errorf("received a resource we should not have. Expected resource: %v. Got: %v",
						resource, name)
				}
			}
			return nil
		}
	}
}
func (r *Runner) ClientDoesNotReceiveAnyMessageFromService(service string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		return err
	}
	done := time.After(800 * time.Millisecond)
	for {
		select {
		case err := <-r.Service.Channels.Err:
			return err
		case <-done:
			if len(r.Validate.Resources[typeUrl]) > 0 {
				return fmt.Errorf("resources received is greater than 0: %v", r.Validate.Resources[typeUrl])
			}
			return nil
		}
	}
}

func (r *Runner) ClientReceivesNoticeThatResourceWasRemovedForService(resource, service string) error {
	stream := r.Service

	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		err := fmt.Errorf("cannot determine typeURL for given service: %v", service)
		return err
	}
	done := time.After(800 * time.Millisecond)
	for {
		select {
		case err := <-stream.Channels.Err:
			return fmt.Errorf("could not find expected response within grace period of 10 seconds. %v", err)
		case <-done:
			actualRemoved := r.Validate.RemovedResources[typeUrl]
			if _, ok := actualRemoved[resource]; !ok {
				return fmt.Errorf("expected resource not in removed resources. Expected: %v, Actual removed: %v", resource, actualRemoved)
			}
			return nil
		}
	}
}

func (r *Runner) ClientDoesNotReceiveResourceOfServiceAtVersion(resource, service, version string) error {
	stream := r.Service
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		err := fmt.Errorf("cannot determine typeURL for given service: %v", service)
		return err
	}
	done := time.After(800 * time.Millisecond)
	for {
		select {
		case err := <-stream.Channels.Err:
			return fmt.Errorf("could not find expected response within grace period of 10 seconds. %v", err)
		case <-done:
			actual := r.Validate.Resources[typeUrl]
			if actual, ok := actual[resource]; ok {
				return fmt.Errorf("was not expecting to find this resource, as we unsubscribed. This is non-conformant: %v", actual)

			}
			return nil
		}
	}
}

///////////////////////////////////////////////////////////////////////////////////
//# Resources are added or updated
///////////////////////////////////////////////////////////////////////////////////

func (r *Runner) ResourceIsAddedToServiceWithVersion(resource, service, version string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		return err
	}

	c := pb.NewAdapterClient(r.Adapter.Conn)
	in := &pb.ResourceRequest{
		Node:         r.NodeID,
		TypeUrl:      typeUrl,
		ResourceName: resource,
		Version:      version,
	}

	_, err = c.AddResource(context.Background(), in)
	if err != nil {
		return fmt.Errorf("cannot add resource using adapter: %v", err)
	}
	log.Debug().
		Msgf("Adding resource %v with version %v", resource, version)
	return nil
}

func (r *Runner) ResourceOfServiceIsUpdatedToVersion(resource, service, version string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		return err
	}

	c := pb.NewAdapterClient(r.Adapter.Conn)
	in := &pb.ResourceRequest{
		Node:         r.NodeID,
		TypeUrl:      typeUrl,
		ResourceName: resource,
		Version:      version,
	}
	log.Debug().
		Msgf("Updating %v resource %v to version %v", service, resource, version)
	_, err = c.UpdateResource(context.Background(), in)
	if err != nil {
		return fmt.Errorf("cannot update resource using adapter: %v", err)
	}
	return nil
}

func (r *Runner) ResourceIsRemovedFromTheService(resource, service string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		return err
	}
	var currentVersion string
	for k, v := range r.Validate.Resources[typeUrl] {
		if k == resource {
			currentVersion = v.Version
		}
	}

	c := pb.NewAdapterClient(r.Adapter.Conn)
	request := &pb.ResourceRequest{
		Node:         r.NodeID,
		TypeUrl:      typeUrl,
		ResourceName: resource,
		Version:      currentVersion,
	}

	_, err = c.RemoveResource(context.Background(), request)
	if err != nil {
		return fmt.Errorf("cannot remove resource using adapter: %v", err)
	}
	log.Debug().
		Msgf("Removing Resource %v", resource)
	return nil
}

///////////////////////////////////////////////////////////////////////////////////
//# Client/server validation
///////////////////////////////////////////////////////////////////////////////////

// ctx.Step(`^the service never responds more than necessary$`, r.TheServiceNeverRespondsMoreThanNecessary)
func (r *Runner) TheServiceNeverRespondsMoreThanNecessary() error {
	stream := r.Service
	stream.Channels.Done <- true

	// give some time for the final messages to come through, if there's any lingering responses.
	time.Sleep(800 * time.Millisecond)
	log.Debug().
		Msgf("Request Count: %v Response Count: %v", r.Validate.RequestCount, r.Validate.ResponseCount)
	if r.Validate.RequestCount <= r.Validate.ResponseCount {
		err := fmt.Errorf("there are more responses than requests.  This indicates the server responded to the last ack")
		return err
	}
	return nil
}

func (r *Runner) ResourcesAndVersionForServiceCameInASingleResponse(resources, version, service string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		return err
	}
	expected := strings.Split(resources, ",")
	actual := r.Validate.Resources[typeUrl]

	responses := make(map[string]bool)
	for _, resource := range expected {
		info, ok := actual[resource]
		// TODO(critter): remove version checks from this thing
		//if !ok || info.Version != version {
		if !ok {
			return fmt.Errorf("could not find correct resource in validation struct. This is rare; perhaps recheck how the test was written")
		}
		responses[info.Nonce] = true
	}
	if len(responses) != 1 {
		return fmt.Errorf("resources came via multiple responses. This is not conformant for CDS and  LDS tests")
	}
	return nil
}

func (r *Runner) NoOtherResourceHasSameVersionOrNonce(service, resource string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		return err
	}
	resources := r.Validate.Resources[typeUrl]
	chosen := resources[resource]
	for r, v := range resources {
		if r == resource {
			continue
		} else if v.Nonce == chosen.Nonce {
			fmt.Printf("%+v\n", resources)
			return fmt.Errorf("found other resource with same nonce, meaning it came back in same response: %v", r)
		} else if v.Version == chosen.Version {
			fmt.Printf("%+v\n", resources)
			return fmt.Errorf("found other resource with same version: %v", r)
		}
	}
	return nil
}

func (r *Runner) NoOtherResourceHasSameNonce(service, resource string) error {
	typeUrl, err := parser.ServiceToTypeURL(service)
	if err != nil {
		return err
	}
	resources := r.Validate.Resources[typeUrl]
	chosen := resources[resource]
	for r, v := range resources {
		if r == resource {
			continue
		} else if v.Nonce == chosen.Nonce {
			return fmt.Errorf("found other resource with same nonce, meaning it came back in same response: %v", r)
		}
	}
	return nil
}
