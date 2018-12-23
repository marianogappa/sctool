package analyzer

import (
	"encoding/csv"
	"encoding/json"
	"io"
)

type Output interface {
	Pre(analyzerWrappers []analyzerWrapper) error
	ReplayResults(results []string) error
	Post() error
}

type NoOutput struct{}

func NewNoOutput() *NoOutput                                     { return &NoOutput{} }
func (o *NoOutput) Pre(analyzerWrappers []analyzerWrapper) error { return nil }
func (o *NoOutput) ReplayResults(results []string) error         { return nil }
func (o *NoOutput) Post() error                                  { return nil }

type CSVOutput struct {
	w                *csv.Writer
	analyzerWrappers []analyzerWrapper
}

func NewCSVOutput(w io.Writer) *CSVOutput {
	return &CSVOutput{csv.NewWriter(w), nil}
}

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

func (o *CSVOutput) ReplayResults(_results []string) error {
	results := []string{}
	for i, r := range _results {
		if !o.analyzerWrappers[i].isFilter || !o.analyzerWrappers[i].isFilterNot {
			results = append(results, r)
		}
	}
	return o.w.Write(results)
}

func (o *CSVOutput) Post() error {
	o.w.Flush()
	return o.w.Error()
}

type JSONOutput struct {
	w                 io.Writer
	firstJSONRow      bool
	analyzerWrappers  []analyzerWrapper
	fieldDisplayNames []string
}

func NewJSONOutput(w io.Writer) *JSONOutput {
	return &JSONOutput{w, true, nil, nil}
}

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

func (o *JSONOutput) Post() error {
	if _, err := o.w.Write([]byte("\n]\n")); err != nil {
		return err
	}
	return nil
}
