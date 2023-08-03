// Test suite builders for each of the variants.
package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	pb "github.com/greymatter-io/xds-test-harness/api/adapter"
	"github.com/greymatter-io/xds-test-harness/pkg/types"
	"github.com/rs/zerolog/log"

	_ "embed"
)

var (
	//go:embed features/delta.feature
	deltaFeature []byte

	//go:embed features/subscribing.feature
	subscribingFeature []byte

	//go:embed features/unsubcribing.feature
	unsubscribingFeature []byte
)

type Suite struct {
	Variant     types.Variant
	Runner      *Runner
	Aggregated  bool
	Incremental bool
	TestWriting bool
	Buffer      bytes.Buffer
	Tags        string
	TestSuite   godog.TestSuite
}

func (s *Suite) StartRunner(node, adapter, target string) error {
	s.Runner = FreshRunner()
	s.Runner.NodeID = node
	s.Runner.Aggregated = s.Aggregated
	s.Runner.Incremental = s.Incremental

	if err := s.Runner.ConnectClient("target", target); err != nil {
		return fmt.Errorf("cannot connect to target: %v", err)
	}
	if err := s.Runner.ConnectClient("adapter", adapter); err != nil {
		return fmt.Errorf("cannot connect to adapter: %v", err)
	}
	log.Info().
		Msgf("Connected to target at %s and adapter at %s", target, adapter)

	if s.Runner.Aggregated {
		log.Info().
			Msgf("Tests will be run via an aggregated streams")
	} else {
		log.Info().
			Msgf("Tests will be run via separate streams")
	}
	return nil
}

func (s *Suite) SetTags(base string) error {
	tagList := []string{}
	variants := strings.Split(string(s.Variant), " ")

	if len(variants) < 1 {
		err := fmt.Errorf("no variant type found to create tags from. This means the suite was not initialized properly")
		return err
	}
	for _, tag := range variants {
		tag = "@" + tag
		tagList = append(tagList, tag)
	}
	if base != "" {
		tagList = append(tagList, base)
	}
	s.Tags = strings.Join(tagList, " && ")
	return nil
}

func (s *Suite) ConfigureSuite() {
	initScenario := func(ctx *godog.ScenarioContext) {
		ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			log.Debug().
				Msg("Creating Fresh Runner!")
			s.Runner = FreshRunner(s.Runner)
			return ctx, nil
		})
		ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
			if err != nil {
				log.Err(err).Msg("error passed in scenario After hook")
			}
			c := pb.NewAdapterClient(s.Runner.Adapter.Conn)
			clearRequest := &pb.ClearStateRequest{Node: s.Runner.NodeID}
			clear, err := c.ClearState(context.Background(), clearRequest)
			if err != nil {
				log.Err(err).
					Msg("Couldn't clear state")
			}
			log.Debug().
				Msgf("Clearing State: %v\n", clear.Response)
			return ctx, nil
		})
		s.Runner.LoadSteps(ctx)
	}

	godogOpts := godog.Options{
		ShowStepDefinitions: false,
		Randomize:           0,
		StopOnFailure:       false,
		Strict:              false,
		NoColors:            false,
		Tags:                s.Tags,
		Format:              "pretty",
		Concurrency:         0,
	}

	// use our embedded features
	godogOpts.FeatureContents = []godog.Feature{
		{
			Name:     "Delta",
			Contents: deltaFeature,
		},
		{
			Name:     "Subscribing",
			Contents: deltaFeature,
		},
		{
			Name:     "Unsubscribing",
			Contents: deltaFeature,
		},
	}

	if !s.TestWriting { // default is pretty output to stdout.
		// Only use default when writing tests, otherwise print to our special buffer.
		outputFile := variantToOutputFile(s.Variant)
		godogOpts.Format = "xds,cucumber:" + outputFile
		godogOpts.Output = &s.Buffer
	}

	suite := godog.TestSuite{
		Name:                fmt.Sprintf("xds Test Suite [%v]", s.Variant),
		ScenarioInitializer: initScenario,
		Options:             &godogOpts,
	}

	s.TestSuite = suite
}

func (s *Suite) Run(t *testing.T) (results types.VariantResults, err error) {
	s.TestSuite.Options.TestingT = t
	s.TestSuite.Run()
	if s.TestWriting {
		return results, err
	}

	if err = json.Unmarshal(s.Buffer.Bytes(), &results); err != nil {
		err = fmt.Errorf("error unmarshalling test results: %v", err)
		return results, err
	}

	results.Name = string(s.Variant)
	return results, err
}

func NewSotwNonAggregatedSuite(testWriting bool) *Suite {
	return &Suite{
		Variant:     types.SotwNonAggregated,
		Aggregated:  false,
		Incremental: false,
		TestWriting: testWriting,
		Buffer:      *bytes.NewBuffer(nil),
	}
}

func NewSotwAggregatedSuite(testWriting bool) *Suite {
	return &Suite{
		Variant:     types.SotwAggregated,
		Aggregated:  true,
		Incremental: false,
		TestWriting: testWriting,
		Buffer:      *bytes.NewBuffer(nil),
	}

}

func NewIncrementalNonAggregatedSuite(testWriting bool) *Suite {
	return &Suite{
		Variant:     types.IncrementalNonAggregated,
		Aggregated:  false,
		Incremental: true,
		TestWriting: testWriting,
		Buffer:      *bytes.NewBuffer(nil),
	}

}

func NewIncrementalAggregatedSuite(testWriting bool) *Suite {
	return &Suite{
		Variant:     types.IncrementalAggregated,
		Aggregated:  true,
		Incremental: true,
		TestWriting: testWriting,
		Buffer:      *bytes.NewBuffer(nil),
	}
}

func NewSuite(variant types.Variant, testWriting bool) *Suite {
	switch variant {
	case types.SotwNonAggregated:
		return NewSotwNonAggregatedSuite(testWriting)
	case types.SotwAggregated:
		return NewSotwAggregatedSuite(testWriting)
	case types.IncrementalNonAggregated:
		return NewIncrementalNonAggregatedSuite(testWriting)
	case types.IncrementalAggregated:
		return NewIncrementalAggregatedSuite(testWriting)
	default:
		return nil
	}
}

func UpdateResults(current types.Results, variantResults types.VariantResults) types.Results {
	return types.Results{
		Total:            current.Total + int64(variantResults.Total),
		Passed:           current.Passed + int64(variantResults.Passed),
		Failed:           current.Failed + int64(variantResults.Failed),
		Skipped:          current.Skipped + int64(variantResults.Skipped),
		Undefined:        current.Undefined + int64(variantResults.Undefined),
		Pending:          current.Pending + int64(variantResults.Pending),
		Variants:         append(current.Variants, variantResults.Name),
		ResultsByVariant: append(current.ResultsByVariant, variantResults),
	}
}

func variantToOutputFile(v types.Variant) string {
	parts := strings.Split(string(v), " ")
	fileName := strings.Join(parts, "-")
	return fileName + ".json"
}
