package analyzer

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// -------------------------------------------------------------------------------------------------------------------
type IsThereARace struct {
	done   bool
	race   string
	result string
}

func (a IsThereARace) Name() string { return "is-there-a-race" }
func (a IsThereARace) Description() string {
	return "Analyzes if there is a specific race in the replay."
}
func (a IsThereARace) DependsOn() map[string]struct{} { return map[string]struct{}{} }
func (a IsThereARace) IsDone() (string, bool)         { return a.result, a.done }
func (a IsThereARace) Version() int                   { return 1 }
func (a IsThereARace) IsBooleanResult() bool          { return true }
func (a IsThereARace) IsStringFlag() bool             { return true }
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
func (a *IsThereARace) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *IsThereARace) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = "false"
	for _, p := range replay.Header.OrigPlayers {
		if p.Race.Name == a.race {
			a.result = "true"
		}
	}
	a.done = true
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MyAPM struct {
	done   bool
	result string
}

func (a MyAPM) Name() string                                     { return "my-apm" }
func (a MyAPM) Description() string                              { return "Analyzes the APM of the -me player." }
func (a MyAPM) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a MyAPM) IsDone() (string, bool)                           { return a.result, a.done }
func (a MyAPM) Version() int                                     { return 1 }
func (a MyAPM) IsBooleanResult() bool                            { return false }
func (a MyAPM) IsStringFlag() bool                               { return false }
func (a *MyAPM) SetArguments(args []string) error                { return nil }
func (a *MyAPM) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *MyAPM) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	if replay.Computed == nil {
		a.result = "-1"
		a.done = true
		return nil, true
	}
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		return fmt.Errorf("-me player not present in this replay"), true
	}
	for _, pDesc := range replay.Computed.PlayerDescs {
		if pDesc.PlayerID == playerID {
			a.result = fmt.Sprintf("%v", pDesc.APM)
			a.done = true
			break
		}
	}
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MyRace struct {
	done   bool
	result string
}

func (a MyRace) Name() string                                     { return "my-race" }
func (a MyRace) Description() string                              { return "Analyzes the race of the -me player." }
func (a MyRace) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a MyRace) IsDone() (string, bool)                           { return a.result, a.done }
func (a MyRace) Version() int                                     { return 1 }
func (a MyRace) IsBooleanResult() bool                            { return false }
func (a MyRace) IsStringFlag() bool                               { return false }
func (a *MyRace) SetArguments(args []string) error                { return nil }
func (a *MyRace) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *MyRace) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = ""
	a.done = true
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		return fmt.Errorf("-me player not present in this replay"), true
	}
	a.result = replay.Header.PIDPlayers[playerID].Race.Name
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MyRaceIs struct {
	done   bool
	result string
	race   string
}

func (a MyRaceIs) Name() string { return "my-race-is" }
func (a MyRaceIs) Description() string {
	return "Analyzes if the race of the -me player is the one specified."
}
func (a MyRaceIs) DependsOn() map[string]struct{} { return map[string]struct{}{} }
func (a MyRaceIs) IsDone() (string, bool)         { return a.result, a.done }
func (a MyRaceIs) Version() int                   { return 1 }
func (a MyRaceIs) IsBooleanResult() bool          { return true }
func (a MyRaceIs) IsStringFlag() bool             { return true }
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
func (a *MyRaceIs) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *MyRaceIs) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = ""
	playerID := findPlayerID(replay, ctx.Me)
	if playerID == 127 {
		a.done = true
		return fmt.Errorf("-me player not present in this replay"), true
	}
	a.done = true
	a.result = "false"
	if replay.Header.PIDPlayers[playerID].Race.Name == a.race {
		a.result = "true"
	}
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type DateTime struct {
	done   bool
	result string
}

func (a DateTime) Name() string                                     { return "date-time" }
func (a DateTime) Description() string                              { return "Analyzes the datetime of the replay." }
func (a DateTime) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a DateTime) IsDone() (string, bool)                           { return a.result, a.done }
func (a DateTime) Version() int                                     { return 1 }
func (a DateTime) IsBooleanResult() bool                            { return false }
func (a DateTime) IsStringFlag() bool                               { return false }
func (a *DateTime) SetArguments(args []string) error                { return nil }
func (a *DateTime) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *DateTime) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = fmt.Sprintf("%v", replay.Header.StartTime)
	a.done = true
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type DurationMinutes struct {
	done   bool
	result string
}

func (a DurationMinutes) Name() string { return "duration-minutes" }
func (a DurationMinutes) Description() string {
	return "Analyzes the duration of the replay in minutes."
}
func (a DurationMinutes) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a DurationMinutes) IsDone() (string, bool)                           { return a.result, a.done }
func (a DurationMinutes) Version() int                                     { return 1 }
func (a DurationMinutes) IsBooleanResult() bool                            { return false }
func (a DurationMinutes) IsStringFlag() bool                               { return false }
func (a *DurationMinutes) SetArguments(args []string) error                { return nil }
func (a *DurationMinutes) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *DurationMinutes) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = fmt.Sprintf("%v", int(replay.Header.Duration().Minutes()))
	a.done = true
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type DurationMinutesIsGreaterThan struct {
	done    bool
	result  string
	minutes int
}

func (a DurationMinutesIsGreaterThan) Name() string { return "duration-minutes-is-greater-than" }
func (a DurationMinutesIsGreaterThan) Description() string {
	return "Analyzes if the duration of the replay in minutes is greater than specified."
}
func (a DurationMinutesIsGreaterThan) DependsOn() map[string]struct{} { return map[string]struct{}{} }
func (a DurationMinutesIsGreaterThan) IsDone() (string, bool)         { return a.result, a.done }
func (a DurationMinutesIsGreaterThan) Version() int                   { return 1 }
func (a DurationMinutesIsGreaterThan) IsBooleanResult() bool          { return true }
func (a DurationMinutesIsGreaterThan) IsStringFlag() bool             { return true }
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
func (a *DurationMinutesIsGreaterThan) ProcessCommand(command repcmd.Cmd) (error, bool) {
	return nil, true
}
func (a *DurationMinutesIsGreaterThan) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	actualMinutes := int(replay.Header.Duration().Minutes())
	a.result = fmt.Sprintf("%v", actualMinutes > a.minutes)
	a.done = true
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type DurationMinutesIsLowerThan struct {
	done    bool
	result  string
	minutes int
}

func (a DurationMinutesIsLowerThan) Name() string { return "duration-minutes-is-lower-than" }
func (a DurationMinutesIsLowerThan) Description() string {
	return "Analyzes if the duration of the replay in minutes is lower than specified."
}
func (a DurationMinutesIsLowerThan) DependsOn() map[string]struct{} { return map[string]struct{}{} }
func (a DurationMinutesIsLowerThan) IsDone() (string, bool)         { return a.result, a.done }
func (a DurationMinutesIsLowerThan) Version() int                   { return 1 }
func (a DurationMinutesIsLowerThan) IsBooleanResult() bool          { return true }
func (a DurationMinutesIsLowerThan) IsStringFlag() bool             { return true }
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
func (a *DurationMinutesIsLowerThan) ProcessCommand(command repcmd.Cmd) (error, bool) {
	return nil, true
}
func (a *DurationMinutesIsLowerThan) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	actualMinutes := int(replay.Header.Duration().Minutes())
	a.result = fmt.Sprintf("%v", actualMinutes < a.minutes)
	a.done = true
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MyName struct {
	done   bool
	result string
}

func (a MyName) Name() string                                     { return "my-name" }
func (a MyName) Description() string                              { return "Analyzes the name of the -me player." }
func (a MyName) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a MyName) IsDone() (string, bool)                           { return a.result, a.done }
func (a MyName) Version() int                                     { return 1 }
func (a MyName) IsBooleanResult() bool                            { return false }
func (a MyName) IsStringFlag() bool                               { return false }
func (a *MyName) SetArguments(args []string) error                { return nil }
func (a *MyName) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *MyName) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = ""
	if replay.Computed == nil {
		a.done = true
		return nil, true
	}
	for _, p := range replay.Header.Players {
		if _, ok := ctx.Me[p.Name]; ok {
			a.result = p.Name
			a.done = true
			break
		}
	}
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type ReplayName struct {
	done   bool
	result string
}

func (a ReplayName) Name() string                                     { return "replay-name" }
func (a ReplayName) Description() string                              { return "Analyzes the replay's name." }
func (a ReplayName) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a ReplayName) IsDone() (string, bool)                           { return a.result, a.done }
func (a ReplayName) Version() int                                     { return 1 }
func (a ReplayName) IsBooleanResult() bool                            { return false }
func (a ReplayName) IsStringFlag() bool                               { return false }
func (a *ReplayName) SetArguments(args []string) error                { return nil }
func (a *ReplayName) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *ReplayName) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = path.Base(replayPath)
	a.result = a.result[:len(a.result)-4]
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type ReplayPath struct {
	done   bool
	result string
}

func (a ReplayPath) Name() string                                     { return "replay-path" }
func (a ReplayPath) Description() string                              { return "Analyzes the replay's path." }
func (a ReplayPath) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a ReplayPath) IsDone() (string, bool)                           { return a.result, a.done }
func (a ReplayPath) Version() int                                     { return 1 }
func (a ReplayPath) IsBooleanResult() bool                            { return false }
func (a ReplayPath) IsStringFlag() bool                               { return false }
func (a *ReplayPath) SetArguments(args []string) error                { return nil }
func (a *ReplayPath) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *ReplayPath) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = replayPath
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MyWin struct {
	done   bool
	result string
}

func (a MyWin) Name() string                                     { return "my-win" }
func (a MyWin) Description() string                              { return "Analyzes if the -me player won the game." }
func (a MyWin) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a MyWin) IsDone() (string, bool)                           { return a.result, a.done }
func (a MyWin) Version() int                                     { return 1 }
func (a MyWin) IsBooleanResult() bool                            { return true }
func (a MyWin) IsStringFlag() bool                               { return false }
func (a *MyWin) SetArguments(args []string) error                { return nil }
func (a *MyWin) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *MyWin) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	if replay.Computed == nil || replay.Computed.WinnerTeam == 0 {
		a.result = "unknown"
		a.done = true
		return nil, true
	}
	a.result = "false"
	for _, p := range replay.Header.Players {
		if _, ok := ctx.Me[p.Name]; ok && p.Team == replay.Computed.WinnerTeam {
			a.result = "true"
			break
		}
	}
	a.done = true
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MyGame struct {
	done   bool
	result string
}

func (a MyGame) Name() string                                     { return "my-game" }
func (a MyGame) Description() string                              { return "Analyzes if the -me player played the game." }
func (a MyGame) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a MyGame) IsDone() (string, bool)                           { return a.result, a.done }
func (a MyGame) Version() int                                     { return 1 }
func (a MyGame) IsBooleanResult() bool                            { return true }
func (a MyGame) IsStringFlag() bool                               { return false }
func (a *MyGame) SetArguments(args []string) error                { return nil }
func (a *MyGame) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *MyGame) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = "false"
	for _, p := range replay.Header.Players {
		if _, ok := ctx.Me[p.Name]; ok {
			a.result = "true"
			break
		}
	}
	a.done = true
	return nil, a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MapName struct {
	done   bool
	result string
}

func (a MapName) Name() string                                     { return "map-name" }
func (a MapName) Description() string                              { return "Analyzes the map's name." }
func (a MapName) DependsOn() map[string]struct{}                   { return map[string]struct{}{} }
func (a MapName) IsDone() (string, bool)                           { return a.result, a.done }
func (a MapName) Version() int                                     { return 1 }
func (a MapName) IsBooleanResult() bool                            { return false }
func (a MapName) IsStringFlag() bool                               { return false }
func (a *MapName) SetArguments(args []string) error                { return nil }
func (a *MapName) ProcessCommand(command repcmd.Cmd) (error, bool) { return nil, true }
func (a *MapName) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) (error, bool) {
	a.result = replay.Header.Map
	a.done = true
	return nil, a.done
}
