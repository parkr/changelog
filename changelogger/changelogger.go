package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/parkr/changelog"
)

func main() {
	// Read options
	var filename string
	flag.StringVar(&filename, "file", "", "The path to your changelog")
	var output string
	flag.StringVar(&output, "out", "", "Where to write the changelog")
	flag.Parse()

	// Find History.markdown
	if filename == "" {
		filename = changelog.HistoryFilename()
	}

	// Read History.markdown
	history, err := changelog.NewChangelog(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Write History.markdown
	var writer io.Writer
	if output == "" {
		writer = os.Stderr
	} else {
		f, err := os.Create(output)
		if err != nil {
			log.Fatal(err)
		}
		writer = f
		defer f.Close()
	}
	fmt.Fprintf(writer, "%s", history)
}
