// gpx-import converts a GPX track file into the Topaz Racing V1 boats.json format.
//
// Usage:
//
//	gpx-import --id topaz --name "Topaz" --color "#4fd1ff" --self input.gpx
//	gpx-import --id wildcard --name "Wildcard" --color "#ff8a5b" --merge existing-boats.json input.gpx
//
// The converted file is written to stdout unless --output is specified.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"goatsthatcode.github.io/topaz-racing/gpximport"
)

func main() {
	var (
		id        = flag.String("id", "", "boat ID (required)")
		name      = flag.String("name", "", "boat display name (defaults to GPX track name or boat ID)")
		color     = flag.String("color", "", "hex color for this boat, e.g. #4fd1ff (default: #4fd1ff)")
		boatType  = flag.String("boat-type", "", "boat class or type, e.g. \"Express 27\" (default: unknown)")
		self      = flag.Bool("self", false, "mark this boat as isSelf (Loren's boat)")
		mergeFile = flag.String("merge", "", "path to existing boats.json to merge this boat into")
		output    = flag.String("output", "", "output path for boats.json (default: stdout)")
	)
	flag.Parse()

	if *id == "" {
		log.Fatal("--id is required")
	}

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("usage: gpx-import [flags] <input.gpx>")
	}

	f, err := os.Open(args[0])
	if err != nil {
		log.Fatalf("opening GPX file: %v", err)
	}
	defer f.Close()

	boat, err := gpximport.ConvertGPX(f, gpximport.BoatOptions{
		ID:       *id,
		Name:     *name,
		Color:    *color,
		BoatType: *boatType,
		IsSelf:   *self,
	})
	if err != nil {
		log.Fatalf("converting GPX: %v", err)
	}

	var boatsFile gpximport.BoatsFile

	if *mergeFile != "" {
		data, err := os.ReadFile(*mergeFile)
		if err != nil {
			log.Fatalf("reading merge file %q: %v", *mergeFile, err)
		}
		if err := json.Unmarshal(data, &boatsFile); err != nil {
			log.Fatalf("parsing merge file %q: %v", *mergeFile, err)
		}
	}

	gpximport.MergeBoat(&boatsFile, boat)

	out, err := json.MarshalIndent(boatsFile, "", "  ")
	if err != nil {
		log.Fatalf("encoding boats.json: %v", err)
	}
	out = append(out, '\n')

	if *output != "" {
		if err := os.WriteFile(*output, out, 0o644); err != nil {
			log.Fatalf("writing output file %q: %v", *output, err)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", *output)
	} else {
		os.Stdout.Write(out)
	}
}
