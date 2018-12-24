package analyzer

import (
	"encoding/csv"
	"encoding/json"
	"io"
)

// Output is an interface for outputting the results of Analyzers. Some implementations are:
// CSVOutput: outputs results in CSV format with header.
// JSONOutput: outputs results in JSON format as an array of objects.
// NoOutput: swallows output. Usually used together with AnalyzerExecutor.ExecuteWithResults().
type Output interface {
	Pre(analyzerWrappers []analyzerWrapper) error
	ReplayResults(results []string) error
	Post() error
}

// NoOutput swallows output. Usually used together with AnalyzerExecutor.ExecuteWithResults().
type NoOutput struct{}

// NewNoOutput is the NoOutput constructor.
func NewNoOutput() *NoOutput { return &NoOutput{} }

// Pre runs at the beginning of the replay analyzing cycle.
func (o *NoOutput) Pre(analyzerWrappers []analyzerWrapper) error { return nil }

// ReplayResults runs at each replay result cycle.
func (o *NoOutput) ReplayResults(results []string) error { return nil }

// Post runs at the end of the replay analyzing cycle.
func (o *NoOutput) Post() error { return nil }

// CSVOutput outputs results in CSV format with header.
type CSVOutput struct {
	w                *csv.Writer
	analyzerWrappers []analyzerWrapper
}

// NewCSVOutput is the CSVOutput constructor.
func NewCSVOutput(w io.Writer) *CSVOutput {
	return &CSVOutput{csv.NewWriter(w), nil}
}

// Pre runs at the beginning of the replay analyzing cycle.
func (o *CSVOutput) Pre(analyzerWrappers []analyzerWrapper) error {
	o.analyzerWrappers = analyzerWrappers
	fieldDisplayNames := []string{}
	for _, wrapper := range analyzerWrappers {
		if !wrapper.isFilter && !wrapper.isFilterNot {
			fieldDisplayNames = append(fieldDisplayNames, wrapper.displayName)
		}
	}
	return o.w.Write(fieldDisplayNames)
}

// ReplayResults runs at each replay result cycle.
func (o *CSVOutput) ReplayResults(_results []string) error {
	results := []string{}
	for i, wrapper := range o.analyzerWrappers {
		if !wrapper.isFilter && !wrapper.isFilterNot {
			results = append(results, _results[i])
		}
	}
	return o.w.Write(results)
}

// Post runs at the end of the replay analyzing cycle.
func (o *CSVOutput) Post() error {
	o.w.Flush()
	return o.w.Error()
}

// JSONOutput outputs results in JSON format as an array of objects.
type JSONOutput struct {
	w                 io.Writer
	firstJSONRow      bool
	analyzerWrappers  []analyzerWrapper
	fieldDisplayNames []string
}

// NewJSONOutput is the JSONOutput constructor.
func NewJSONOutput(w io.Writer) *JSONOutput {
	return &JSONOutput{w, true, nil, nil}
}

// Pre runs at the beginning of the replay analyzing cycle.
func (o *JSONOutput) Pre(analyzerWrappers []analyzerWrapper) error {
	o.analyzerWrappers = analyzerWrappers
	for _, wrapper := range analyzerWrappers {
		o.fieldDisplayNames = append(o.fieldDisplayNames, wrapper.displayName)
	}
	if _, err := o.w.Write([]byte("[\n")); err != nil {
		return err
	}
	return nil
}

// ReplayResults runs at each replay result cycle.
func (o *JSONOutput) ReplayResults(_results []string) error {
	if !o.firstJSONRow {
		if _, err := o.w.Write([]byte(",\n")); err != nil {
			return err
		}
	}
	results := map[string]string{}
	for i, result := range _results {
		if !o.analyzerWrappers[i].isFilter && !o.analyzerWrappers[i].isFilterNot {
			results[o.fieldDisplayNames[i]] = result
		}
	}
	bs, err := json.Marshal(results) // TODO improvement: map traversal
	if err != nil {
		return err
	}
	if _, err := o.w.Write(bs); err != nil {
		return err
	}
	o.firstJSONRow = false
	return nil
}

// Post runs at the end of the replay analyzing cycle.
func (o *JSONOutput) Post() error {
	if _, err := o.w.Write([]byte("\n]\n")); err != nil {
		return err
	}
	return nil
}
