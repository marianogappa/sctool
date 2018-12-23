package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/repparser"
	"github.com/marianogappa/raszagal/analyzer"
)

func main() {
	var (
		_analyzers = map[string]analyzer.Analyzer{
			(&analyzer.MyAPM{}).Name():                        &analyzer.MyAPM{},
			(&analyzer.MyRace{}).Name():                       &analyzer.MyRace{},
			(&analyzer.DateTime{}).Name():                     &analyzer.DateTime{},
			(&analyzer.DurationMinutes{}).Name():              &analyzer.DurationMinutes{},
			(&analyzer.DurationMinutesIsGreaterThan{}).Name(): &analyzer.DurationMinutesIsGreaterThan{},
			(&analyzer.DurationMinutesIsLowerThan{}).Name():   &analyzer.DurationMinutesIsLowerThan{},
			(&analyzer.MyName{}).Name():                       &analyzer.MyName{},
			(&analyzer.IsThereARace{}).Name():                 &analyzer.IsThereARace{},
			(&analyzer.MyRaceIs{}).Name():                     &analyzer.MyRaceIs{},
			(&analyzer.ReplayName{}).Name():                   &analyzer.ReplayName{},
			(&analyzer.ReplayPath{}).Name():                   &analyzer.ReplayPath{},
			(&analyzer.MyWin{}).Name():                        &analyzer.MyWin{},
			(&analyzer.MyGame{}).Name():                       &analyzer.MyGame{},
			(&analyzer.MapName{}).Name():                      &analyzer.MapName{},
			(&analyzer.MyFirstSpecificUnitSeconds{}).Name():   &analyzer.MyFirstSpecificUnitSeconds{},
			(&analyzer.Matchup{}).Name():                      &analyzer.Matchup{},
			(&analyzer.MyMatchup{}).Name():                    &analyzer.MyMatchup{},
			(&analyzer.MatchupIs{}).Name():                    &analyzer.MatchupIs{},
			// TODO MatchupIs
			// TODO Is1v1
			// TODO Is2v2
			// TODO MyBOIs9Pool
			// TODO MyBOIs12Pool
			// TODO MyBOIsOverpool
			// TODO MyBOIs12Hatch
			// TODO MyBOIs3HatchBeforePool
			// TODO MyBOIs2HatchBeforePool
			// TODO MyBOIs1-1-1
			// TODO MyBOIs2HatchSpire
		}
		boolFlags               = map[string]*bool{}
		stringFlags             = map[string]*string{}
		fReplay                 = flag.String("replay", "", "(>= 1 replays required) path to replay file")
		fReplays                = flag.String("replays", "", "(>= 1 replays required) comma-separated paths to replay files")
		fReplayDir              = flag.String("replay-dir", "", "(>= 1 replays required) path to folder with replays (recursive)")
		fMe                     = flag.String("me", "", "comma-separated list of player names to identify as the main player")
		fJSON                   = flag.Bool("json", false, "outputs a JSON instead of the default CSV")
		fCopyToIfMatchesFilters = flag.String("copy-to-if-matches-filters", "",
			"copy replay files matched by -filter-- and not matched by -filter--not-- filters to specified directory")
	)
	for name, a := range _analyzers {
		if a.IsStringFlag() {
			stringFlags[name] = flag.String(name, "", a.Description())
		} else {
			boolFlags[name] = flag.Bool(name, false, a.Description())
		}
		if a.IsBooleanResult() {
			boolFlags["filter--"+name] = flag.Bool("filter--"+name, false, "Filter for: "+a.Description())
			boolFlags["filter-not--"+name] = flag.Bool("filter-not--"+name, false, "Filter-Not for: "+a.Description())
		}
	}
	flag.Parse()
	var (
		analyzers                  = map[string]analyzer.Analyzer{}
		filters                    = map[string]struct{}{}
		filterNots                 = map[string]struct{}{}
		fieldNames                 = []string{} // TODO add filename
		shouldCopyToOutputLocation = true
	)
	if ok, err := isDirExist(*fCopyToIfMatchesFilters); *fCopyToIfMatchesFilters == "" || !ok || err != nil {
		shouldCopyToOutputLocation = false
		if *fCopyToIfMatchesFilters != "" && !ok {
			log.Printf("Output directory doesn't exist: %v\n", *fCopyToIfMatchesFilters)
		}
		if *fCopyToIfMatchesFilters != "" && err != nil {
			log.Printf("Error locating output directory (%v): %v\n", *fCopyToIfMatchesFilters, err)
		}
	}
	for name, f := range boolFlags {
		if *f {
			addToFieldNames := true
			if strings.HasPrefix(name, "filter--") {
				name = name[len("filter--"):] // side-effect so that the analyzer runs
				filters[name] = struct{}{}
				addToFieldNames = false
			} else if strings.HasPrefix(name, "filter-not--") {
				name = name[len("filter-not--"):] // side-effect so that the analyzer runs
				filterNots[name] = struct{}{}
				addToFieldNames = false
			}
			analyzers[name] = _analyzers[name]
			if addToFieldNames {
				fieldNames = append(fieldNames, name)
			}
		}
	}
	for name, f := range stringFlags {
		if *f != "" {
			a := _analyzers[name]
			if err := a.SetArguments(unmarshalArguments(*f)); err != nil {
				log.Printf("Invalid arguments '%v' for Analyzer %v: %v", *f, name, err)
				continue
			}
			analyzers[name] = _analyzers[name]
			fieldNames = append(fieldNames, name)
		}
	}

	// TODO implement library for programmatic use
	// Prepares for CSV output
	sort.Strings(fieldNames)
	w := csv.NewWriter(os.Stdout)
	if !*fJSON {
		w.Write(fieldNames)
	}

	// Prepares for JSON output
	firstJSONRow := true
	if *fJSON {
		fmt.Println("[")
	}

	// Prepares AnalyzerContext
	ctx := analyzer.AnalyzerContext{Me: map[string]struct{}{}}
	if fMe != nil && len(*fMe) > 0 {
		for _, name := range strings.Split(*fMe, ",") {
			ctx.Me[strings.TrimSpace(name)] = struct{}{}
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

	// Main loop parsing replays
	// TODO evaluate if it's important that replay traversal is non-deterministic
replayLoop:
	for replay := range replays {
		analyzerInstances := make(map[string]analyzer.Analyzer, len(analyzers))
		analyzerNames := []string{}
		for n, a := range analyzers {
			analyzerInstances[n] = a
			analyzerNames = append(analyzerNames, n)
		}
		sort.Slice(analyzerNames, func(i, j int) bool { // Optimization: execute analyzers that are filters first
			_, iIsInFilters := filters[analyzerNames[i]]
			_, iIsInFilterNots := filterNots[analyzerNames[i]]
			_, jIsInFilters := filters[analyzerNames[j]]
			_, jIsInFilterNots := filterNots[analyzerNames[j]]
			iIsImportant := iIsInFilters || iIsInFilterNots
			jIsImportant := jIsInFilters || jIsInFilterNots
			return (iIsImportant && !jIsImportant) || (iIsImportant == jIsImportant && analyzerNames[i] <= analyzerNames[j])
		})

		r, err := repparser.ParseFile(replay)
		if err != nil {
			log.Printf("Failed to parse replay: %v\n", err)
			continue
		}
		tryCompute(r)

		var results = map[string]string{}
		for _, name := range analyzerNames {
			a, ok := analyzerInstances[name]
			if !ok { // N.B. Analyzer might have been removed
				continue
			}
			_, isInFilters := filters[name]
			_, isInFilterNots := filterNots[name]
			err, done := a.StartReadingReplay(r, ctx, replay)
			if err != nil {
				log.Printf("Error beginning to read replay %v with Analyzer %v: %v\n", replay, name, err)
			}
			if done {
				results[name], _ = a.IsDone()
			}
			if done && ((isInFilterNots && results[name] == "true") || (isInFilters && results[name] != "true")) {
				continue replayLoop // Optimization: move to next replay if excluded by any filters already
			}
			if done || err != nil {
				delete(analyzerInstances, name)
			}
		}
		for _, c := range r.Commands.Cmds { // N.B. This is the expensive loop in the algorithm; optimize here!
			if len(analyzerInstances) == 0 {
				break // Optimization: don't loop over commands if there's nothing to do!
			}
			for _, name := range analyzerNames {
				a, ok := analyzerInstances[name]
				if !ok { // N.B. Analyzer might have been removed
					continue
				}
				_, isInFilters := filters[name]
				_, isInFilterNots := filterNots[name]
				err, done := a.ProcessCommand(c)
				if err != nil {
					log.Printf("Error reading command on replay %v with Analyzer %v: %v\n", replay, a.Name(), err)
				}
				if done {
					results[name], _ = a.IsDone()
				}
				if done && ((isInFilterNots && results[name] == "true") || (isInFilters && results[name] != "true")) {
					continue replayLoop // Optimization: move to next replay if excluded by any filters already
				}
				if done || err != nil { // Optimization: delete Analyzers that finished or errored out
					delete(analyzerInstances, name)
				}
			}
		}

		// Outputs a line of result (i.e. results for one replay)
		if shouldCopyToOutputLocation {
			err := copyFile(replay, fmt.Sprintf("%v/%v", *fCopyToIfMatchesFilters, filepath.Base(replay)))
			if err != nil {
				log.Printf("Error copying replay with path %v to %v: %v\n", replay, *fCopyToIfMatchesFilters, err)
			}
		}
		if *fJSON {
			if !firstJSONRow {
				fmt.Println(",")
			}
			firstJSONRow = false
			row := map[string]string{}
			for _, field := range fieldNames {
				row[field] = results[field]
			}
			bs, _ := json.Marshal(row)
			fmt.Printf("%s", bs)
		} else {
			csvRow := make([]string, 0, len(fieldNames))
			for _, field := range fieldNames {
				csvRow = append(csvRow, results[field])
			}
			w.Write(csvRow)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	if *fJSON {
		fmt.Println("\n]")
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
