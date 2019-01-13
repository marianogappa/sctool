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

// Analyzers are all implemented replay analyzers. Should be cloned before being used.
var Analyzers = map[string]Analyzer{
	"is-there-a-race": newAnalyzerImpl(
		"is-there-a-race",
		"Analyzes if there is a specific race in the replay.",
		1, // version
		map[string]struct{}{}, // dependsOn
		true,  // isStringFlag
		true,  // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorRace{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				for _, p := range replay.Header.OrigPlayers {
					if p.Race.Name == args[0] {
						return "true", true, nil, nil
					}
				}
				return "false", true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-apm": newAnalyzerImpl(
		"my-apm",
		"Analyzes the APM of the -me player.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		true,  // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "-1",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				if replay.Computed == nil {
					return "-1", true, nil, nil
				}
				playerID := findPlayerID(replay, ctx.Me)
				if playerID == 127 {
					return "-1", true, nil, fmt.Errorf("-me player not present in this replay")
				}
				for _, pDesc := range replay.Computed.PlayerDescs {
					if pDesc.PlayerID == playerID {
						return fmt.Sprintf("%v", pDesc.APM), true, nil, nil
					}
				}
				return "-1", true, nil, fmt.Errorf("unexpected error")
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-race": newAnalyzerImpl(
		"my-race",
		"Analyzes the race of the -me player.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				playerID := findPlayerID(replay, ctx.Me)
				if playerID == 127 {
					return "", true, nil, fmt.Errorf("-me player not present in this replay")
				}
				return replay.Header.PIDPlayers[playerID].Race.Name, true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-race-is": newAnalyzerImpl(
		"my-race-is",
		"Analyzes if the race of the -me player is the one specified.",
		1, // version
		map[string]struct{}{}, // dependsOn
		true,  // isStringFlag
		true,  // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorRace{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				playerID := findPlayerID(replay, ctx.Me)
				if playerID == 127 {
					return "", true, nil, fmt.Errorf("-me player not present in this replay")
				}
				result := "false"
				if replay.Header.PIDPlayers[playerID].Race.Name == args[0] {
					result = "true"
				}
				return result, true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"date": newAnalyzerImpl(
		"date",
		"Analyzes the date of the replay. Uses yyyy-mm-dd pattern because it's lexicographically sorted.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				return replay.Header.StartTime.Format("2006-01-02"), true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-name": newAnalyzerImpl(
		"my-name",
		"Analyzes the name of the -me player.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				playerID := findPlayerID(replay, ctx.Me)
				if playerID == 127 {
					return "", true, nil, fmt.Errorf("-me player not present in this replay")
				}
				return replay.Header.PIDPlayers[playerID].Name, true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"replay-name": newAnalyzerImpl(
		"replay-name",
		"Analyzes the replay's name.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				result := path.Base(replayPath)
				return result[:len(result)-4], true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"replay-path": newAnalyzerImpl(
		"replay-path",
		"Analyzes the replay's path.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				return replayPath, true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-win": newAnalyzerImpl(
		"my-win",
		"Analyzes if the -me player won the game.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				if replay.Computed == nil || replay.Computed.WinnerTeam == 0 {
					return "unknown", true, nil, nil
				}
				playerID := findPlayerID(replay, ctx.Me)
				if playerID == 127 {
					return "", true, nil, fmt.Errorf("-me player not present in this replay")
				}
				if replay.Header.PIDPlayers[playerID].Team == replay.Computed.WinnerTeam {
					return "true", true, nil, nil
				}
				return "false", true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-game": newAnalyzerImpl(
		"my-game",
		"Analyzes if the -me player played the game.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		true,  // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				if playerID := findPlayerID(replay, ctx.Me); playerID == 127 {
					return "false", true, nil, nil
				}
				return "true", true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"map-name": newAnalyzerImpl(
		"map-name",
		"Analyzes the map's name. Note that it doesn't do anything clever, so many versions of a map can have slightly different names, or two maps with the same name might be actually different.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				return replay.Header.Map, true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"is-1v1": newAnalyzerImpl(
		"is-1v1",
		"Analyzes if the replay is of an 1v1 match.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				if len(replay.Header.Players) == 2 && replay.Header.Players[0].Team != replay.Header.Players[1].Team {
					return "true", true, nil, nil
				}
				return "false", true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"is-2v2": newAnalyzerImpl(
		"is-2v2",
		"Analyzes if the replay is of a 2v2 match.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				if len(replay.Header.Players) == 4 && replay.Header.Players[0].Team == replay.Header.Players[1].Team &&
					replay.Header.Players[1].Team != replay.Header.Players[2].Team &&
					replay.Header.Players[2].Team == replay.Header.Players[3].Team {
					return "true", true, nil, nil
				}
				return "false", true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"duration-minutes": newAnalyzerImpl(
		"duration-minutes",
		"Analyzes the duration of the replay in minutes.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				return fmt.Sprintf("%v", int(replay.Header.Duration().Minutes())), true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"duration-minutes-is-greater-than": newAnalyzerImpl(
		"duration-minutes-is-greater-than",
		"Analyzes if the duration of the replay in minutes is greater than specified.",
		1, // version
		map[string]struct{}{}, // dependsOn
		true,  // isStringFlag
		true,  // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorMinutes{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				actualMinutes := int(replay.Header.Duration().Minutes())
				expectedMinutes, _ := strconv.Atoi(args[0]) // N.B. Validator already checked it's ok
				return fmt.Sprintf("%v", actualMinutes > expectedMinutes), true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"duration-minutes-is-lower-than": newAnalyzerImpl(
		"duration-minutes-is-lower-than",
		"Analyzes if the duration of the replay in minutes is lower than specified.",
		1, // version
		map[string]struct{}{}, // dependsOn
		true,  // isStringFlag
		true,  // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorMinutes{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				actualMinutes := int(replay.Header.Duration().Minutes())
				expectedMinutes, _ := strconv.Atoi(args[0]) // N.B. Validator already checked it's ok
				return fmt.Sprintf("%v", actualMinutes < expectedMinutes), true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"matchup": newAnalyzerImpl(
		"matchup",
		"Analyzes the replay's matchup. On an 1v1, it will sort the races lexicographically, so it will return TvZ rather than ZvT. Other than 1v1, it will simply return whatever screp returns.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				if len(replay.Header.Players) == 2 {
					r0 := strings.ToUpper(string(replay.Header.Players[0].Race.Letter))
					r1 := strings.ToUpper(string(replay.Header.Players[1].Race.Letter))
					if r0 > r1 {
						return r1 + "v" + r0, true, nil, nil
					}
					return r0 + "v" + r1, true, nil, nil
				}
				return replay.Header.Matchup(), true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-matchup": newAnalyzerImpl(
		"my-matchup",
		"Analyzes the replay's matchup from the point of view of the -me player. For example, if the -me player is Z and the opponent is T it will return ZvT rather than TvZ. At the moment, the behaviour other than 1v1 is unexpected: it returns the matchup as returned by screp.",
		1, // version
		map[string]struct{}{}, // dependsOn
		false, // isStringFlag
		false, // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorNoArguments{},
		&analyzerProcessorImpl{
			result: "",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				playerID := findPlayerID(replay, ctx.Me)
				if playerID == 127 {
					return "", true, nil, nil
				}
				if len(replay.Header.Players) == 2 {
					r0 := strings.ToUpper(string(replay.Header.Players[0].Race.Letter))
					r1 := strings.ToUpper(string(replay.Header.Players[1].Race.Letter))
					if playerID == 1 {
						return r1 + "v" + r0, true, nil, nil
					}
					return r0 + "v" + r1, true, nil, nil
				}
				return replay.Header.Matchup(), true, nil, nil // TODO put -me player on the left side
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"matchup-is": newAnalyzerImpl(
		"matchup-is",
		"Analyzes if the replay's MatchupIs is equal to the specified one (only works for 1v1 for now). The specified matchup can be in either order (i.e. ZvT == TvZ).",
		1, // version
		map[string]struct{}{}, // dependsOn
		true,  // isStringFlag
		true,  // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidator1v1Matchup{},
		&analyzerProcessorImpl{
			result: "false",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				if len(replay.Header.Players) != 2 {
					return "", true, nil, nil
				}
				actualRaces := []string{
					strings.ToUpper(string(replay.Header.Players[0].Race.Letter)),
					strings.ToUpper(string(replay.Header.Players[1].Race.Letter)),
				}
				sort.Strings(actualRaces)
				return fmt.Sprintf("%v", reflect.DeepEqual(args, actualRaces)), true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-matchup-is": newAnalyzerImpl(
		"my-matchup-is",
		"Analyzes if the replay's matchup is equal to the specified one, from the -me player perspective (only works for 1v1 for now). The specified matchup must contain the -me player's race first.",
		1, // version
		map[string]struct{}{}, // dependsOn
		true,  // isStringFlag
		true,  // isBooleanResult
		false, // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidator1v1Matchup{},
		&analyzerProcessorImpl{
			result: "false",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				playerID := findPlayerID(replay, ctx.Me)
				if playerID == 127 || len(replay.Header.Players) != 2 {
					return "", true, nil, nil
				}
				actualRaces := []string{
					strings.ToUpper(string(replay.Header.Players[0].Race.Letter)),
					strings.ToUpper(string(replay.Header.Players[1].Race.Letter)),
				}
				sort.Strings(actualRaces)
				return fmt.Sprintf("%v", reflect.DeepEqual(args, actualRaces)), true, nil, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				return result, true, nil
			},
		},
	),
	"my-first-specific-unit-seconds": newAnalyzerImpl(
		"my-first-specific-unit-seconds",
		"Analyzes the time the first specified unit/building/evolution was built, in seconds. Refer to the unit name list in utils.go#nameToUnitID. -1 if the unit never appears.",
		1, // version
		map[string]struct{}{}, // dependsOn
		true,  // isStringFlag
		false, // isBooleanResult
		true,  // requiresParsingCommands
		false, // requiresParsingMapData
		&argumentValidatorUnit{},
		&analyzerProcessorImpl{
			result: "-1",
			done:   false,
			startReadingReplay: func(replay *rep.Replay, ctx Context, replayPath string, args []string) (string, bool, interface{}, error) {
				unitID, _ := strconv.Atoi(args[0]) // N.B. Validator already checked it's ok
				playerID := findPlayerID(replay, ctx.Me)
				state := []int{unitID, int(playerID)}
				return "-1", playerID == 127, state, nil
			},
			processCommand: func(command repcmd.Cmd, args []string, result string, state interface{}) (string, bool, error) {
				_state := state.([]int)
				unitID, playerID := _state[0], _state[1]
				result, done := maybePlayersUnitSeconds(command, byte(playerID), uint16(unitID))
				return result, done, nil
			},
		},
	),
}
