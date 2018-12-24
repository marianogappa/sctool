package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marianogappa/sctool/analyzer"
)

func main() {
	var (
		boolFlags               = map[string]*bool{}
		stringFlags             = map[string]*string{}
		fReplay                 = flag.String("replay", "", "(>= 1 replays required) path to replay file")
		fReplays                = flag.String("replays", "", "(>= 1 replays required) comma-separated paths to replay files")
		fReplayDir              = flag.String("replay-dir", "", "(>= 1 replays required) path to folder with replays (recursive)")
		fMe                     = flag.String("me", "", "comma-separated list of player names to identify as the main player")
		fJSON                   = flag.Bool("json", false, "outputs a JSON instead of the default CSV")
		fCopyToIfMatchesFilters = flag.String("copy-to-if-matches-filters", "",
			"copy replay files matched by -filter-- and not matched by -filter--not-- filters to specified directory")
		fQuiet = flag.Bool("quiet", false, "don't print any errors (discouraged: note that you can silence with 2>/dev/null).")
		fHelp  = flag.Bool("help", false, "Returns help usage and exits.")
	)
	for name, a := range analyzer.Analyzers {
		if a.IsStringFlag() {
			stringFlags[name] = flag.String(name, "", a.Description())
			if a.IsBooleanResult() {
				stringFlags["filter--"+name] = flag.String("filter--"+name, "", "Filter for: "+a.Description())
				stringFlags["filter-not--"+name] = flag.String("filter-not--"+name, "", "Filter-Not for: "+a.Description())
			}
		} else {
			boolFlags[name] = flag.Bool(name, false, a.Description())
			if a.IsBooleanResult() {
				boolFlags["filter--"+name] = flag.Bool("filter--"+name, false, "Filter for: "+a.Description())
				boolFlags["filter-not--"+name] = flag.Bool("filter-not--"+name, false, "Filter-Not for: "+a.Description())
			}
		}
	}
	flag.Parse()
	if *fHelp {
		flag.Usage()
		os.Exit(0)
	}

	ctx := analyzer.Context{Me: map[string]struct{}{}}
	if fMe != nil && len(*fMe) > 0 {
		for _, name := range strings.Split(*fMe, ",") {
			ctx.Me[strings.TrimSpace(name)] = struct{}{}
		}
	}
	var output analyzer.Output = analyzer.NewCSVOutput(os.Stdout)
	if *fJSON {
		output = analyzer.NewJSONOutput(os.Stdout)
	}
	replayPaths := resolveReplayPaths(*fReplay, *fReplays, *fReplayDir)

	analyzerRequests := [][]string{}
	for name, f := range boolFlags {
		if *f {
			analyzerRequests = append(analyzerRequests, []string{name})
		}
	}
	for name, f := range stringFlags {
		if *f != "" {
			s := []string{name}
			s = append(s, unmarshalArguments(*f)...)
			analyzerRequests = append(analyzerRequests, s)
		}
	}

	executor, errs := analyzer.NewExecutor(replayPaths, analyzerRequests, ctx, output, *fCopyToIfMatchesFilters)
	if len(errs) > 0 {
		if !*fQuiet {
			log.Println("Errors encountered while preparing to execute analyzers:")
			log.Println()
			for _, err := range errs {
				log.Println(err.Error())
			}
		}
		return
	}
	errs = executor.Execute()
	if len(errs) > 0 {
		if !*fQuiet {
			log.Println("Errors encountered while executing analyzers:")
			log.Println()
			for _, err := range errs {
				log.Println(err.Error())
			}
		}
		return
	}
}

func resolveReplayPaths(fReplay, fReplays, fReplayDir string) []string {
	var replays = map[string]struct{}{}
	fReplay = strings.TrimSpace(fReplay)
	if len(fReplay) >= 5 && (fReplay)[len(fReplay)-4:] == ".rep" {
		replays[fReplay] = struct{}{}
	}
	if fReplays != "" {
		for _, r := range strings.Split(fReplays, ",") {
			r = strings.TrimSpace(r)
			if len(r) >= 5 && r[len(r)-4:] == ".rep" {
				replays[r] = struct{}{}
			}
		}
	}
	if fReplayDir != "" {
		e := filepath.Walk(fReplayDir, func(path string, info os.FileInfo, err error) error {
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
	var res = []string{}
	for replay := range replays {
		res = append(res, replay)
	}
	return res
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
