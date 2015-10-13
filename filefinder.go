package changelog

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var historyFilenameRegexp = regexp.MustCompile("(?i:(History|Changelog).m(ar)?k?d(own)?)")

func HistoryFilename() string {
	infos, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Problem finding your history file.")
		os.Exit(1)
	}
	for _, info := range infos {
		if isHistoryFile(info.Name()) {
			return info.Name()
		}
	}
	return "History.markdown"
}

func isHistoryFile(filename string) bool {
	return historyFilenameRegexp.FindString(filename) != ""
}
