package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
			(&analyzer.MyAPM{}).Name():                      &analyzer.MyAPM{},
			(&analyzer.MyRace{}).Name():                     &analyzer.MyRace{},
			(&analyzer.DateTime{}).Name():                   &analyzer.DateTime{},
			(&analyzer.DurationMinutes{}).Name():            &analyzer.DurationMinutes{},
			(&analyzer.MyName{}).Name():                     &analyzer.MyName{},
			(&analyzer.IsThereARace{}).Name():               &analyzer.IsThereARace{},
			(&analyzer.MyRaceIsZerg{}).Name():               &analyzer.MyRaceIsZerg{},
			(&analyzer.MyRaceIsTerran{}).Name():             &analyzer.MyRaceIsTerran{},
			(&analyzer.MyRaceIsProtoss{}).Name():            &analyzer.MyRaceIsProtoss{},
			(&analyzer.ReplayName{}).Name():                 &analyzer.ReplayName{},
			(&analyzer.ReplayPath{}).Name():                 &analyzer.ReplayPath{},
			(&analyzer.MyWin{}).Name():                      &analyzer.MyWin{},
			(&analyzer.MyGame{}).Name():                     &analyzer.MyGame{},
			(&analyzer.MapName{}).Name():                    &analyzer.MapName{},
			(&analyzer.MyFirstSpecificUnitSeconds{}).Name(): &analyzer.MyFirstSpecificUnitSeconds{},
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
			if strings.HasPrefix(name, "filter--") {
				filters[name[len("filter--"):]] = struct{}{}
			} else if strings.HasPrefix(name, "filter-not--") {
				filterNots[name[len("filter-not--"):]] = struct{}{}
			} else {
				analyzers[name] = _analyzers[name]
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
	// TODO break if there are no Analyzers at the beginning or after an iteration
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

		var results = map[string]string{}
		for name, a := range analyzerInstances {
			err, done := a.StartReadingReplay(r, ctx, replay)
			if err != nil {
				log.Printf("Error beginning to read replay %v with Analyzer %v: %v\n", replay, a.Name(), err)
			}
			if done {
				results[name], _ = a.IsDone()
			}
			if done || err != nil {
				delete(analyzerInstances, name)
			}
		}
		for _, c := range r.Commands.Cmds { // N.B. This is the expensive loop in the algorithm; optimize here!
			if len(analyzerInstances) == 0 {
				break // Optimization: don't loop over commands if there's nothing to do!
			}
			for name, a := range analyzerInstances {
				err, done := a.ProcessCommand(c)
				if err != nil {
					log.Printf("Error reading command on replay %v with Analyzer %v: %v\n", replay, a.Name(), err)
				}
				if done {
					results[name], _ = a.IsDone()
				}
				if done || err != nil { // Optimization: delete Analyzers that finished or errored out
					delete(analyzerInstances, name)
				}
			}
		}

		// Decides if this replay should be output based on filter flags
		shouldShowBasedOnFilterNots := true
		for filterNot := range filterNots {
			if res, ok := results[filterNot]; ok && res == "true" {
				shouldShowBasedOnFilterNots = false
			}
		}
		shouldShowBasedOnFilters := true
		for filter := range filters {
			if res, ok := results[filter]; !ok || res != "true" {
				shouldShowBasedOnFilters = false
			}
		}

		// Outputs a line of result (i.e. results for one replay)
		if shouldShowBasedOnFilterNots && shouldShowBasedOnFilters {
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

func unmarshalArguments(s string) []string {
	ss := []string{}
	for _, _si := range strings.Split(s, ",") {
		si := strings.TrimSpace(_si)
		if si != "" {
			ss = append(ss, si)
		}
	}
	return ss
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	err = copyFileContents(src, dst)
	return
}

// https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// https://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-exists
// exists returns whether the given file or directory exists
func isDirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
