package report

import (
	"io/ioutil"
	"testing"

	"github.com/sclevine/spec"
)

// Log reports specs via the testing log methods and only affects verbose runs.
type Log struct{}

func (Log) Start(t *testing.T, plan spec.Plan) {
	t.Helper()
	t.Log("Suite:", plan.Text)
	t.Logf("Total: %d | Focused: %d | Pending: %d", plan.Total, plan.Focused, plan.Pending)
	if plan.HasRandom {
		t.Log("Random seed:", plan.Seed)
	}
	if plan.HasFocus {
		t.Log("Focus is active.")
	}
}

func (Log) Specs(t *testing.T, specs <-chan spec.Spec) {
	t.Helper()
	var passed, failed, skipped int
	for s := range specs {
		switch {
		case s.Failed:
			failed++
			if testing.Verbose() {
				if out, err := ioutil.ReadAll(s.Out); err == nil {
					t.Logf("%s", out)
				}
			}
		case s.Skipped:
			skipped++
		default:
			passed++
		}
	}
	t.Logf("Passed: %d | Failed: %d | Skipped: %d", passed, failed, skipped)
}
