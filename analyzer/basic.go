package analyzer

import (
	"fmt"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// -------------------------------------------------------------------------------------------------------------------
type IsThereAZerg struct {
	done   bool
	result string
}

func (a IsThereAZerg) Name() string                            { return "is-there-a-zerg" }
func (a IsThereAZerg) Description() string                     { return "Analyzes if there is a zerg player in the replay." }
func (a IsThereAZerg) DependsOn() map[string]struct{}          { return map[string]struct{}{} }
func (a IsThereAZerg) IsDone() (Result, bool)                  { return stringResult{a.result}, a.done }
func (a IsThereAZerg) Version() int                            { return 1 }
func (a *IsThereAZerg) SetArguments(args []string)             {}
func (a *IsThereAZerg) ProcessCommand(command repcmd.Cmd) bool { return true }
func (a *IsThereAZerg) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext) bool {
	a.result = "false"
	for _, p := range replay.Header.OrigPlayers {
		if p.Race.ShortName == "zerg" {
			a.result = "true"
		}
	}
	a.done = true
	return a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MyAPM struct {
	done   bool
	result string
}

func (a MyAPM) Name() string                            { return "my-apm" }
func (a MyAPM) Description() string                     { return "Analyzes the APM of the -me player." }
func (a MyAPM) DependsOn() map[string]struct{}          { return map[string]struct{}{} }
func (a MyAPM) IsDone() (Result, bool)                  { return stringResult{a.result}, a.done }
func (a MyAPM) Version() int                            { return 1 }
func (a *MyAPM) SetArguments(args []string)             {}
func (a *MyAPM) ProcessCommand(command repcmd.Cmd) bool { return true }
func (a *MyAPM) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext) bool {
	if replay.Computed == nil {
		a.result = "-1"
		a.done = true
		return true
	}
	var playerID byte
	for _, p := range replay.Header.Players {
		if _, ok := ctx.Me[p.Name]; ok {
			playerID = p.ID
		}
	}
	for _, pDesc := range replay.Computed.PlayerDescs {
		if pDesc.PlayerID == playerID {
			a.result = fmt.Sprintf("%v", pDesc.APM)
			a.done = true
			break
		}
	}
	return a.done
}

// -------------------------------------------------------------------------------------------------------------------
type MyRace struct {
	done   bool
	result string
}

func (a MyRace) Name() string                            { return "my-race" }
func (a MyRace) Description() string                     { return "Analyzes the race of the -me player." }
func (a MyRace) DependsOn() map[string]struct{}          { return map[string]struct{}{} }
func (a MyRace) IsDone() (Result, bool)                  { return stringResult{a.result}, a.done }
func (a MyRace) Version() int                            { return 1 }
func (a *MyRace) SetArguments(args []string)             {}
func (a *MyRace) ProcessCommand(command repcmd.Cmd) bool { return true }
func (a *MyRace) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext) bool {
	if replay.Computed == nil {
		a.result = ""
		a.done = true
		return true
	}
	for _, p := range replay.Header.Players {
		if _, ok := ctx.Me[p.Name]; ok {
			a.result = p.Race.Name
			a.done = true
			break
		}
	}
	return a.done
}

// -------------------------------------------------------------------------------------------------------------------
type DateTime struct {
	done   bool
	result string
}

func (a DateTime) Name() string                            { return "date-time" }
func (a DateTime) Description() string                     { return "Analyzes the datetime of the replay." }
func (a DateTime) DependsOn() map[string]struct{}          { return map[string]struct{}{} }
func (a DateTime) IsDone() (Result, bool)                  { return stringResult{a.result}, a.done }
func (a DateTime) Version() int                            { return 1 }
func (a *DateTime) SetArguments(args []string)             {}
func (a *DateTime) ProcessCommand(command repcmd.Cmd) bool { return true }
func (a *DateTime) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext) bool {
	a.result = fmt.Sprintf("%v", replay.Header.StartTime)
	a.done = true
	return a.done
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
func (a DurationMinutes) DependsOn() map[string]struct{}          { return map[string]struct{}{} }
func (a DurationMinutes) IsDone() (Result, bool)                  { return stringResult{a.result}, a.done }
func (a DurationMinutes) Version() int                            { return 1 }
func (a *DurationMinutes) SetArguments(args []string)             {}
func (a *DurationMinutes) ProcessCommand(command repcmd.Cmd) bool { return true }
func (a *DurationMinutes) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext) bool {
	a.result = fmt.Sprintf("%v", int(replay.Header.Duration().Minutes()))
	a.done = true
	return a.done
}
