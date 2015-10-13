package changelog

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

var (
	versionRegexp    = regexp.MustCompile(`## (?i:(HEAD|\d+.\d+(.\d+)?)( / (\d{4}-\d{2}-\d{2}))?)`)
	subheaderRegexp  = regexp.MustCompile(`### ([0-9A-Za-z_ ]+)`)
	changeLineRegexp = regexp.MustCompile(`\* (.+)( \(((#[0-9]+)|(@[[:word:]]+))\))?`)
)

func parseChangelog(file io.Reader, history *Changelog) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	currentHeader := ""
	currentSubHeader := ""
	var currentLine *ChangeLine
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Println(txt)
		log.Println("isHeader", versionRegexp.MatchString(txt))
		if versionRegexp.MatchString(txt) {
			matches := versionRegexp.FindStringSubmatch(txt)
			log.Println("headerMatches:", matches, len(matches))
			currentHeader = matches[1]
			currentSubHeader = ""
			log.Printf("currentHeader: '%s'", currentHeader)
			var date string
			if len(matches) == 5 {
				date = matches[4]
			}
			history.Versions = append(history.Versions, &Version{
				Version:     currentHeader,
				Date:        date,
				History:     []*ChangeLine{},
				Subsections: []*Subsection{},
			})
			continue
		}

		log.Println("isSubHeader", subheaderRegexp.MatchString(txt))
		if subheaderRegexp.MatchString(txt) {
			matches := subheaderRegexp.FindStringSubmatch(txt)
			log.Println("subHeaderMatches:", matches, len(matches))
			currentSubHeader = matches[1]
			log.Printf("currentSubHeader: '%s'", currentSubHeader)
			history.AddSubsection(currentHeader, currentSubHeader)
			continue
		}

		if changeLineRegexp.MatchString(txt) {
			matches := changeLineRegexp.FindStringSubmatch(txt)
			log.Println("changeLineMatches:", matches, len(matches))
			line := &ChangeLine{
				Summary:   matches[1],
				Reference: matches[3],
			}
			log.Println("newChangeLine:", line)
			currentLine = line
			if currentSubHeader == "" {
				history.AddLineToVersion(currentHeader, line)
			} else {
				history.AddLineToSubsection(currentHeader, currentSubHeader, line)
			}
			continue
		} else {
			if strings.TrimSpace(txt) != "" && currentLine != nil {
				currentLine.Summary += " " + strings.TrimSpace(txt)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error reading history:", err)
	}
	return nil
}
