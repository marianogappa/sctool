package analyzer

import (
	"fmt"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// -------------------------------------------------------------------------------------------------------------------

// IsThereARace is "true" if the specified race appears in the replay.
type IsThereARace struct {
	done   bool
	race   string
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a IsThereARace) Name() string { return "is-there-a-race" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a IsThereARace) Description() string {
	return "Analyzes if there is a specific race in the replay."
}

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a IsThereARace) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a IsThereARace) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a IsThereARace) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a IsThereARace) Clone() Analyzer { return &IsThereARace{a.done, a.race, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a IsThereARace) IsBooleanResult() bool { return true }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a IsThereARace) IsStringFlag() bool { return true }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *IsThereARace) SetArguments(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a valid race name e.g. Zerg/Protoss/Terran") // TODO provide list
	}
	r := strings.ToLower(args[0])
	if _, ok := raceNameTranslations[r]; !ok {
		return fmt.Errorf("invalid race name %v", args[0]) // TODO provide list
	}
	a.race = raceNameTranslations[r]
	return nil
}

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *IsThereARace) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *IsThereARace) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = "false"
	for _, p := range replay.Header.OrigPlayers {
		if p.Race.Name == a.race {
			a.result = "true"
		}
	}
	a.done = true
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MyAPM returns the APM of the -me player. -1 if the player is not in the replay.
type MyAPM struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyAPM) Name() string { return "my-apm" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyAPM) Description() string { return "Analyzes the APM of the -me player." }

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyAPM) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyAPM) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyAPM) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyAPM) Clone() Analyzer { return &MyAPM{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyAPM) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyAPM) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyAPM) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyAPM) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyAPM) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = "-1"
	if replay.Computed == nil {
		a.done = true
		return true, nil
	}
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		return true, fmt.Errorf("-me player not present in this replay")
	}
	for _, pDesc := range replay.Computed.PlayerDescs {
		if pDesc.PlayerID == playerID {
			a.result = fmt.Sprintf("%v", pDesc.APM)
			a.done = true
			break
		}
	}
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MyRace returns the race of the -me player. Either Zerg, Terran or Protoss. "" if there's no -me player.
type MyRace struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyRace) Name() string { return "my-race" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyRace) Description() string { return "Analyzes the race of the -me player." }

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyRace) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyRace) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyRace) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyRace) Clone() Analyzer { return &MyRace{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyRace) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyRace) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyRace) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyRace) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyRace) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = ""
	a.done = true
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		return true, fmt.Errorf("-me player not present in this replay")
	}
	a.result = replay.Header.PIDPlayers[playerID].Race.Name
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MyRaceIs returns true if the -me player has the specified race.
type MyRaceIs struct {
	done   bool
	result string
	race   string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyRaceIs) Name() string { return "my-race-is" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyRaceIs) Description() string {
	return "Analyzes if the race of the -me player is the one specified."
}

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyRaceIs) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyRaceIs) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyRaceIs) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyRaceIs) Clone() Analyzer { return &MyRaceIs{a.done, a.result, a.race} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyRaceIs) IsBooleanResult() bool { return true }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyRaceIs) IsStringFlag() bool { return true }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyRaceIs) SetArguments(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a valid race name e.g. Zerg/Protoss/Terran") // TODO provide list
	}
	r := strings.ToLower(args[0])
	if _, ok := raceNameTranslations[r]; !ok {
		return fmt.Errorf("invalid race name %v", args[0]) // TODO provide list
	}
	a.race = raceNameTranslations[r]
	return nil
}

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyRaceIs) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyRaceIs) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = ""
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		a.done = true
		return true, fmt.Errorf("-me player not present in this replay")
	}
	a.done = true
	a.result = "false"
	if replay.Header.PIDPlayers[playerID].Race.Name == a.race {
		a.result = "true"
	}
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// DateTime returns the DateTime in which the replay was played.
type DateTime struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a DateTime) Name() string { return "date-time" }

// Description is used fr flag usages
func (a DateTime) Description() string { return "Analyzes the datetime of the replay." }

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a DateTime) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a DateTime) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a DateTime) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a DateTime) Clone() Analyzer { return &DateTime{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a DateTime) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a DateTime) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DateTime) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DateTime) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *DateTime) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = fmt.Sprintf("%v", replay.Header.StartTime)
	a.done = true
	return a.done, nil
}

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

// -------------------------------------------------------------------------------------------------------------------

// MyName returns the name of the -me player. Note that -me might contain many names, but MyName will be only one.
type MyName struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyName) Name() string { return "my-name" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyName) Description() string { return "Analyzes the name of the -me player." }

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyName) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyName) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyName) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyName) Clone() Analyzer { return &MyName{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyName) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyName) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyName) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyName) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyName) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = ""
	a.done = true
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		return true, fmt.Errorf("-me player not present in this replay")
	}
	a.result = replay.Header.PIDPlayers[playerID].Name
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// ReplayName returns the name of the replay file. Note that two replays might have the same name on a different path;
// if you want an unique identifier ReplayPath is better suited (but an md5 would be best).
type ReplayName struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a ReplayName) Name() string { return "replay-name" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a ReplayName) Description() string { return "Analyzes the replay's name." }

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a ReplayName) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a ReplayName) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a ReplayName) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a ReplayName) Clone() Analyzer { return &ReplayName{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a ReplayName) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a ReplayName) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *ReplayName) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *ReplayName) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *ReplayName) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = path.Base(replayPath)
	a.result = a.result[:len(a.result)-4]
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// ReplayPath returns the path to the replay. Prefer this to ReplayName as an unique identifier.
type ReplayPath struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a ReplayPath) Name() string { return "replay-path" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a ReplayPath) Description() string { return "Analyzes the replay's path." }

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a ReplayPath) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a ReplayPath) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a ReplayPath) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a ReplayPath) Clone() Analyzer { return &ReplayPath{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a ReplayPath) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a ReplayPath) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *ReplayPath) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *ReplayPath) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *ReplayPath) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = replayPath
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MyWin returns true if the -me player won the game. Note that, unfortunately, calculating the winner of a replay is
// either quite error-prone or impossible, so this analyzer is not very useful.
type MyWin struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyWin) Name() string { return "my-win" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyWin) Description() string { return "Analyzes if the -me player won the game." }

// DependsOn is used to determine the Analyzer dependency DAG
// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyWin) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyWin) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyWin) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyWin) Clone() Analyzer { return &MyWin{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyWin) IsBooleanResult() bool { return true }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyWin) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyWin) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyWin) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyWin) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.done = true
	if replay.Computed == nil || replay.Computed.WinnerTeam == 0 {
		a.result = "unknown"
		return true, nil
	}
	a.result = "false"
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		return true, fmt.Errorf("-me player not present in this replay")
	}
	if replay.Header.PIDPlayers[playerID].Team == replay.Computed.WinnerTeam {
		a.result = "true"
	}
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MyGame returns true if the -me player participated in this game.
type MyGame struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyGame) Name() string { return "my-game" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyGame) Description() string { return "Analyzes if the -me player played the game." }

// DependsOn is use to determine the Analyzer dependency DAG
// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyGame) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyGame) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyGame) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyGame) Clone() Analyzer { return &MyGame{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyGame) IsBooleanResult() bool { return true }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyGame) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyGame) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyGame) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyGame) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.done = true
	if playerID := findPlayerID(replay, ctx.Me); playerID == 127 {
		a.result = "false"
		return true, nil
	}
	a.result = "true"
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MapName returns the map's name. Note that it doesn't do anything clever, so many versions of a map can have
// slightly different names, or two maps with the same name might be actually different.
type MapName struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MapName) Name() string { return "map-name" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MapName) Description() string { return "Analyzes the map's name." }

// DependsOn is used to determine the Analyzer dependecy DAG
// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MapName) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MapName) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MapName) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MapName) Clone() Analyzer { return &MapName{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MapName) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MapName) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MapName) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MapName) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MapName) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.result = replay.Header.Map
	a.done = true
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// Matchup returns the matchup of the game. On an 1v1, it will sort the races lexicographically, so it will return
// TvZ rather than ZvT. Other than 1v1, it will simply return whatever screp returns.
type Matchup struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a Matchup) Name() string { return "matchup" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a Matchup) Description() string { return "Analyzes the replay's matchup." }

// DependsOn is used to determine the Analyzer dependency DA
// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a Matchup) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a Matchup) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a Matchup) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a Matchup) Clone() Analyzer { return &Matchup{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a Matchup) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a Matchup) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *Matchup) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *Matchup) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *Matchup) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.done = true
	if len(replay.Header.Players) == 2 {
		r0 := strings.ToUpper(string(replay.Header.Players[0].Race.Letter))
		r1 := strings.ToUpper(string(replay.Header.Players[1].Race.Letter))
		a.result = r0 + "v" + r1
		if r0 > r1 {
			a.result = r1 + "v" + r0
		}
	} else {
		a.result = replay.Header.Matchup()
	}
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MyMatchup returns the replay's matchup from the point of view of the -me player. For example, if the -me player is
// Z and the opponent is T it will return ZvT rather than TvZ. At the moment, the behaviour other than 1v1 is
// unexpected: it returns the matchup as returned by screp.
type MyMatchup struct {
	done   bool
	result string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyMatchup) Name() string { return "my-matchup" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyMatchup) Description() string {
	return "Analyzes the replay's matchup from the point of view of the -me player"
}

// DependsOn is used to determine the Aalyzer dependency DAG
// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyMatchup) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyMatchup) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyMatchup) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyMatchup) Clone() Analyzer { return &MyMatchup{a.done, a.result} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyMatchup) IsBooleanResult() bool { return false }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyMatchup) IsStringFlag() bool { return false }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyMatchup) SetArguments(args []string) error { return nil }

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyMatchup) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyMatchup) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.done = true
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		return true, fmt.Errorf("-me player not present in this replay")
	}
	if len(replay.Header.Players) == 2 {
		r0 := strings.ToUpper(string(replay.Header.Players[0].Race.Letter))
		r1 := strings.ToUpper(string(replay.Header.Players[1].Race.Letter))
		a.result = r0 + "v" + r1
		if playerID == 1 {
			a.result = r1 + "v" + r0
		}
	} else {
		a.result = replay.Header.Matchup() // TODO put -me player on the left side
	}
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MatchupIs returns true if the matchup is equal to the specified one. Note it only works for 1v1 for now.
// The specified matchup can be in either order (i.e. ZvT == TvZ).
type MatchupIs struct {
	done   bool
	result string
	races  []string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MatchupIs) Name() string { return "matchup-is" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MatchupIs) Description() string {
	return "Analyzes if the replay's MatchupIs is equal to the specified one (only works for 1v1 for now)."
}

// DependsOn is used to determine he Analyzer dependency DAG
// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MatchupIs) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MatchupIs) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MatchupIs) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MatchupIs) Clone() Analyzer { return &MatchupIs{a.done, a.result, cloneStringSlice(a.races)} }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MatchupIs) IsBooleanResult() bool { return true }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MatchupIs) IsStringFlag() bool { return true }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MatchupIs) SetArguments(args []string) error {
	if len(args) < 1 || len(args[0]) != 3 {
		return fmt.Errorf("please provide a valid matchup e.g. TvZ (only works for 1v1 for now)")
	}
	args[0] = strings.ToUpper(args[0])
	a.races = append(a.races, string(args[0][0]), string(args[0][2]))
	sort.Strings(a.races)
	return nil
}

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MatchupIs) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MatchupIs) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.done = true
	a.result = "false"
	if len(replay.Header.Players) != 2 {
		return true, nil
	}
	actualRaces := []string{
		strings.ToUpper(string(replay.Header.Players[0].Race.Letter)),
		strings.ToUpper(string(replay.Header.Players[1].Race.Letter)),
	}
	sort.Strings(actualRaces)
	a.result = fmt.Sprintf("%v", reflect.DeepEqual(a.races, actualRaces))
	return a.done, nil
}

// -------------------------------------------------------------------------------------------------------------------

// MyMatchupIs returns true if the matchup is equal to the specified one, from the point of view of the -me player.
// Note it only works for 1v1 for now. The specified matchup must contain the -me player's race first.
type MyMatchupIs struct {
	done   bool
	result string
	races  []string
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a MyMatchupIs) Name() string { return "my-matchup-is" }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a MyMatchupIs) Description() string {
	return "Analyzes if the replay's matchup is equal to the specified one, from the -me player perspective (only works for 1v1 for now)."
}

// DependsOn is used to determine theAnalyzer dependency DAG
// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a MyMatchupIs) DependsOn() map[string]struct{} { return map[string]struct{}{} }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a MyMatchupIs) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a MyMatchupIs) Version() int { return 1 }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a MyMatchupIs) Clone() Analyzer {
	return &MyMatchupIs{a.done, a.result, cloneStringSlice(a.races)}
}

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a MyMatchupIs) IsBooleanResult() bool { return true }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a MyMatchupIs) IsStringFlag() bool { return true }

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyMatchupIs) SetArguments(args []string) error {
	if len(args) < 1 || len(args[0]) != 3 {
		return fmt.Errorf("please provide a valid matchup e.g. TvZ (only works for 1v1 for now)")
	}
	args[0] = strings.ToUpper(args[0])
	a.races = append(a.races, string(args[0][0]), string(args[0][2]))
	return nil
}

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyMatchupIs) ProcessCommand(command repcmd.Cmd) (bool, error) { return true, nil }

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *MyMatchupIs) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	a.done = true
	a.result = "false"
	if len(replay.Header.Players) != 2 {
		return true, nil
	}
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		return true, fmt.Errorf("-me player not present in this replay")
	}
	actualRaces := []string{
		strings.ToUpper(string(replay.Header.Players[0].Race.Letter)),
		strings.ToUpper(string(replay.Header.Players[1].Race.Letter)),
	}
	if playerID == 1 {
		actualRaces[0], actualRaces[1] = actualRaces[1], actualRaces[0]
	}
	a.result = fmt.Sprintf("%v", reflect.DeepEqual(a.races, actualRaces))
	return a.done, nil
}
