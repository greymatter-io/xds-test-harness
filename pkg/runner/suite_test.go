package runner

import (
	"testing"
	"github.com/greymatter-io/xds-test-harness/pkg/types"
)

func TestSetTags (t *testing.T) {
	// type of suite don't matter, this is just convenient
	tags := "@sotw && @non-aggregated"
	base := "@wip"
	suite := NewSuite(types.SotwNonAggregated, true)
	suite.SetTags(base)
	expected := tags + " && " + base
	if suite.Tags != expected {
	  t.Errorf("Created tags not matching what is expected. Expected: %v, Actual: %v", expected, suite.Tags)
	}
}
