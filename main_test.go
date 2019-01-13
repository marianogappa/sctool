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
			name:     "tests -map-name",
			args:     []string{"-replay", "testdata/larvavsMini.rep", "-map-name", "-o", "none"},
			expected: [][]string{{"Transistor1.2"}},
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
