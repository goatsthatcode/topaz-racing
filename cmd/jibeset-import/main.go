// jibeset-import converts a Jibeset multi-boat track export file into the
// Topaz Racing V1 boats.json format.
//
// Usage:
//
//	jibeset-import --self 53234 --interval 600 input.txt
//	jibeset-import --self 53234 --race-date 2025-01-18 --output boats.json input.txt
//	jibeset-import --self 53234 --include 53234,7406 --output boats.json input.txt
//
// When a file contains multiple races (identified by different FINH start dates),
// use --race-date YYYY-MM-DD to select only the tracks for one specific race.
//
// Each boat in the file gets an auto-assigned ID (sail number lowercased), name
// (from the file), color (from the built-in palette), and boatType "unknown".
// Override any of these per-boat with --color, --name, and --boat-type flags,
// each in the form SAILNUM:value.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"goatsthatcode.github.io/topaz-racing/gpximport"
	"goatsthatcode.github.io/topaz-racing/jibesetimport"
)

// defaultPalette is assigned in order to non-self boats.
var defaultPalette = []string{
	"#ff8a5b", // orange
	"#7ef58a", // green
	"#c084fc", // purple
	"#ffd966", // amber
	"#ff6b9d", // pink
	"#4dd0e1", // teal
}

func main() {
	var (
		selfFlag     = flag.String("self", "", "sail number of isSelf boat (Loren's boat)")
		raceDateFlag = flag.String("race-date", "", "YYYY-MM-DD — only include tracks whose FINH start date matches (use when file contains multiple races)")
		intervalFlag = flag.Int("interval", 600, "minimum seconds between track points (0 = no downsampling)")
		outputFlag   = flag.String("output", "", "output path for boats.json (default: stdout)")
		mergeFlag    = flag.String("merge", "", "path to existing boats.json to merge into")
		includeFlag  = flag.String("include", "", "comma-separated sail numbers to include (default: all)")
		colorFlag    = multiFlag{}
		nameFlag     = multiFlag{}
		boatTypeFlag = multiFlag{}
	)
	flag.Var(&colorFlag, "color", "SAILNUM:#hexcolor — override color for a boat (repeatable)")
	flag.Var(&nameFlag, "name", "SAILNUM:Display Name — override display name for a boat (repeatable)")
	flag.Var(&boatTypeFlag, "boat-type", "SAILNUM:type — override boat type for a boat (repeatable)")
	flag.Parse()

	// Parse --race-date if supplied.
	var filterDate time.Time
	if *raceDateFlag != "" {
		var parseErr error
		filterDate, parseErr = time.Parse("2006-01-02", *raceDateFlag)
		if parseErr != nil {
			log.Fatalf("--race-date %q is not a valid YYYY-MM-DD date: %v", *raceDateFlag, parseErr)
		}
	}

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("usage: jibeset-import [flags] <input.txt>")
	}

	f, err := os.Open(args[0])
	if err != nil {
		log.Fatalf("opening input file: %v", err)
	}
	defer f.Close()

	tracks, err := jibesetimport.ParseFile(f)
	if err != nil {
		log.Fatalf("parsing jibeset file: %v", err)
	}
	fmt.Fprintf(os.Stderr, "parsed %d track(s) from %s\n", len(tracks), args[0])
	for _, t := range tracks {
		fmt.Fprintf(os.Stderr, "  sail=%s name=%q start=%s points=%d\n",
			t.SailNumber, t.Name, t.StartTime.Format("2006-01-02"), len(t.Points))
	}

	// Build include set.
	includeSet := map[string]bool{}
	if *includeFlag != "" {
		for _, s := range strings.Split(*includeFlag, ",") {
			includeSet[strings.TrimSpace(s)] = true
		}
	}

	var boatsFile gpximport.BoatsFile
	if *mergeFlag != "" {
		data, readErr := os.ReadFile(*mergeFlag)
		if readErr != nil {
			log.Fatalf("reading merge file %q: %v", *mergeFlag, readErr)
		}
		if parseErr := json.Unmarshal(data, &boatsFile); parseErr != nil {
			log.Fatalf("parsing merge file %q: %v", *mergeFlag, parseErr)
		}
	}

	paletteIdx := 0
	for _, track := range tracks {
		// Filter by race date if --race-date was specified.
		if !filterDate.IsZero() {
			trackDate := track.StartTime.UTC().Truncate(24 * time.Hour)
			wantDate := filterDate.UTC().Truncate(24 * time.Hour)
			if !trackDate.Equal(wantDate) {
				continue
			}
		}

		if len(includeSet) > 0 && !includeSet[track.SailNumber] {
			fmt.Fprintf(os.Stderr, "  skipping sail=%s (not in --include)\n", track.SailNumber)
			continue
		}

		isSelf := *selfFlag != "" && track.SailNumber == *selfFlag

		color := colorFlag.lookup(track.SailNumber)
		if color == "" {
			if isSelf {
				color = "#4fd1ff"
			} else {
				color = defaultPalette[paletteIdx%len(defaultPalette)]
				paletteIdx++
			}
		}

		opts := gpximport.BoatOptions{
			ID:       strings.ToLower(strings.ReplaceAll(track.SailNumber, " ", "-")),
			Name:     nameFlag.lookup(track.SailNumber),
			Color:    color,
			BoatType: boatTypeFlag.lookup(track.SailNumber),
			IsSelf:   isSelf,
		}

		boat, convertErr := jibesetimport.ConvertTrack(track, opts, *intervalFlag)
		if convertErr != nil {
			log.Fatalf("converting track sail=%s: %v", track.SailNumber, convertErr)
		}
		gpximport.MergeBoat(&boatsFile, boat)
		fmt.Fprintf(os.Stderr, "  converted sail=%s → id=%s points=%d\n",
			track.SailNumber, boat.ID, len(boat.Track))
	}

	if len(boatsFile.Boats) == 0 {
		log.Fatal("no boats were produced; check --include and input file")
	}

	out, err := json.MarshalIndent(boatsFile, "", "  ")
	if err != nil {
		log.Fatalf("encoding boats.json: %v", err)
	}
	out = append(out, '\n')

	if *outputFlag != "" {
		if writeErr := os.WriteFile(*outputFlag, out, 0o644); writeErr != nil {
			log.Fatalf("writing output file %q: %v", *outputFlag, writeErr)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", *outputFlag)
	} else {
		os.Stdout.Write(out)
	}
}

// multiFlag is a flag.Value that collects SAILNUM:value pairs.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ", ") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// lookup returns the value for the given sail number key, or "" if not set.
func (m multiFlag) lookup(sailNum string) string {
	prefix := sailNum + ":"
	for _, entry := range m {
		if strings.HasPrefix(entry, prefix) {
			return strings.TrimPrefix(entry, prefix)
		}
	}
	return ""
}
