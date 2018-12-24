package analyzer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

// Note to contributors (and to self from the past): if possible, ignore this struct and file altogether.
// The reason these structs and interfaces exist is simply to write Analyzers that comply with the large
// Analyzer interface without writing so much code or comments (and godocs).
type analyzerImpl struct {
	name                    string
	description             string
	args                    []string
	version                 int
	dependsOn               map[string]struct{}
	isStringFlag            bool
	isBooleanResult         bool
	requiresParsingCommands bool
	requiresParsingMapData  bool
	argumentValidator       argumentValidator
	analyzerProcessor       analyzerProcessor

	result string
	done   bool
}

// Name is used for DependsOn() and as argument to CLI, so must be
// hyphenated and without spaces or special characters.
func (a analyzerImpl) Name() string { return a.name }

// Description is a human readable description for what the analyzer is useful for. Used
// in command line usage help.
func (a analyzerImpl) Description() string { return a.description }

// DependsOn are the Analyzer Name's whose Results this Analyzer depends on: for building DAG.
func (a analyzerImpl) DependsOn() map[string]struct{} { return a.dependsOn }

// IsDone Returns true if the analyzer is finished calculating the result, and
// returns it. Shouldn't be called before calling StartReadingReplay.
func (a analyzerImpl) IsDone() (string, bool) { return a.result, a.done }

// Version is useful for managing updates to an Analyzer: whenever an update is made to an
// analyzer, the Version should be numerically higher. Then, if there's a cached
// Result of an Analyzer on a Replay, the result should be recomputed.
func (a analyzerImpl) Version() int { return a.version }

// RequiresParsingCommands is true if this Analyzer requires parsing commands from the replay
func (a analyzerImpl) RequiresParsingCommands() bool { return a.requiresParsingCommands }

// RequiresParsingMapData is true if this Analyzer requires parsing map data from the replay
func (a analyzerImpl) RequiresParsingMapData() bool { return a.requiresParsingMapData }

// IsBooleanResult Determines if the result type is "true"/"false". Used for providing -filter-- and -filter-not--
// flags.
func (a analyzerImpl) IsBooleanResult() bool { return a.isBooleanResult }

// IsStringFlag determines the type of the CLI flag. It can either be Bool (default) or String.
func (a analyzerImpl) IsStringFlag() bool { return a.isStringFlag }

// Clone is a convenience method just so there can be a map[string]analyzer.Analyzer in createSortedAnalyzerWrappers
func (a analyzerImpl) Clone() Analyzer {
	return &analyzerImpl{
		name:                    a.name,
		description:             a.description,
		args:                    a.args,
		version:                 a.version,
		dependsOn:               a.dependsOn,
		isStringFlag:            a.isStringFlag,
		isBooleanResult:         a.isBooleanResult,
		requiresParsingCommands: a.requiresParsingCommands,
		requiresParsingMapData:  a.requiresParsingMapData,
		argumentValidator:       a.argumentValidator,
		analyzerProcessor:       a.analyzerProcessor.Clone(),
		done:                    a.done,
		result:                  a.result,
	}
}

// SetArguments for running: should be called before StartReadingReplay().
// It may error, signaling that this Analyzer should not be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *analyzerImpl) SetArguments(args []string) error {
	var err error
	a.args, err = a.argumentValidator.ValidateAndSet(args)
	return err
}

// ProcessCommand should be called for every command during a Replay analizing cycle.
// StartReadingReplay should be called before processing any command, to refresh
// any state and to decide if processing commands are necessary to determine result.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process further commands).
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *analyzerImpl) ProcessCommand(command repcmd.Cmd) (bool, error) {
	var err error
	a.result, a.done, err = a.analyzerProcessor.ProcessCommand(command, a.args, a.result)
	return a.done, err
}

// StartReadingReplay is called at the beginning of a Replay analyzing cycle.
// Should not read anything from Commands; use ProcessCommand for that.
// Returns true if the analyzer is finished calculating the result (i.e. no need
// to process commands)
// It may error, signaling that this Analyzer should no longer be used, and an error
// should be shown to the client, but execution of the rest may continue.
func (a *analyzerImpl) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string) (bool, error) {
	var err error
	a.result, a.done, err = a.analyzerProcessor.StartReadingReplay(replay, ctx, replayPath, a.args)
	return a.done, err
}

func newAnalyzerImpl(
	name,
	description string,
	version int,
	dependsOn map[string]struct{},
	isStringFlag,
	isBooleanResult,
	requiresParsingCommands,
	requiresParsingMapData bool,
	argumentValidator argumentValidator,
	analyzerProcessor analyzerProcessor,
) Analyzer {
	return &analyzerImpl{
		name:                    name,
		description:             description,
		args:                    []string{},
		version:                 version,
		dependsOn:               dependsOn,
		isStringFlag:            isStringFlag,
		isBooleanResult:         isBooleanResult,
		requiresParsingCommands: requiresParsingCommands,
		requiresParsingMapData:  requiresParsingMapData,
		argumentValidator:       argumentValidator,
		analyzerProcessor:       analyzerProcessor,
	}
}

type argumentValidator interface {
	// SetArguments for running: should be called before StartReadingReplay().
	// It may error, signaling that this Analyzer should not be used, and an error
	// should be shown to the client, but execution of the rest may continue.
	ValidateAndSet(args []string) ([]string, error)
}

type argumentValidatorNoArguments struct{}

func (a *argumentValidatorNoArguments) ValidateAndSet(args []string) ([]string, error) {
	return []string{}, nil
}

type argumentValidatorRace struct{}

func (a *argumentValidatorRace) ValidateAndSet(args []string) ([]string, error) {
	if len(args) < 1 {
		return []string{}, fmt.Errorf("please provide a valid race name e.g. Zerg/Protoss/Terran") // TODO provide list
	}
	r := strings.ToLower(args[0])
	if _, ok := raceNameTranslations[r]; !ok {
		return []string{}, fmt.Errorf("invalid race name %v", args[0]) // TODO provide list
	}
	return []string{raceNameTranslations[r]}, nil
}

type argumentValidatorMinutes struct{}

func (a *argumentValidatorMinutes) ValidateAndSet(args []string) ([]string, error) {
	if len(args) < 1 {
		return []string{}, fmt.Errorf("please provide a valid number of minutes")
	}
	if _, err := strconv.Atoi(args[0]); err != nil {
		return []string{}, fmt.Errorf("invalid number of minutes: %v", args[0])
	}
	return args, nil
}

type argumentValidator1v1Matchup struct{}

func (a *argumentValidator1v1Matchup) ValidateAndSet(args []string) ([]string, error) {
	if len(args) < 1 || len(args[0]) != 3 {
		return []string{}, fmt.Errorf("please provide a valid 1v1 matchup e.g. TvZ")
	}
	args[0] = strings.ToUpper(args[0])
	var res []string
	res = append(res, string(args[0][0]), string(args[0][2]))
	sort.Strings(res)
	return res, nil
}

type argumentValidatorUnit struct{}

func (a *argumentValidatorUnit) ValidateAndSet(args []string) ([]string, error) {
	if len(args) < 1 {
		return []string{}, fmt.Errorf("please provide a valid unit/building/evolution name e.g. Zergling") // TODO provide list
	}
	if _, ok := nameToUnitID[args[0]]; !ok {
		return []string{}, fmt.Errorf("invalid unit/building/evolution name") // TODO provide list
	}
	return []string{fmt.Sprintf("%v", nameToUnitID[args[0]])}, nil
}

type analyzerProcessor interface {
	ProcessCommand(command repcmd.Cmd, args []string, result string) (string, bool, error)
	StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, error)
	Clone() analyzerProcessor
}

type analyzerProcessorImpl struct {
	result             string
	done               bool
	startReadingReplay func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, error)
	processCommand     func(command repcmd.Cmd, args []string, result string) (string, bool, error)
}

func (a *analyzerProcessorImpl) Clone() analyzerProcessor {
	return &analyzerProcessorImpl{"", false, a.startReadingReplay, a.processCommand}
}

func (a *analyzerProcessorImpl) StartReadingReplay(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, error) {
	var err error
	a.result, a.done, err = a.startReadingReplay(replay, ctx, replayPath, args)
	return a.result, a.done, err
}

func (a *analyzerProcessorImpl) ProcessCommand(command repcmd.Cmd, args []string, result string) (string, bool, error) {
	var err error
	a.result, a.done, err = a.processCommand(command, args, result)
	return a.result, a.done, err
}
