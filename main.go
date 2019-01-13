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
	executor, fQuiet, errs := buildAnalyzerExecutor(os.Args[1:])
	if len(errs) > 0 {
		if !fQuiet {
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
		if !fQuiet {
			log.Println("Errors encountered while executing analyzers:")
			log.Println()
			for _, err := range errs {
				log.Println(err.Error())
			}
		}
		return
	}
}

func buildAnalyzerExecutor(args []string) (*analyzer.Executor, bool, []error) {
	fs, stringFlags, boolFlags, fHelp, fOutput, fQuiet, fCopyToIfMatchesFilters := newFlagSet()
	fs.Parse(args)
	if *fHelp {
		flag.Usage()
		os.Exit(0)
	}

	var output analyzer.Output
	switch *fOutput {
	case "json":
		output = analyzer.NewJSONOutput(os.Stdout)
	case "none":
		output = analyzer.NewNoOutput()
	default:
		output = analyzer.NewCSVOutput(os.Stdout)
	}

	executor, errs := analyzer.NewExecutor(
		resolveReplayPaths(fs),
		resolveAnalyzerRequests(stringFlags, boolFlags),
		resolveContext(fs),
		output,
		*fCopyToIfMatchesFilters,
	)
	return executor, *fQuiet, errs
}

func newFlagSet() (*flag.FlagSet, map[string]*string, map[string]*bool, *bool, *string, *bool, *string) {
	var (
		fs                      = flag.NewFlagSet("", flag.ExitOnError)
		stringFlags             = map[string]*string{}
		boolFlags               = map[string]*bool{}
		fOutput                 = fs.String("o", "csv", "output format {csv|json|none} default: csv")
		fCopyToIfMatchesFilters = fs.String("copy-to-if-matches-filters", "",
			"copy replay files matched by -filter-- and not matched by -filter--not-- filters to specified directory")
		fQuiet = fs.Bool("quiet", false, "don't print any errors (discouraged: note that you can silence with 2>/dev/null).")
		fHelp  = fs.Bool("help", false, "Returns help usage and exits.")
	)
	fs.String("replay", "", "(>= 1 replays required) path to replay file")
	fs.String("replays", "", "(>= 1 replays required) comma-separated paths to replay files")
	fs.String("replay-dir", "", "(>= 1 replays required) path to folder with replays (recursive)")
	fs.String("me", "", "comma-separated list of player names to identify as the main player")
	for name, a := range analyzer.Analyzers {
		if a.IsStringFlag() {
			stringFlags[name] = fs.String(name, "", a.Description())
			if a.IsBooleanResult() {
				stringFlags["filter--"+name] = fs.String("filter--"+name, "", "Filter for: "+a.Description())
				stringFlags["filter-not--"+name] = fs.String("filter-not--"+name, "", "Filter-Not for: "+a.Description())
			}
		} else {
			boolFlags[name] = fs.Bool(name, false, a.Description())
			if a.IsBooleanResult() {
				boolFlags["filter--"+name] = fs.Bool("filter--"+name, false, "Filter for: "+a.Description())
				boolFlags["filter-not--"+name] = fs.Bool("filter-not--"+name, false, "Filter-Not for: "+a.Description())
			}
		}
	}
	return fs, stringFlags, boolFlags, fHelp, fOutput, fQuiet, fCopyToIfMatchesFilters
}

func resolveContext(fs *flag.FlagSet) analyzer.Context {
	var fMe string
	if fs.Lookup("me") != nil {
		fMe = fs.Lookup("me").Value.String()
	}
	ctx := analyzer.Context{Me: map[string]struct{}{}}
	if len(fMe) > 0 {
		for _, name := range strings.Split(fMe, ",") {
			ctx.Me[strings.TrimSpace(name)] = struct{}{}
		}
	}
	return ctx
}

func resolveAnalyzerRequests(stringFlags map[string]*string, boolFlags map[string]*bool) [][]string {
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
	return analyzerRequests
}

func resolveReplayPaths(fs *flag.FlagSet) []string {
	var fReplay, fReplays, fReplayDir string
	if fs.Lookup("replay") != nil {
		fReplay = fs.Lookup("replay").Value.String()
	}
	if fs.Lookup("replays") != nil {
		fReplays = fs.Lookup("replays").Value.String()
	}
	if fs.Lookup("replay-dir") != nil {
		fReplayDir = fs.Lookup("replay-dir").Value.String()
	}

	if fReplay == "" && fReplays == "" && fReplayDir == "" {
		fReplayDir = "." // default to find replays recursively in the current directory
	}
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
