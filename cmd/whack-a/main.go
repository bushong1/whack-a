// Modified in 2026 from the original old-man-yells-at for whack-a purposes.
// Command whack-a creates a whack-a-mole emoji from a PNG image.
package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	whacka "github.com/bushong1/whack-a"
)

func main() {
	options := whacka.DefaultOptions()
	flags := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	flags.IntVar(&options.LeftHiddenPct, "left-hidden-pct", options.LeftHiddenPct, "percent of the lower-left emoji hidden behind its hole")
	flags.IntVar(&options.LeftHiddenPct, "l", options.LeftHiddenPct, "shortcut for --left-hidden-pct")
	flags.IntVar(&options.RightExposedPct, "right-exposed-pct", options.RightExposedPct, "percent of the right emoji exposed above its hole")
	flags.IntVar(&options.RightExposedPct, "r", options.RightExposedPct, "shortcut for --right-exposed-pct")
	if err := flags.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}
	if flags.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <emoji.png>\n", filepath.Base(os.Args[0]))
		flags.PrintDefaults()
		os.Exit(2)
	}

	inputName := flags.Arg(0)
	f, err := os.Open(inputName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading input file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	input, err := png.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decoding PNG: %v\n", err)
		os.Exit(1)
	}
	whacked, err := whacka.WhackAWithOptions(input, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid display options: %v\n", err)
		os.Exit(2)
	}

	base := filepath.Base(inputName)
	outputName := "whack-a-" + strings.TrimSuffix(base, filepath.Ext(base)) + ".png"
	output, err := os.Create(outputName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating output file: %v\n", err)
		os.Exit(1)
	}
	defer output.Close()
	if err := png.Encode(output, whacked); err != nil {
		fmt.Fprintf(os.Stderr, "writing output image: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Whack-a-mole created: %s\n", outputName)
}
