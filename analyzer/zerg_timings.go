package analyzer

import (
	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// -------------------------------------------------------------------------------------------------------------------
type MySpawningPoolSeconds struct {
	done     bool
	result   string
	playerID byte
	unitID   uint16
}

func (a MySpawningPoolSeconds) Name() string { return "my-spawning-pool-seconds" }
func (a MySpawningPoolSeconds) Description() string {
	return "Analyzes the time the first Spawning Pool was built, in seconds."
}
func (a MySpawningPoolSeconds) DependsOn() map[string]struct{} { return map[string]struct{}{} }
func (a MySpawningPoolSeconds) IsDone() (Result, bool)         { return stringResult{a.result}, a.done }
func (a MySpawningPoolSeconds) Version() int                   { return 1 }
func (a *MySpawningPoolSeconds) SetArguments(args []string)    {}
func (a *MySpawningPoolSeconds) StartReadingReplay(replay *rep.Replay, ctx AnalyzerContext, replayPath string) bool {
	a.result = "-1"
	a.unitID = 0x8E
	a.playerID = findPlayerID(replay, ctx.Me)
	a.done = a.playerID == 127 // If we don't find it, no need to see commands
	return a.done
}
func (a *MySpawningPoolSeconds) ProcessCommand(command repcmd.Cmd) bool {
	a.result, a.done = maybePlayersUnitSeconds(command, a.playerID, a.unitID)
	return a.done
}
