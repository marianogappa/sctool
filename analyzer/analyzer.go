package analyzer

import (
	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// An Analyzer is structure code that, given a replay, it determines something
// about it. It is optimized for modularity, extensibility and performance.
// An example would be an Analyzer that answers if the game is a 1v1, or
// if a player did a 5rax BO.
type Analyzer interface {
	// Analyzer name: used for DependsOn() and as argument to CLI, so must be
	// hyphenated and without spaces or special characters.
	Name() string

	// Human readable description for what the analyzer is useful for. Used
	// in command line usage help.
	Description() string

	// Arguments for running: should be called before StartReadingReplay().
	SetArguments(args []string)

	// Analyzer Name's whose Results this Analyzer depends on: for building DAG.
	DependsOn() map[string]struct{}

	// Called at the beginning of a Replay analizing cycle.
	// Should be called before any ProcessCommand and only if IsDone is false.
	// Returns true if the analyzer is finished calculating the result (i.e. no need
	// to process commands)
	StartReadingReplay(replay *rep.Replay) bool

	// Should be called for every command during a Replay analizing cycle.
	// StartReadingReplay should be called before processing any command, to refresh
	// any state and to decide if processing commands are necessary to determine result.
	// Returns true if the analyzer is finished calculating the result (i.e. no need
	// to process further commands)
	ProcessCommand(command repcmd.Cmd) bool

	// Returns true if the analyzer is finished calculating the result, and
	// returns it. Shouldn't be called before calling StartReadingReplay.
	IsDone() (Result, bool)

	// Useful for managing updates to an Analyzer: whenever an update is made to an
	// analyzer, the Version should be numerically higher. Then, if there's a cached
	// Result of an Analyzer on a Replay, the result should be recomputed.
	Version() int
}

// Result of an Analyzer
type Result interface {
	Value() string
}

type stringResult struct {
	result string
}

func (r stringResult) Value() string { return r.result }
