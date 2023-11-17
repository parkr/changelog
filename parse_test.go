package changelog

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRegexpOutput struct {
	text    string
	matched []string
}

var (
	versions = []testRegexpOutput{
		{
			text:    "# HEAD",
			matched: []string{"# HEAD", "HEAD", ""},
		},
		{
			text:    "# [Unreleased]",
			matched: []string{"# [Unreleased]", "[Unreleased]", ""},
		},
		{
			text:    "# 1.0.0",
			matched: []string{"# 1.0.0", "1.0.0", ""},
		},
		{
			text:    "# 80.92.12 / 2015-02-30",
			matched: []string{"# 80.92.12 / 2015-02-30", "80.92.12", "2015-02-30"},
		},
		{
			text:    "# 80.92.12 - 2015-02-30",
			matched: []string{"# 80.92.12 - 2015-02-30", "80.92.12", "2015-02-30"},
		},
		{
			text:    "# 80.92.12 (2015-02-30)",
			matched: []string{"# 80.92.12 (2015-02-30)", "80.92.12", "2015-02-30"},
		},
		{
			text:    "# v0.6",
			matched: []string{"# v0.6", "v0.6", ""},
		},
		{
			text:    " # v0.6 / 2015-02-30",
			matched: []string{"# v0.6 / 2015-02-30", "v0.6", "2015-02-30"},
		},
		{
			text:    " # v0.6 / 2015-02-30",
			matched: []string{"# v0.6 / 2015-02-30", "v0.6", "2015-02-30"},
		},
		{
			text:    "## HEAD",
			matched: []string{"## HEAD", "HEAD", ""},
		},
		{
			text:    "## 1.0.0",
			matched: []string{"## 1.0.0", "1.0.0", ""},
		},
		{
			text:    "## 80.92.12 / 2015-02-30",
			matched: []string{"## 80.92.12 / 2015-02-30", "80.92.12", "2015-02-30"},
		},
		{
			text:    "## 80.92.12 - 2015-02-30",
			matched: []string{"## 80.92.12 - 2015-02-30", "80.92.12", "2015-02-30"},
		},
		{
			text:    "## 80.92.12 (2015-02-30)",
			matched: []string{"## 80.92.12 (2015-02-30)", "80.92.12", "2015-02-30"},
		},
		{
			text:    "## v0.6",
			matched: []string{"## v0.6", "v0.6", ""},
		},
		{
			text:    " ## v0.6 / 2015-02-30",
			matched: []string{"## v0.6 / 2015-02-30", "v0.6", "2015-02-30"},
		},
		{
			text:    " ## v0.6 - 2015-02-30",
			matched: []string{"## v0.6 - 2015-02-30", "v0.6", "2015-02-30"},
		},
		{
			text:    " ## v0.6 (2015-02-30)",
			matched: []string{"## v0.6 (2015-02-30)", "v0.6", "2015-02-30"},
		},
		{
			text:    "## v0.6.0-beta",
			matched: []string{"## v0.6.0-beta", "v0.6.0-beta", ""},
		},
		{
			text:    "## v0.6.0-beta - 2015-02-30",
			matched: []string{"## v0.6.0-beta - 2015-02-30", "v0.6.0-beta", "2015-02-30"},
		},
		{
			text:    "## [v0.6.0-beta] - 2015-02-30",
			matched: []string{"## [v0.6.0-beta] - 2015-02-30", "v0.6.0-beta", "2015-02-30"},
		},
		{
			text:    "## [1.0.0-rc.1+build.1] - 2015-02-30",
			matched: []string{"## [1.0.0-rc.1+build.1] - 2015-02-30", "1.0.0-rc.1+build.1", "2015-02-30"},
		},
	}
	subheaders = []testRegexpOutput{
		{
			text:    "### Minor",
			matched: []string{"### Minor", "Minor"},
		},
		{
			text:    " ### Minor Enhancements",
			matched: []string{"### Minor Enhancements", "Minor Enhancements"},
		},
	}
	changelines = []testRegexpOutput{
		{
			text:    "* I made a really cool change!",
			matched: []string{"* I made a really cool change!", "I made a really cool change!"},
		},
		{
			text:    "  * I made a really cool change!",
			matched: []string{"* I made a really cool change!", "I made a really cool change!"},
		},
		{
			text:    "  * The `coolest` change eVAR :smile: (#123)",
			matched: []string{"* The `coolest` change eVAR :smile: (#123)", "The `coolest` change eVAR :smile:", " (#123)", "#123", "#123", ""},
		},
		{
			text:    "  * The `coolest` change eVAR :smile: (abcdef23)",
			matched: []string{"* The `coolest` change eVAR :smile: (abcdef23)", "The `coolest` change eVAR :smile:", " (abcdef23)", "abcdef23", "", "abcdef23"},
		},
		{
			text:    "    * Fixed that narsty bug with tokenization (@carla)",
			matched: []string{"* Fixed that narsty bug with tokenization (@carla)", "Fixed that narsty bug with tokenization", " (@carla)", "@carla", "", "@carla"},
		},
		{
			text:    "- I made a really cool change!",
			matched: []string{"- I made a really cool change!", "I made a really cool change!"},
		},
		{
			text:    "  - I made a really cool change!",
			matched: []string{"- I made a really cool change!", "I made a really cool change!"},
		},
		{
			text:    "  - The `coolest` change eVAR :smile: (#123)",
			matched: []string{"- The `coolest` change eVAR :smile: (#123)", "The `coolest` change eVAR :smile:", " (#123)", "#123", "#123", ""},
		},
		{
			text:    "  - The `coolest` change eVAR :smile: (abcdef23)",
			matched: []string{"- The `coolest` change eVAR :smile: (abcdef23)", "The `coolest` change eVAR :smile:", " (abcdef23)", "abcdef23", "", "abcdef23"},
		},
		{
			text:    "    - Fixed that narsty bug with tokenization (@carla)",
			matched: []string{"- Fixed that narsty bug with tokenization (@carla)", "Fixed that narsty bug with tokenization", " (@carla)", "@carla", "", "@carla"},
		},
	}
	representativeChangelog = `## HEAD

### Major Enhancements

  * Liquid profiler (i.e. know how fast or slow your templates render) (#3762)
  * Iterate over site.collections as an array instead of a hash. (#3670)
  * Added permalink time variables (#3990)

### Minor Enhancements

  * Added basic microdata to post template in site template (#3189)

## 1.0 / 2012-02-03

  * I did some cool stuffs.

## v0.9 / 2012-01-01

* Birthday!!!!!
`
	trickyChangelog = `Hello there!

## [*Previously:*](https://example.com/previously?baz=boing)

- [ ] [The Mechanical Apple Watch | Watchfinder & Co. - YouTube](https://youtube.com/)
- [ ] [The Internet’s Own Example](https://example.com)
- [ ] [Little Shop of Horrors: Tiny Desk (Home) Concert - YouTube](https://m.youtube.com/watch?v=ymqKPz5kRXE)
- [ ] [D) Sector 4 (AQA)](https://www.ign.com/wikis/metroid-fusion/D)_Sector_4_(AQA))
- [ ] [Issue #123](https://github.com/octocat/mona/issues/123) (#123)

A link: https://example.com/2?foo=bar
`
)

func TestVersionRegexp(t *testing.T) {
	for _, version := range versions {
		assert.Regexp(t, versionRegexp, version.text)
		matches, ok := matchLine(versionRegexp, version.text)
		assert.True(t, ok, "regexp should match %q", version.text)
		assert.Equal(t, version.matched, matches)
	}
}

func TestSubheaderRegexp(t *testing.T) {
	for _, subheader := range subheaders {
		assert.Regexp(t, subheaderRegexp, subheader.text)
		matches, ok := matchLine(subheaderRegexp, subheader.text)
		assert.True(t, ok)
		assert.Equal(t, matches, subheader.matched)
	}
}

func TestChangelineRegexp(t *testing.T) {
	for _, changeline := range changelines {
		assert.Regexp(t, changeLineRegexp, changeline.text)
		if len(changeline.matched) > 5 {
			// Has ref
			assert.Regexp(t, changeLineRegexpWithRef, changeline.text)
			matches, ok := matchLine(changeLineRegexpWithRef, changeline.text)
			assert.True(t, ok)
			assert.Equal(t, matches, changeline.matched)
		} else {
			// No ref
			matches, ok := matchLine(changeLineRegexp, changeline.text)
			assert.True(t, ok)
			assert.Equal(t, matches, changeline.matched)
		}
	}
}

func TestParseChangelog(t *testing.T) {
	changes := NewChangelog()
	err := parseChangelog(strings.NewReader(representativeChangelog), changes)
	assert.NoError(t, err)
	assert.Len(t, changes.Versions, 3)

	// Test HEAD
	assert.Equal(t, "HEAD", changes.Versions[0].Version)
	assert.Len(t, changes.Versions[0].Subsections, 2)
	assert.Len(t, changes.Versions[0].Subsections[0].History, 3)
	assert.Equal(t,
		"Liquid profiler (i.e. know how fast or slow your templates render)",
		changes.Versions[0].Subsections[0].History[0].Summary,
	)
	assert.Equal(t, "#3762", changes.Versions[0].Subsections[0].History[0].Reference)
	assert.Equal(t,
		"Added basic microdata to post template in site template",
		changes.Versions[0].Subsections[1].History[0].Summary,
	)
	assert.Equal(t, "#3189", changes.Versions[0].Subsections[1].History[0].Reference)

	assert.Equal(t, "1.0", changes.Versions[1].Version)
	assert.Equal(t, "2012-02-03", changes.Versions[1].Date)
	assert.Len(t, changes.Versions[1].History, 1)

	assert.Equal(t, "v0.9", changes.Versions[2].Version)
	assert.Equal(t, "2012-01-01", changes.Versions[2].Date)
	assert.Len(t, changes.Versions[2].History, 1)
	assert.Equal(t, "Birthday!!!!!", changes.Versions[2].History[0].Summary)
}

func TestParseChangelog_WithoutHeaders(t *testing.T) {
	fd, err := os.Open("testdata/changelog-checkboxes-no-headers.md")
	if err != nil {
		t.Fatal(err)
	}
	expected := NewChangelog()
	expected.AddLineToVersion("", &ChangeLine{Summary: "[ ] [Issue 1 complete](https://github.com/foo/bar/issue/1)"})
	expected.AddLineToVersion("", &ChangeLine{Summary: "[ ] [Issue 2 complete, too!](https://github.com/foo/bar/issue/2)"})
	expected.AddLineToVersion("", &ChangeLine{Summary: "[ ] [Secretly Issue 3, but masquerading as issue 4](https://github.com/foo/bar/issue/3)\n\n[You can see more later sometime maybe!](https://hi.there/foo)"})
	changes := NewChangelog()

	err = parseChangelog(fd, changes)

	assert.NoError(t, err)
	assert.Equal(t, expected, changes)
}

func TestParseChangelog_TrickyInput(t *testing.T) {
	changes := NewChangelog()
	SetVerbose(true)

	err := parseChangelog(strings.NewReader(trickyChangelog), changes)

	assert.NoError(t, err)
	assert.Len(t, changes.Versions, 1)
	previouslyVersion := changes.Versions[0]
	assert.Equal(t, previouslyVersion, &Version{
		Version: "", // TODO: Think about whether I should try to take any h1/h2 as a header, or only specifically organized ones.
		History: []*ChangeLine{
			{Summary: "[ ] [The Mechanical Apple Watch | Watchfinder & Co. - YouTube](https://youtube.com/)", Reference: ""},
			{Summary: "[ ] [The Internet’s Own Example](https://example.com)", Reference: ""},
			{Summary: "[ ] [Little Shop of Horrors: Tiny Desk (Home) Concert - YouTube](https://m.youtube.com/watch?v=ymqKPz5kRXE)", Reference: ""},
			{Summary: "[ ] [D) Sector 4 (AQA)](https://www.ign.com/wikis/metroid-fusion/D)_Sector_4_(AQA))", Reference: ""},
			{Summary: "[ ] [Issue #123](https://github.com/octocat/mona/issues/123)\n\nA link: https://example.com/2?foo=bar", Reference: "#123"},
		},
		Subsections: []*Subsection{},
		sortOrder:   -1,
	})
}
