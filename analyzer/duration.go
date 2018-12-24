package analyzer

import (
	"fmt"
	"strconv"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// -------------------------------------------------------------------------------------------------------------------

// DurationMinutes returns the duration of the replay in minutes.
type DurationMinutes struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a DurationMinutes) Name() string { return "duration-minutes" }

// Description is used for flg usages
func (a DurationMinutes) Description() string {
	return "Analyzes the duration of the replay in minutes."
}

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a DurationMinutes) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a DurationMinutes) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a DurationMinutes) Version() int { return 1 }

// RequiresParsingCommands is true if this Analyzer requires parsing commands from the replay
func (a DurationMinutes) RequiresParsingCommands() bool { return false }

// RequiresParsingMapData is true if this Analyzer requires parsing map data from the replay
func (a DurationMinutes) RequiresParsingMapData() bool { return false }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a DurationMinutes) Clone() Analyzer { return &DurationMinutes{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a DurationMinutes) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a DurationMinutes) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutes) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutes) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutes) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = fmt.Sprintf("%v", int(replay.Header.Duration().Minutes()))
	a.done = true
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// DurationMinutesIsGreaterThan returns true if the replay's duration is greater than the specified amount of minutes.
type DurationMinutesIsGreaterThan struct {
	done    bool
	result  string
	minutes int
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a DurationMinutesIsGreaterThan) Name() string { return "duration-minutes-is-greater-than" }

// Description is ued for flag usages
func (a DurationMinutesIsGreaterThan) Description() string {
	return "Analyzes if the duration of the replay in minutes is greater than specified."
}

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a DurationMinutesIsGreaterThan) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a DurationMinutesIsGreaterThan) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a DurationMinutesIsGreaterThan) Version() int { return 1 }

// RequiresParsingCommands is true if this Analyzer requires parsing commands from the replay
func (a DurationMinutesIsGreaterThan) RequiresParsingCommands() bool { return false }

// RequiresParsingMapData is true if this Analyzer requires parsing map data from the replay
func (a DurationMinutesIsGreaterThan) RequiresParsingMapData() bool { return false }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a DurationMinutesIsGreaterThan) Clone() Analyzer {
	return &DurationMinutesIsGreaterThan{a.done, a.result, a.minutes}
}

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a DurationMinutesIsGreaterThan) IsBooleanResult() bool { return true }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a DurationMinutesIsGreaterThan) IsStringFlag() bool { return true }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutesIsGreaterThan) SetArguments(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a valid number of minutes")
	}
	n, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid number of minutes: %v", args[0])
	}
	a.minutes = n
	return nil
}

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutesIsGreaterThan) ProcessCommand(command repcmd.Cmd) (bool, error) {
	return true, nil
}

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutesIsGreaterThan) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	actualMinutes := int(replay.Header.Duration().Minutes())
	a.result = fmt.Sprintf("%v", actualMinutes > a.minutes)
	a.done = true
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// DurationMinutesIsLowerThan returns true if the replay's duration is lower than the specified amount of minutes.
type DurationMinutesIsLowerThan struct {
	done    bool
	result  string
	minutes int
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a DurationMinutesIsLowerThan) Name() string { return "duration-minutes-is-lower-than" }

// Description is used for flag usage
func (a DurationMinutesIsLowerThan) Description() string {
	return "Analyzes if the duration of the replay in minutes is lower than specified."
}

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a DurationMinutesIsLowerThan) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a DurationMinutesIsLowerThan) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a DurationMinutesIsLowerThan) Version() int { return 1 }

// RequiresParsingCommands is true if this Analyzer requires parsing commands from the replay
func (a DurationMinutesIsLowerThan) RequiresParsingCommands() bool { return false }

// RequiresParsingMapData is true if this Analyzer requires parsing map data from the replay
func (a DurationMinutesIsLowerThan) RequiresParsingMapData() bool { return false }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a DurationMinutesIsLowerThan) Clone() Analyzer {
	return &DurationMinutesIsLowerThan{a.done, a.result, a.minutes}
}

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a DurationMinutesIsLowerThan) IsBooleanResult() bool { return true }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a DurationMinutesIsLowerThan) IsStringFlag() bool { return true }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutesIsLowerThan) SetArguments(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a valid number of minutes")
	}
	n, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid number of minutes: %v", args[0])
	}
	a.minutes = n
	return nil
}

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutesIsLowerThan) ProcessCommand(command repcmd.Cmd) (bool, error) {
	return true, nil
}

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DurationMinutesIsLowerThan) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	actualMinutes := int(replay.Header.Duration().Minutes())
	a.result = fmt.Sprintf("%v", actualMinutes < a.minutes)
	a.done = true
	return a.done, nil
}
