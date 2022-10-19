package changelog

import (
	"io/ioutil"
	"regexp"
)

var historyFilenameRegexp = regexp.MustCompile("(?i:(History|Changelog).m(ar)?k?d(own)?)")

// HistoryFilename discovers the correct filename for your history file
// based on files in the current working directory. It iterates through the
// files in your current directory looking for a file with some
// case-insensitive form of History.markdown or Changelog.markdown with any
// series of supported markdown file extensions.
func HistoryFilename() string {
	infos, err := ioutil.ReadDir(".")
	if err != nil {
		logFatal("Problem finding your history file.")
		// System exit!
	}
	for _, info := range infos {
		if isHistoryFile(info.Name()) {
			return info.Name()
		}
	}
	return "History.markdown"
}

// isHistoryFile checks whether a given filename is a valid changelog name.
func isHistoryFile(filename string) bool {
	return historyFilenameRegexp.FindString(filename) != ""
}
