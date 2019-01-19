package main

import (
	"reflect"
	"testing"
)

func TestAnalyzers(t *testing.T) {
	ts := []struct {
		name     string
		args     []string
		expected [][]string
	}{
		{
			name: "tests -map-name",
			args: []string{
				"-is-1v1",
				"-is-2v2",
				"-is-there-a-race", "protoss",
				"-map-name",
				"-matchup",
				"-my-game",
				"-my-matchup", // TODO: this should be ZvP; review!
				"-my-name",
				"-my-race",
				"-my-race-is", "zerg",
				"-me", "adultrabbit",
				"-replay", "testdata/larvavsMini.rep", "-o", "none",
			},
			expected: [][]string{{"true", "false", "true", "Transistor1.2", "PvZ", "true", "PvZ", "adultrabbit", "Zerg", "true"}},
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			executor, _, errs := buildAnalyzerExecutor(tc.args)
			if len(errs) != 0 {
				t.Errorf("Expected no errors building AnalyzerExecutor but: %v", errs)
				t.FailNow()
			}
			results, errs := executor.ExecuteWithResults()
			if len(errs) != 0 {
				t.Errorf("Expected no errors executing AnalyzerExecutor but: %v", errs)
				t.FailNow()
			}
			if !reflect.DeepEqual(tc.expected, results) {
				t.Errorf("Expected: %v, but got: %v", tc.expected, results)
				t.FailNow()
			}
		})
	}
}
