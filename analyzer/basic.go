package analyzer

import (
	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

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
func (a *IsThereAZerg) StartReadingReplay(replay *rep.Replay) bool {
	a.result = "false"
	for _, p := range replay.Header.OrigPlayers {
		if p.Race.ShortName == "zerg" {
			a.result = "true"
		}
	}
	a.done = true
	return a.done
}
