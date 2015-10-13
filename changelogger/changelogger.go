package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/parkr/changelog"
)

func main() {
	// Read options
	var filename string
	flag.StringVar(&filename, "file", "", "The path to your changelog")
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

	fmt.Fprintf(os.Stderr, "%s", history)

	// Add Changelog entry to correct part
	// Write History.markdown
}
