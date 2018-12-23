package analyzer

import (
	"fmt"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// -------------------------------------------------------------------------------------------------------------------
type MyFirstSpecificUnitSeconds struct {
	done     bool
	result   string
	playerID byte
	unitID   uint16
}

func (a MyFirstSpecificUnitSeconds) Name() string { return "my-first-specific-unit-seconds" }
func (a MyFirstSpecificUnitSeconds) Description() string {
	return "Analyzes the time the first specified unit/building/evolution was built, in seconds."
}
func (a MyFirstSpecificUnitSeconds) DependsOn() map[string]struct{} { return map[string]struct{}{} }
func (a MyFirstSpecificUnitSeconds) IsDone() (string, bool)         { return a.result, a.done }
func (a MyFirstSpecificUnitSeconds) Version() int                   { return 1 }
func (a MyFirstSpecificUnitSeconds) Clone() Analyzer {
	return &MyFirstSpecificUnitSeconds{a.done, a.result, a.playerID, a.unitID}
}
func (a MyFirstSpecificUnitSeconds) IsBooleanResult() bool { return false }
func (a MyFirstSpecificUnitSeconds) IsStringFlag() bool    { return true }
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
func (a *MyFirstSpecificUnitSeconds) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (error, bool) {
	a.result = "-1"
	a.playerID = findPlayerID(replay, ctx.Me)
	a.done = a.playerID == 127 // If we don't find it, no need to see commands
	return nil, a.done
}
func (a *MyFirstSpecificUnitSeconds) ProcessCommand(command repcmd.Cmd) (error, bool) {
	a.result, a.done = maybePlayersUnitSeconds(command, a.playerID, a.unitID)
	return nil, a.done
}
