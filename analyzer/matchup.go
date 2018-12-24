package analyzer

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

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
