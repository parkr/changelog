package changelog

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

var (
	versionRegexp           = regexp.MustCompile(`##? (?i:(HEAD|v?\d+\.\d+(?:\.\d+)?))(?U:.*)(\d{4}-\d{2}-\d{2})?.?$`)
	subheaderRegexp         = regexp.MustCompile(`### ([0-9A-Za-z_ ]+)`)
	changeLineRegexp        = regexp.MustCompile(`[\*|\-] (.+)`)
	changeLineRegexpWithRef = regexp.MustCompile(`[\*|\-] (.+)( \(((#[0-9]+)|(@?[[:word:]]+))\))`)

	verbose = false
)

// SetVerbose sets the verbose flag to the value passed.
// If true is passed, verbose logging will be enabled.
func SetVerbose(v bool) {
	verbose = v
}

func logVerbose(args ...interface{}) {
	if verbose == true {
		log.Println(args...)
	}
}

func matchLine(regexp *regexp.Regexp, line string) (matches []string, doesMatch bool) {
	if regexp.MatchString(line) {
		return regexp.FindStringSubmatch(line), true
	}
	return nil, false
}

func versionDateFromMatches(matches []string) string {
	var date string
	if len(matches) == 3 {
		date = matches[2]
	}
	return date
}

func parseChangelog(file io.Reader, history *Changelog) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	currentHeader := ""
	currentSubHeader := ""
	var currentLine *ChangeLine
	for scanner.Scan() {
		txt := scanner.Text()
		logVerbose(txt)
		logVerbose("isHeader", versionRegexp.MatchString(txt))
		if matches, ok := matchLine(versionRegexp, txt); ok {
			logVerbose("headerMatches:", matches, len(matches))
			currentHeader = matches[1]
			currentSubHeader = ""
			logVerbose("currentHeader:", currentHeader)
			newVersion := history.GetVersionOrCreate(currentHeader)
			newVersion.Date = versionDateFromMatches(matches)
			continue
		}

		logVerbose("isSubHeader", subheaderRegexp.MatchString(txt))
		if matches, ok := matchLine(subheaderRegexp, txt); ok {
			logVerbose("subHeaderMatches:", matches, len(matches))
			currentSubHeader = matches[1]
			logVerbose("currentSubHeader:", currentSubHeader)
			history.GetSubsectionOrCreate(currentHeader, currentSubHeader)
			continue
		}

		logVerbose("isChangeLine", changeLineRegexp.MatchString(txt))
		if matches, ok := matchLine(changeLineRegexp, txt); ok {
			logVerbose("changeLineMatches:", matches, len(matches))
			var line *ChangeLine
			if more, ok := matchLine(changeLineRegexpWithRef, txt); ok {
				// Has ref
				line = &ChangeLine{
					Summary:   more[1],
					Reference: more[3],
				}
			} else {
				// No ref
				line = &ChangeLine{
					Summary: matches[1],
				}
			}
			logVerbose("newChangeLine:", line)
			currentLine = line
			if currentSubHeader == "" {
				history.AddLineToVersion(currentHeader, line)
			} else {
				history.AddLineToSubsection(currentHeader, currentSubHeader, line)
			}
			continue
		} else {
			if strings.TrimSpace(txt) != "" && currentLine != nil {
				currentLine.Summary += "\n\n" + txt
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error reading history:", err)
	}
	return nil
}
