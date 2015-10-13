package changelog

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type Changelog struct {
	Versions []*Version
}

func (c *Changelog) String() string {
	return join(c.Versions, "\n\n") + "\n"
}

func (c *Changelog) getVersion(versionNum string) *Version {
	for _, v := range c.Versions {
		if v.Version == versionNum {
			return v
		}
	}
	return nil
}

func (c *Changelog) getSubsection(versionNum, subsectionName string) *Subsection {
	for _, s := range c.getVersion(versionNum).Subsections {
		if s.Name == subsectionName {
			return s
		}
	}
	return nil
}

func (c *Changelog) AddSubsection(versionNum string, subsection string) {
	version := c.getVersion(versionNum)
	version.Subsections = append(version.Subsections, &Subsection{Name: subsection})
}

func (c *Changelog) AddLineToVersion(versionNum string, line *ChangeLine) {
	c.addToChangelines(&c.getVersion(versionNum).History, line)
}

func (c *Changelog) AddLineToSubsection(versionNum, subsectionName string, line *ChangeLine) {
	s := c.getSubsection(versionNum, subsectionName)
	c.addToChangelines(&s.History, line)
}

func (c *Changelog) addToChangelines(lines *[]*ChangeLine, line *ChangeLine) {
	*lines = append(*lines, line)
}

type Version struct {
	Version     string
	Date        string
	History     []*ChangeLine
	Subsections []*Subsection
}

func (v *Version) String() string {
	out := fmt.Sprintf("## %s", v.Version)
	if v.Date != "" {
		out += " / " + v.Date
	}
	if len(v.History) > 0 {
		out += "\n\n" + join(v.History, "\n")
	}
	if len(v.Subsections) > 0 {
		out += "\n\n" + join(v.Subsections, "\n\n")
	}
	return out
}

type Subsection struct {
	Name    string
	History []*ChangeLine
}

func (s *Subsection) String() string {
	if len(s.History) > 0 {
		return fmt.Sprintf(
			"### %s\n\n%s",
			s.Name,
			join(s.History, "\n"),
		)
	} else {
		return ""
	}
}

type ChangeLine struct {
	Summary   string
	Reference string
}

func (l *ChangeLine) String() string {
	if l.Reference == "" {
		return fmt.Sprintf(
			"  * %s",
			l.Summary,
		)
	} else {
		return fmt.Sprintf(
			"  * %s (%s)",
			l.Summary,
			l.Reference,
		)
	}
}

func join(lines interface{}, sep string) string {
	s := reflect.ValueOf(lines)
	if s.Kind() != reflect.Slice {
		panic("join given a non-slice type")
	}

	ret := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		vals := s.Index(i).MethodByName("String").Call(nil)
		ret[i] = vals[0].String()
	}

	return strings.Join(ret, sep)
}

func NewChangelog(filename string) (*Changelog, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return NewChangelogFromReader(file)
}

func NewChangelogFromReader(file io.Reader) (*Changelog, error) {
	history := &Changelog{Versions: []*Version{}}
	err := parseChangelog(file, history)
	if err != nil {
		return nil, err
	}
	return history, nil
}
