package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/toros100/easynbt/easynbt/generator"
)

func main() {
	// TODO: set up a nicer CLI (cobra)
	// add optional target package(s) pattern and build flags args

	var types string
	flag.StringVar(&types, "types", "", "Comma-separated list of target types to generate unmarshalling code for")

	var out string
	flag.StringVar(&out, "out", "", "Output file path: optional, default is {package name}_nbt_gen.go in the working directory")

	var dry bool
	flag.BoolVar(&dry, "dry", false, "Dry run: no output file produced")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Verbose: additional logs, e.g. related to non-fatal errors")

	var version bool
	flag.BoolVar(&version, "version", false, "Print version and exit")

	flag.Parse()

	if version {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Fprintln(os.Stderr, "Failed to read build info")
			os.Exit(1)
		}

		fmt.Println("easynbt", info.Main.Version)
		fmt.Println("built with", info.GoVersion)
		os.Exit(0)
	}

	if len(types) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	t := strings.Split(types, ",")
	g := generator.New()

	opts := &generator.Options{
		OutFile: out,
		DryRun:  dry,
		Verbose: verbose,
	}

	err := g.Generate(opts, ".", t)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate: %v\n", err)
		os.Exit(1)
	}

}
