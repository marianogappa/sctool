package analyzer

import (
	"fmt"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// -------------------------------------------------------------------------------------------------------------------

// MyFirstSpecificUnitSeconds is used for returning the second that the first specified unit was trained/morphed by
// the -me player. Refer to the unit name list in utils.go#nameToUnitID. -1 if the unit never appears.
type MyFirstSpecificUnitSeconds struct {
	done     bool
	result   string
	playerID byte
	unitID   uint16
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyFirstSpecificUnitSeconds) Name() string { return "my-first-specific-unit-seconds" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyFirstSpecificUnitSeconds) Description() string {
	return "Analyzes the time the first specified unit/building/evolution was built, in seconds."
}

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyFirstSpecificUnitSeconds) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyFirstSpecificUnitSeconds) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyFirstSpecificUnitSeconds) Version() int { return 1 }

// RequiresParsingCommands is true if this Analyzer requires parsing commands from the replay
func (a MyFirstSpecificUnitSeconds) RequiresParsingCommands() bool { return true }

// RequiresParsingMapData is true if this Analyzer requires parsing map data from the replay
func (a MyFirstSpecificUnitSeconds) RequiresParsingMapData() bool { return false }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyFirstSpecificUnitSeconds) Clone() Analyzer {
	return &MyFirstSpecificUnitSeconds{a.done, a.result, a.playerID, a.unitID}
}

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyFirstSpecificUnitSeconds) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyFirstSpecificUnitSeconds) IsStringFlag() bool { return true }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyFirstSpecificUnitSeconds) SetArguments(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a valid unit/building/evolution name e.g. Zergling") // TODO provide list
	}
	if _, ok := nameToUnitID[args[0]]; !ok {
		return fmt.Errorf("invalid unit/building/evolution name") // TODO provide list
	}
	a.unitID = nameToUnitID[args[0]]
	return nil
}

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyFirstSpecificUnitSeconds) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = "-1"
	a.playerID = findPlayerID(replay, ctx.Me)
	a.done = a.playerID == 127 // If we don't find it, no need to see commands
	return a.done, nil
}

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyFirstSpecificUnitSeconds) ProcessCommand(command repcmd.Cmd) (bool, error) {
	a.result, a.done = maybePlayersUnitSeconds(command, a.playerID, a.unitID)
	return a.done, nil
}
