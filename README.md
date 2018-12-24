# sctool [![Build Status](https://img.shields.io/travis/marianogappa/sctool.svg)](https://travis-ci.org/marianogappa/sctool) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/marianogappa/sctool/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/marianogappa/sctool?style=flat-square)](https://goreportcard.com/report/github.com/marianogappa/sctool) [![GoDoc](https://godoc.org/github.com/marianogappa/sctool/analyzer?status.svg)](https://godoc.org/github.com/marianogappa/sctool/analyzer)


sctool is a Starcraft: Remastered replay analyzer library and CLI tool built with 4 main use cases:

1. enabling Data Science research
2. organising replay directories
3. learning by analyzing and tinkering with pro gamer replays' data
4. improving one's skill by discovering strengths and weaknesses on one's past games

## Status

sctool is NOT ready for mainstream use. Help beta-testing is appreciated.

## Installation

Please visit the releases page (TODO: not yet available) and download the right binary for your operating system.

## Installation for Go devs

```
$ go get -u github.com/marianogappa/sctool
$ sctool -help
```

## Summary

- sctool allows you to analyze a replay, list of replays or directory (recursive) of replays with a number of "analyzers". An example analyzer is: `--matchup-is 1v1`.

- All analyzers that return true/false results are also available as "filters" and "filter-nots". For example, `--filter--matchup-is 1v1` will only match replays with only 2 players on 2 separate teams (i.e. 1v1 games), and `--filter-not--matchup-is 1v1` will only match replays that contain 3 or more players.

- sctool further allows you to copy all replays that matched your filter criteria to a given folder. This should enable you to organise a replay folder by whatever supported criteria you want, e.g. UMS games, 1v1 games, TvZ games, games on a particular map, etc.

- sctool's output is CSV by default, making it ideal for streamlining into a Data Science research project, but it can also return JSON, which is handy to compose with [jq](https://stedolan.github.io/jq/) and then possibly into [chart](https://github.com/marianogappa/chart) for charting.

- Thanks to DateTime analyzers and different kinds of filtering and segmentation, sctool can track your progress: for example, you can see your APM improvement on 1v1 games on this season's maps for the matchup you're having difficulties with.

## Usage

```
$ sctool -help
Usage of sctool:
  -copy-to-if-matches-filters string
    	copy replay files matched by -filter-- and not matched by -filter--not-- filters to specified directory
  -date-time
    	Analyzes the datetime of the replay.
  -duration-minutes
    	Analyzes the duration of the replay in minutes.
  -duration-minutes-is-greater-than string
    	Analyzes if the duration of the replay in minutes is greater than specified.
  -duration-minutes-is-lower-than string
    	Analyzes if the duration of the replay in minutes is lower than specified.
  -filter--duration-minutes-is-greater-than string
    	Filter for: Analyzes if the duration of the replay in minutes is greater than specified.
  -filter--duration-minutes-is-lower-than string
    	Filter for: Analyzes if the duration of the replay in minutes is lower than specified.
  -filter--is-1v1
    	Filter for: Analyzes if the replay is of an 1v1 match.
  -filter--is-2v2
    	Filter for: Analyzes if the replay is of a 2v2 match.
  -filter--is-there-a-race string
    	Filter for: Analyzes if there is a specific race in the replay.
  -filter--matchup-is string
    	Filter for: Analyzes if the replay's MatchupIs is equal to the specified one (only works for 1v1 for now).
  -filter--my-game
    	Filter for: Analyzes if the -me player played the game.
  -filter--my-matchup-is string
    	Filter for: Analyzes if the replay's matchup is equal to the specified one, from the -me player perspective (only works for 1v1 for now).
  -filter--my-race-is string
    	Filter for: Analyzes if the race of the -me player is the one specified.
  -filter--my-win
    	Filter for: Analyzes if the -me player won the game.
  -filter-not--duration-minutes-is-greater-than string
    	Filter-Not for: Analyzes if the duration of the replay in minutes is greater than specified.
  -filter-not--duration-minutes-is-lower-than string
    	Filter-Not for: Analyzes if the duration of the replay in minutes is lower than specified.
  -filter-not--is-1v1
    	Filter-Not for: Analyzes if the replay is of an 1v1 match.
  -filter-not--is-2v2
    	Filter-Not for: Analyzes if the replay is of a 2v2 match.
  -filter-not--is-there-a-race string
    	Filter-Not for: Analyzes if there is a specific race in the replay.
  -filter-not--matchup-is string
    	Filter-Not for: Analyzes if the replay's MatchupIs is equal to the specified one (only works for 1v1 for now).
  -filter-not--my-game
    	Filter-Not for: Analyzes if the -me player played the game.
  -filter-not--my-matchup-is string
    	Filter-Not for: Analyzes if the replay's matchup is equal to the specified one, from the -me player perspective (only works for 1v1 for now).
  -filter-not--my-race-is string
    	Filter-Not for: Analyzes if the race of the -me player is the one specified.
  -filter-not--my-win
    	Filter-Not for: Analyzes if the -me player won the game.
  -help
    	Returns help usage and exits.
  -is-1v1
    	Analyzes if the replay is of an 1v1 match.
  -is-2v2
    	Analyzes if the replay is of a 2v2 match.
  -is-there-a-race string
    	Analyzes if there is a specific race in the replay.
  -json
    	outputs a JSON instead of the default CSV
  -map-name
    	Analyzes the map's name.
  -matchup
    	Analyzes the replay's matchup.
  -matchup-is string
    	Analyzes if the replay's MatchupIs is equal to the specified one (only works for 1v1 for now).
  -me string
    	comma-separated list of player names to identify as the main player
  -my-apm
    	Analyzes the APM of the -me player.
  -my-first-specific-unit-seconds string
    	Analyzes the time the first specified unit/building/evolution was built, in seconds.
  -my-game
    	Analyzes if the -me player played the game.
  -my-matchup
    	Analyzes the replay's matchup from the point of view of the -me player
  -my-matchup-is string
    	Analyzes if the replay's matchup is equal to the specified one, from the -me player perspective (only works for 1v1 for now).
  -my-name
    	Analyzes the name of the -me player.
  -my-race
    	Analyzes the race of the -me player.
  -my-race-is string
    	Analyzes if the race of the -me player is the one specified.
  -my-win
    	Analyzes if the -me player won the game.
  -quiet
    	don't print any errors (discouraged: note that you can silence with 2>/dev/null).
  -replay string
    	(>= 1 replays required) path to replay file
  -replay-dir string
    	(>= 1 replays required) path to folder with replays (recursive)
  -replay-name
    	Analyzes the replay's name.
  -replay-path
    	Analyzes the replay's path.
  -replays string
    	(>= 1 replays required) comma-separated paths to replay files
```

## Building on top of the analyzer library

If rather than creating a CSV/JSON you're looking to access the replay analysing results programatically, don't use
the sctool; instead, import the sctool/analyzer library in your go code. Refer to [main.go](main.go) as an example
use and to the [godoc](https://godoc.org/github.com/marianogappa/sctool/analyzer#Analyzer) as reference.

TODO sample code

## Contribute

I'm happy to review Issues and PRs!

## Acknowledgements

sctool is only possible thanks to the amazing effort by [András Belicza](https://github.com/icza) who built the [screp](https://github.com/icza/screp) library that sctool is built on top of.
