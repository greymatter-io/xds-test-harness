package parser

import (
	"fmt"
	"strings"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/greymatter.io/xds-test-harness/internal/types"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/rs/zerolog/log"
)

const (
	TypeUrlLDS = "type.googleapis.com/envoy.config.listener.v3.Listener"
	TypeUrlCDS = "type.googleapis.com/envoy.config.cluster.v3.Cluster"
	TypeUrlRDS = "type.googleapis.com/envoy.config.route.v3.RouteConfiguration"
	TypeUrlEDS = "type.googleapis.com/envoy.config.endpoint.v3.ClusterLoadAssignment"
)

func ServiceToTypeURL(service string) (typeURL string, err error) {
	typeURLs := map[string]string{
		"lds": TypeUrlLDS,
		"cds": TypeUrlCDS,
		"eds": TypeUrlEDS,
		"rds": TypeUrlRDS,
	}
	service = strings.ToLower(service)

	typeURL, ok := typeURLs[service]
	if !ok {
		err = fmt.Errorf("cannot find type URL for given service: %v", service)
		return typeURL, err
	}
	return typeURL, err
}

func ResourceNames(res *envoy_service_discovery_v3.DiscoveryResponse) (resourceNames []string, err error) {
	typeUrl := res.TypeUrl
	switch typeUrl {
	case TypeUrlLDS:
		for _, resource := range res.GetResources() {
			listener := &listener.Listener{}
			if err := resource.UnmarshalTo(listener); err != nil {
				return nil, fmt.Errorf("could not get resource name from %v. err: %v", resource, err)
			}
			resourceNames = append(resourceNames, listener.Name)
		}
	case TypeUrlCDS:
		for _, resource := range res.GetResources() {
			cluster := &cluster.Cluster{}
			if err := resource.UnmarshalTo(cluster); err != nil {
				return nil, fmt.Errorf("could not get resource name from %v. err: %v", resource, err)
			}
			resourceNames = append(resourceNames, cluster.Name)
		}
	case TypeUrlEDS:
		for _, resource := range res.GetResources() {
			endpointConfig := &endpoint.ClusterLoadAssignment{}
			if err := resource.UnmarshalTo(endpointConfig); err != nil {
				return nil, fmt.Errorf("could not get resource name from %v. err: %v", resource, err)
			}
			resourceNames = append(resourceNames, endpointConfig.ClusterName)
		}
	case TypeUrlRDS:
		for _, resource := range res.GetResources() {
			route := &route.RouteConfiguration{}
			if err := resource.UnmarshalTo(route); err != nil {
				return nil, fmt.Errorf("could not get resource name from %v. err: %v", resource, err)
			}
			resourceNames = append(resourceNames, route.Name)
		}
	}
	return resourceNames, err
}

func ParseSupportedVariants(variants []string) (supported []types.Variant, err error) {
	variantMap := map[string]types.Variant{
		"sotw non-aggregated":        types.SotwNonAggregated,
		"sotw aggregated":            types.SotwAggregated,
		"incremental non-aggregated": types.IncrementalNonAggregated,
		"incremental aggregated":     types.IncrementalAggregated,
	}

	for _, v := range variants {
		variant, ok := variantMap[strings.ToLower(v)]
		if !ok {
			err = fmt.Errorf("config included unrecognized variant. Please remove it and try again: %v", variant)
			return supported, err
		}
		supported = append(supported, variant)
	}
	return supported, err
}

func ValuesFromConfig(config string) (target, adapter, nodeID string, supportedVariants []types.Variant) {
	c, err := yaml.ReadFile(config)
	if err != nil {
		log.Fatal().
			Msgf("cannot read config: %v", config)
	}
	nodeID, err = c.Get("nodeID")
	if err != nil {
		log.Fatal().
			Msgf("error reading config file for Node ID: %v\n", err)
	}
	target, err = c.Get("targetAddress")
	if err != nil {
		log.Fatal().
			Msgf("error reading config file for Target Address: %v\n", config)
	}
	adapter, err = c.Get("adapterAddress")
	if err != nil {
		log.Info().
			Msgf("cannot get adapter address from config file: %v\n", err)
	}
	v, err := yaml.Child(c.Root, "variants")
	if err != nil {
		log.Fatal().Msgf("error getting variants from config: %v\n", err)
	}
	variants := []string{}
	varsInYaml, ok := v.(yaml.List)
	if ok {
		for i := 0; i < varsInYaml.Len(); i++ {
			node := varsInYaml.Item(i)
			variant := string(node.(yaml.Scalar))
			variants = append(variants, variant)
		}
	}
	supportedVariants, err = ParseSupportedVariants(variants)
	if err != nil {
		log.Fatal().Msgf("Cannot parse supported variants from config: %v", err)
	}
	return target, adapter, nodeID, supportedVariants
}
