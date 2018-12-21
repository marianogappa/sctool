package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
	"github.com/icza/screp/repparser"
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
	Value() bool
}

type boolResult struct {
	result bool
}

func (r boolResult) Value() bool { return r.result }

type isThereAZerg struct{ done, result bool }

func (a isThereAZerg) Name() string                            { return "is-there-a-zerg" }
func (a isThereAZerg) Description() string                     { return "Analyzes if there is a zerg player in the replay." }
func (a isThereAZerg) DependsOn() map[string]struct{}          { return map[string]struct{}{} }
func (a isThereAZerg) IsDone() (Result, bool)                  { return boolResult{a.result}, a.done }
func (a isThereAZerg) Version() int                            { return 1 }
func (a *isThereAZerg) SetArguments(args []string)             {}
func (a *isThereAZerg) ProcessCommand(command repcmd.Cmd) bool { return true }
func (a *isThereAZerg) StartReadingReplay(replay *rep.Replay) bool {
	for _, p := range replay.Header.OrigPlayers {
		if p.Race.ShortName == "zerg" {
			a.result = true
		}
	}
	a.done = true
	return a.done
}

func main() {
	var (
		_analyzers = map[string]Analyzer{(&isThereAZerg{}).Name(): &isThereAZerg{}}
		flags      = map[string]*bool{}
		fReplay    = flag.String("replay", "", "path to replay file")
		fReplays   = flag.String("replays", "", "comma-separated paths to replay files")
		fReplayDir = flag.String("replay-dir", "", "path to folder with replays (recursive)")
	)
	for name, a := range _analyzers {
		flags[name] = flag.Bool(name, false, a.Description())
	}
	flag.Parse()
	var analyzers = map[string]Analyzer{}
	for name, f := range flags {
		if *f {
			analyzers[name] = _analyzers[name]
		}
	}

	// Parse replay filename flags
	var replays = map[string]struct{}{}
	*fReplay = strings.TrimSpace(*fReplay)
	if len(*fReplay) >= 5 && (*fReplay)[len(*fReplay)-4:] == ".rep" {
		replays[*fReplay] = struct{}{}
	}
	if *fReplays != "" {
		for _, r := range strings.Split(*fReplays, ",") {
			r = strings.TrimSpace(r)
			if len(r) >= 5 && r[len(r)-4:] == ".rep" {
				replays[r] = struct{}{}
			}
		}
	}
	if *fReplayDir != "" {
		e := filepath.Walk(*fReplayDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && len(info.Name()) >= 5 && info.Name()[len(info.Name())-4:] == ".rep" {
				r := path
				replays[r] = struct{}{}
			}
			return nil
		})
		if e != nil {
			log.Fatal(e)
		}
	}
	for replay := range replays {
		analyzerInstances := make(map[string]Analyzer, len(analyzers))
		for n, a := range analyzers {
			analyzerInstances[n] = a
		}

		r, err := repparser.ParseFile(replay)
		if err != nil {
			fmt.Printf("Failed to parse replay: %v\n", err)
			continue
		}
		tryCompute(r)

		var results = map[string]Result{}
		for name, a := range analyzerInstances {
			if a.StartReadingReplay(r) {
				results[name], _ = a.IsDone()
				delete(analyzerInstances, name)
			}
		}
		for _, c := range r.Commands.Cmds {
			for name, a := range analyzerInstances {
				if a.ProcessCommand(c) {
					results[name], _ = a.IsDone()
					delete(analyzerInstances, name)
				}
			}

		}

		for name, result := range results {
			fmt.Printf("Result for replay %v, analyzer %v: %v\n", replay, name, result.Value())
		}
	}
}

func tryCompute(r *rep.Replay) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered panic: %v", r)
		}
	}()
	r.Compute()
}
