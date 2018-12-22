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
			(&analyzer.MyAPM{}).Name():           &analyzer.MyAPM{},
			(&analyzer.MyRace{}).Name():          &analyzer.MyRace{},
			(&analyzer.DateTime{}).Name():        &analyzer.DateTime{},
			(&analyzer.DurationMinutes{}).Name(): &analyzer.DurationMinutes{},
			(&analyzer.MyName{}).Name():          &analyzer.MyName{},
			(&analyzer.IsThereAZerg{}).Name():    &analyzer.IsThereAZerg{},
			(&analyzer.IsThereATerran{}).Name():  &analyzer.IsThereATerran{},
			(&analyzer.IsThereAProtoss{}).Name(): &analyzer.IsThereAProtoss{},
			(&analyzer.MyRaceIsZerg{}).Name():    &analyzer.MyRaceIsZerg{},
			(&analyzer.MyRaceIsTerran{}).Name():  &analyzer.MyRaceIsTerran{},
			(&analyzer.MyRaceIsProtoss{}).Name(): &analyzer.MyRaceIsProtoss{},
			(&analyzer.ReplayName{}).Name():      &analyzer.ReplayName{},
			(&analyzer.ReplayPath{}).Name():      &analyzer.ReplayPath{},
		}
		flags      = map[string]*bool{}
		fReplay    = flag.String("replay", "", "path to replay file")
		fReplays   = flag.String("replays", "", "comma-separated paths to replay files")
		fReplayDir = flag.String("replay-dir", "", "path to folder with replays (recursive)")
		fMe        = flag.String("me", "", "comma-separated list of player names to identify as the main player")
		fJSON      = flag.Bool("json", false, "outputs a JSON instead of the default CSV")
	)
	for name, a := range _analyzers {
		flags[name] = flag.Bool(name, false, a.Description())
	}
	flag.Parse()
	var (
		analyzers     = map[string]analyzer.Analyzer{}
		csvFieldNames = []string{} // TODO add filename
	)
	for name, f := range flags {
		if *f {
			analyzers[name] = _analyzers[name]
			csvFieldNames = append(csvFieldNames, name)
		}
	}

	// Prepares for CSV output
	sort.Strings(csvFieldNames)
	w := csv.NewWriter(os.Stdout)
	if !*fJSON {
		w.Write(csvFieldNames)
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
	for replay := range replays {
		analyzerInstances := make(map[string]analyzer.Analyzer, len(analyzers))
		for n, a := range analyzers {
			analyzerInstances[n] = a
		}

		r, err := repparser.ParseFile(replay)
		if err != nil {
			log.Printf("Failed to parse replay: %v\n", err)
			continue
		}
		tryCompute(r)

		var results = map[string]analyzer.Result{}
		for name, a := range analyzerInstances {
			if a.StartReadingReplay(r, ctx, replay) {
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

		if *fJSON {
			if !firstJSONRow {
				fmt.Println(",")
			}
			firstJSONRow = false
			row := map[string]string{}
			for _, field := range csvFieldNames {
				row[field] = results[field].Value()
			}
			bs, _ := json.Marshal(row)
			fmt.Printf("%s", bs)
		} else {
			csvRow := make([]string, 0, len(csvFieldNames))
			for _, field := range csvFieldNames {
				if results[field] == nil {
					csvRow = append(csvRow, "")
					continue
				}
				csvRow = append(csvRow, results[field].Value())
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
