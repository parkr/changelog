package changelog

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertSortOrder(t *testing.T, history *Changelog, header string, expected int) {
	assert.Equal(t, expected, history.GetVersion(header).sortOrder, "Wrong sortOrder for header: %q", header)
}

func TestChangelogString_Simplest(t *testing.T) {
	changelog := NewChangelog()
	changelog.AddLineToVersion("", &ChangeLine{
		Summary:   "summary 1",
		Reference: "reference 1",
	})

	actual := changelog.String()

	expected := `  * summary 1 (reference 1)` + "\n"
	assert.Equal(t, expected, actual)
}

func TestChangelogString_Complex(t *testing.T) {
	history := NewChangelog()
	history.GetVersionOrCreate("1.2.3").Date = "2022-10-10"
	history.AddLineToVersion("", &ChangeLine{
		Summary:   "summary 1",
		Reference: "#1",
	})
	history.AddLineToVersion("1.2.3", &ChangeLine{Summary: "summary 2", Reference: "#2"})
	history.AddLineToVersion("1.2.3", &ChangeLine{Summary: "summary 3", Reference: "#3"})
	history.AddLineToVersion("3.2.1", &ChangeLine{Summary: "summary 4", Reference: "#4"})
	history.AddLineToSubsection("3.2.1", "Subsection A", &ChangeLine{Summary: "summary 5", Reference: "#5"})

	actual := history.String()

	expected, err := os.ReadFile("testdata/changelog-string-complex.md")
	assert.NoError(t, err)
	assertSortOrder(t, history, "", -1)
	assertSortOrder(t, history, "1.2.3", 1)
	assertSortOrder(t, history, "3.2.1", 2)
	assert.Equal(t, string(expected), actual)
}

func TestChangelog_ParsesWhatItWrites(t *testing.T) {
	history := NewChangelog()
	history.GetVersionOrCreate("1.2.3").Date = "2022-10-10"
	history.AddLineToVersion("", &ChangeLine{
		Summary:   "summary 1",
		Reference: "#1",
	})
	history.AddLineToVersion("1.2.3", &ChangeLine{Summary: "summary 2", Reference: "#2"})
	history.AddLineToVersion("1.2.3", &ChangeLine{Summary: "summary 3", Reference: "#3"})
	history.AddLineToVersion("3.2.1", &ChangeLine{Summary: "summary 4", Reference: "#4"})
	history.AddLineToSubsection("3.2.1", "Subsection A", &ChangeLine{Summary: "summary 5", Reference: "#5"})
	history.AddLineToSubsection("5.4.1", "Subsection C", &ChangeLine{Summary: "summary 6", Reference: "#6"})
	history.AddLineToVersion("", &ChangeLine{Summary: "summary 7"})
	history.AddLineToVersion("", &ChangeLine{Summary: "summary 8"})
	history.AddLineToVersion("HEAD", &ChangeLine{Summary: "summary 9"})
	history.AddLineToVersion("HEAD", &ChangeLine{Summary: "summary 10"})

	actual := NewChangelog()
	parseChangelog(strings.NewReader(history.String()), actual)

	assertSortOrder(t, history, "", -1)
	assertSortOrder(t, history, "HEAD", 0)
	assertSortOrder(t, history, "1.2.3", 1)
	assertSortOrder(t, history, "3.2.1", 2)
	assertSortOrder(t, history, "5.4.1", 3)
	assert.Equal(t, history, actual, "expected:\n%q\nactual:\n%q\n", history.String(), actual.String())
}

func TestChangelog_WritesWhatItParses(t *testing.T) {
	expected, err := os.ReadFile("testdata/History.markdown")
	fd, err := os.Open("testdata/History.markdown")
	assert.NoError(t, err)
	history := NewChangelog()

	assert.NoError(t, parseChangelog(fd, history))
	actual := history.String()

	assert.Equal(t, string(expected), actual)

	wfd, _ := os.Create("testdata/History-rewritten.md")
	fmt.Fprintf(wfd, actual)
}
