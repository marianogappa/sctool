package analyzer

import (
	"fmt"
	"path/filepath"
	"sort"
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
	// Name is used for DependsOn() and as argument to CLI, so must be
	// hyphenated and without spaces or special characters.
	Name() string

	// Description is a human readable description for what the analyzer is useful for. Used
	// in command line usage help.
	Description() string

	// SetArguments for running: should be called before StartReadingReplay().
	// It may error, signaling that this Analyzer should not be used, and an error
	// should be shown to the client, but execution of the rest may continue.
	SetArguments(args []string) error

	// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
	DependsOn() map[string]struct{}

	// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
	// Should not read anything from Commands; use ProcessCommand for that.
	// Returns true if the analyzer is finished calculating the result (i.e. no need
	// to process commands)
	// It may error, signaling that this Analyzer should no longer be used, and an error
	// should be shown to the client, but execution of the rest may continue.
	StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error)

	// ProcessCommand should be called for every command during a Replay analizing cycle.
	// StartReadingReplay should be called before processing any command, to refresh
	// any state and to decide if processing commands are necessary to determine result.
	// Returns true if the analyzer is finished calculating the result (i.e. no need
	// to process further commands).
	// It may error, signaling that this Analyzer should no longer be used, and an error
	// should be shown to the client, but execution of the rest may continue.
	ProcessCommand(command repcmd.Cmd) (bool, error)

	// IsDone Returns true if the analyzer is finished calculating the result, and
	// returns it. Shouldn't be called before calling StartReadingReplay.
	IsDone() (string, bool)

	// Version is useful for managing updates to an Analyzer: whenever an update is made to an
	// analyzer, the Version should be numerically higher. Then, if there's a cached
	// Result of an Analyzer on a Replay, the result should be recomputed.
	Version() int

	// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
	IsStringFlag() bool

	// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
	// flags.
	IsBooleanResult() bool

	// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
	Clone() Analyzer
}

// Analyzers are all implemented replay analyzers. Must not be modified!
var Analyzers = map[string]Analyzer{
	(&MyAPM{}).Name():                        &MyAPM{},
	(&MyRace{}).Name():                       &MyRace{},
	(&DateTime{}).Name():                     &DateTime{},
	(&DurationMinutes{}).Name():              &DurationMinutes{},
	(&DurationMinutesIsGreaterThan{}).Name(): &DurationMinutesIsGreaterThan{},
	(&DurationMinutesIsLowerThan{}).Name():   &DurationMinutesIsLowerThan{},
	(&MyName{}).Name():                       &MyName{},
	(&IsThereARace{}).Name():                 &IsThereARace{},
	(&MyRaceIs{}).Name():                     &MyRaceIs{},
	(&ReplayName{}).Name():                   &ReplayName{},
	(&ReplayPath{}).Name():                   &ReplayPath{},
	(&MyWin{}).Name():                        &MyWin{},
	(&MyGame{}).Name():                       &MyGame{},
	(&MapName{}).Name():                      &MapName{},
	(&MyFirstSpecificUnitSeconds{}).Name():   &MyFirstSpecificUnitSeconds{},
	(&Matchup{}).Name():                      &Matchup{},
	(&MyMatchup{}).Name():                    &MyMatchup{},
	(&MatchupIs{}).Name():                    &MatchupIs{},
	(&MyMatchupIs{}).Name():                  &MyMatchupIs{},
	(&Is1v1{}).Name():                        &Is1v1{},
	(&Is2v2{}).Name():                        &Is2v2{},
	// TODO MyBOIs9Pool
	// TODO MyBOIs12Pool
	// TODO MyBOIsOverpool
	// TODO MyBOIs12Hatch
	// TODO MyBOIs3HatchBeforePool
	// TODO MyBOIs2HatchBeforePool
	// TODO MyBOIs1-1-1
	// TODO MyBOIs2HatchSpire
	// TODO IsTopVsBottomOrOneOnOneOrMelee
}

// Context is all context necessary for analyzers to properly analyze a replay
type Context struct {
	Me map[string]struct{}
}

// NewContext creates an Analyzer Context. Context should be everything unrelated to a replay that an Analyzer should
// know in order to analyze a replay e.g. who is the -me player
func NewContext(me map[string]struct{}) Context {
	return Context{me}
}

// Executor is the main struct the client should interact with: it receives a list of replays and analyzer
// requests, and executes the analyzers on the replays.
type Executor struct {
	replayPaths                []string
	analyzerWrappers           []analyzerWrapper
	ctx                        Context
	output                     Output
	copyPath                   string
	shouldCopyToOutputLocation bool
}

// NewExecutor should be the entrypoint of this library to the client. It creates an Executor.
// It may return several errors: replay paths may not exist, analyzer requests may be for unknown analyzers,
// copy path may not exist, etc.
func NewExecutor(replayPaths []string, analyzerRequests [][]string, ctx Context, output Output, copyPath string) (*Executor, []error) {
	var (
		errs, rpErrs, aeErrs []error
		ae                   = &Executor{}
	)
	ae.replayPaths, rpErrs = ae.filterReplayPaths(replayPaths)
	ae.analyzerWrappers, aeErrs = ae.createSortedAnalyzerWrappers(analyzerRequests)
	ae.ctx = ctx
	ae.output = output
	if ae.output == nil {
		ae.output = NewNoOutput()
	}
	if copyPath != "" {
		ae.shouldCopyToOutputLocation = true
		if ok, err := isFileExist(copyPath); !ok || err != nil {
			ae.shouldCopyToOutputLocation = false
			if copyPath != "" && !ok {
				errs = append(errs, fmt.Errorf("output directory doesn't exist: %v", copyPath))
			}
			if copyPath != "" && err != nil {
				errs = append(errs, fmt.Errorf("error locating output directory (%v): %v", copyPath, err))
			}
		}
		ae.copyPath = copyPath
	}
	errs = append(errs, rpErrs...)
	errs = append(errs, aeErrs...)
	return ae, errs
}

// Execute executes the given Analyzers on the given replays. Use this method if you want JSON/CSV output onto a file
// or Stdout. If you're using the library and want to work with the results programatically, use ExecuteWithResults
// instead.
// Execute may fail for a number of reasons, but the failures may be isolated to a replay or even a single analyzer,
// so it doesn't necessarily mean that the complete result is unusable.
func (e *Executor) Execute() []error {
	_, errs := e.execute(false)
	return errs
}

// ExecuteWithResults executes the given Analyzers on the given replays and returns all results. Use this method if you
// want to work with the results programatically; otherwise use Execute instead.
// ExecuteWithResults may fail for a number of reasons, but the failures may be isolated to a replay or even a single
// analyzer, so it doesn't necessarily mean that the complete result is unusable.
// Usually, when using this method, the Executor's output should be NoOutput, since the results are
// available in the return value of this method.
func (e *Executor) ExecuteWithResults() ([][]string, []error) {
	results, errs := e.execute(false)
	return results, errs
}

func (e *Executor) execute(saveResults bool) ([][]string, []error) {
	var (
		results [][]string
		errs    []error
	)
	if err := e.output.Pre(e.analyzerWrappers); err != nil { // CSV/JSON setup
		errs = append(errs, err)
	}
	for _, replayPath := range e.replayPaths {
		r, err := e.parseReplayFile(replayPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		replayResult, erErrs := e.executeReplay(r, replayPath, e.cloneAnalyzerWrappers())
		errs = append(errs, erErrs...)
		if len(replayResult) == 0 {
			continue
		}
		if err := e.output.ReplayResults(replayResult); err != nil { // CSV/JSON write line
			errs = append(errs, err)
		}
		if saveResults { // If used as library
			results = append(results, replayResult)
		}
		if e.shouldCopyToOutputLocation {
			if err := copyFile(replayPath, fmt.Sprintf("%v/%v", e.copyPath, filepath.Base(replayPath))); err != nil {
				errs = append(errs, fmt.Errorf("error copying replay with path %v to %v: %v", replayPath, e.copyPath, err))
			}
		}
	}
	if err := e.output.Post(); err != nil { // CSV/JSON teardown
		errs = append(errs, err)
	}
	return results, errs
}

func (e Executor) executeReplay(r *rep.Replay, replayPath string, analyzerWrappers []analyzerWrapper) ([]string, []error) {
	var (
		results      = make([]string, len(analyzerWrappers))
		errs         = []error{}
		removedCount = 0
	)

	// Analyze everything except Commands; try to finish early
	for i, aw := range analyzerWrappers {
		done, err := aw.analyzer.StartReadingReplay(r, e.ctx, replayPath)
		if err != nil {
			errs = append(errs,
				fmt.Errorf("error beginning to read replay %v with Analyzer %v: %v", replayPath,
					aw.analyzer.Name(), err))
		}
		if done {
			results[i], _ = aw.analyzer.IsDone()
		}
		if done && ((aw.isFilterNot && results[i] == "true") || (aw.isFilter && results[i] != "true")) {
			return []string{}, errs // Optimization: move to next replay if excluded by any filters already
		}
		aw.removed = done || err != nil // if analyzer is done or had error, signal that commands needn't be processed
		if aw.removed {
			removedCount++
		}
	}

	// Analyze Commands. N.B. This is the expensive loop in the algorithm; optimize here!
	for _, c := range r.Commands.Cmds {
		if len(analyzerWrappers) == removedCount {
			break // Optimization: don't loop over commands if there's nothing to do!
		}
		for i, aw := range analyzerWrappers {
			if aw.removed {
				continue
			}
			done, err := aw.analyzer.ProcessCommand(c)
			if err != nil {
				errs = append(errs,
					fmt.Errorf("error reading command on replay %v with Analyzer %v: %v",
						replayPath, aw.analyzer.Name(), err))
			}
			if done {
				results[i], _ = aw.analyzer.IsDone()
			}
			if done && ((aw.isFilterNot && results[i] == "true") || (aw.isFilter && results[i] != "true")) {
				return []string{}, errs // Optimization: move to next replay if excluded by any filters already
			}
			aw.removed = done || err != nil // if analyzer is done or had error, signal that more commands needn't be processed
			if aw.removed {
				removedCount++
			}
		}
	}
	return results, errs
}

func (e Executor) parseReplayFile(replayPath string) (*rep.Replay, error) {
	r, err := repparser.ParseFile(replayPath)
	if err != nil {
		return nil, fmt.Errorf("screp failed to parse replay %v: %v", replayPath, err)
	}
	if err := tryCompute(r); err != nil {
		return nil, err
	}
	return r, nil
}

type analyzerWrapper struct {
	analyzer    Analyzer
	isFilter    bool
	isFilterNot bool
	displayName string
	name        string
	pos         int
	removed     bool // used to signal that commands needn't be processed
}

func (w analyzerWrapper) less(w2 analyzerWrapper) bool {
	return ((w.isFilter || w.isFilterNot) && !w2.isFilter && !w2.isFilterNot) ||
		(((w.isFilter || w.isFilterNot) == (w2.isFilter || w2.isFilterNot)) && w.analyzer.Name() < w2.analyzer.Name())
}

func (w analyzerWrapper) clone() analyzerWrapper {
	return analyzerWrapper{w.analyzer.Clone(), w.isFilter, w.isFilterNot, w.displayName, w.name, w.pos, false}
}

func (e Executor) createSortedAnalyzerWrappers(analyzerRequests [][]string) ([]analyzerWrapper, []error) {
	var (
		analyzerWrappers = []analyzerWrapper{}
		errs             = []error{}
	)
	for i, analyzerRequest := range analyzerRequests {
		if len(analyzerRequest) == 0 {
			continue
		}
		var isFilter, isFilterNot bool
		if strings.HasPrefix(analyzerRequest[0], "filter--") {
			analyzerRequest[0] = analyzerRequest[0][len("filter--"):]
			isFilter = true
		} else if strings.HasPrefix(analyzerRequest[0], "filter-not--") {
			analyzerRequest[0] = analyzerRequest[0][len("filter-not--"):]
			isFilterNot = true
		}
		_analyzer, ok := Analyzers[analyzerRequest[0]]
		if !ok {
			errs = append(errs, fmt.Errorf("analyzer for name %v not found; ignoring", analyzerRequest[0]))
			continue
		}
		an := _analyzer.Clone()
		if err := an.SetArguments(analyzerRequest[1:]); err != nil {
			errs = append(errs, fmt.Errorf("error setting arguments for analyzer %v: %v; ignoring", an.Name(), err))
			continue
		}
		displayName := an.Name()
		if len(analyzerRequest[1:]) > 0 {
			displayName = fmt.Sprintf("%v(%v)", an.Name(), strings.Join(analyzerRequest[1:], ","))
		}
		analyzerWrappers = append(analyzerWrappers, analyzerWrapper{
			analyzer:    an,
			name:        fmt.Sprintf("%v_%v", an.Name(), i),
			displayName: displayName,
			isFilter:    isFilter,
			isFilterNot: isFilterNot,
		})
	}
	sort.Slice(analyzerWrappers, func(i, j int) bool {
		return analyzerWrappers[i].less(analyzerWrappers[j])
	})
	for i := range analyzerWrappers {
		analyzerWrappers[i].pos = i
	}
	return analyzerWrappers, errs
}

func (e Executor) filterReplayPaths(replayPaths []string) (paths []string, errs []error) {
	var m = make(map[string]struct{}, len(replayPaths))
	for _, r := range replayPaths { // Trim and unique
		m[strings.TrimSpace(r)] = struct{}{}
	}
	for replayPath := range m {
		if replayPath == "" || len(replayPath) < 5 || replayPath[len(replayPath)-4:] != ".rep" {
			continue
		}
		ok, err := isFileExist(replayPath)
		if err != nil {
			errs = append(errs, err)
		} else if !ok {
			errs = append(errs, fmt.Errorf("replay path not found: %v", replayPath))
		} else {
			paths = append(paths, replayPath)
		}
	}
	sort.Strings(paths)
	return
}

func (e Executor) cloneAnalyzerWrappers() (as []analyzerWrapper) {
	for _, analyzerWrapper := range e.analyzerWrappers {
		as = append(as, analyzerWrapper.clone())
	}
	return
}

func tryCompute(r *rep.Replay) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic in screp: %v", r)
		}
	}()
	r.Compute()
	return
}
