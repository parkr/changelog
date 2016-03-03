package changelog

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewVersion(t *testing.T) {
	versionName := "HEAD"
	version := NewVersion(versionName)
	assert.Equal(t, versionName, version.Version)
	assert.Equal(t, "", version.Date)
	assert.NotNil(t, version.History)
	assert.Len(t, version.History, 0)
	assert.NotNil(t, version.Subsections)
	assert.Len(t, version.Subsections, 0)
}

func TestGetVersion(t *testing.T) {
	versionName := "HEAD"
	c := NewChangelog()
	assert.Nil(t, c.GetVersion(versionName))

	origVersion := NewVersion(versionName)
	c.Versions = append(c.Versions, origVersion)

	version := c.GetVersion(versionName)
	assert.NotNil(t, version)
	assert.Equal(t, versionName, version.Version)
	assert.Equal(t, origVersion, version)
}

func TestGetVersionOrCreate(t *testing.T) {
	versionName := "HEAD"
	c := NewChangelog()
	origVersion := NewVersion(versionName)

	version := c.GetVersionOrCreate(versionName)
	assert.NotNil(t, version)
	assert.Equal(t, versionName, version.Version)
	assert.False(t, origVersion == version)

	c = NewChangelog()
	c.Versions = append(c.Versions, origVersion)

	version = c.GetVersionOrCreate(versionName)
	assert.NotNil(t, version)
	assert.Equal(t, versionName, version.Version)
	assert.True(t, origVersion == version)
}

func TestNewSubsection(t *testing.T) {
	subsectionName := "Minor Enhancements"
	subsection := NewSubsection(subsectionName)
	assert.Equal(t, subsectionName, subsection.Name)
	assert.NotNil(t, subsection.History)
	assert.Len(t, subsection.History, 0)
}

func TestGetSubsection(t *testing.T) {
	versionNum := "HEAD"
	subsectionName := "Bug Fixes"

	c := NewChangelog()

	// Nil if no version.
	assert.Nil(t, c.GetVersion(versionNum))
	assert.Nil(t, c.GetSubsection(versionNum, subsectionName))

	// Nil if version but no subsection.
	version := c.GetVersionOrCreate(versionNum)
	assert.NotNil(t, c.GetVersion(versionNum))
	assert.Nil(t, c.GetSubsection(versionNum, subsectionName))

	// Not nil if version and subsection exist.
	subsection := NewSubsection(subsectionName)
	version.Subsections = append(version.Subsections, subsection)
	assert.NotNil(t, c.GetSubsection(versionNum, subsectionName))
	assert.Equal(t, subsection, c.GetSubsection(versionNum, subsectionName))
}

func TestGetSubsectionOrCreate(t *testing.T) {
	versionNum := "HEAD"
	subsectionName := "Bug Fixes"

	c := NewChangelog()

	// Not nil if no version.
	assert.Nil(t, c.GetVersion(versionNum))
	subsection := c.GetSubsectionOrCreate(versionNum, subsectionName)
	assert.NotNil(t, subsection)
	assert.Equal(t, subsectionName, subsection.Name)
	assert.Exactly(t, c.GetVersion(versionNum).Subsections[0], subsection)

	// Not nil if version but no subsection.
	c = NewChangelog()
	version := c.GetVersionOrCreate(versionNum)
	assert.NotNil(t, c.GetVersion(versionNum))
	assert.Exactly(t, version, c.GetVersion(versionNum))
	subsection = c.GetSubsectionOrCreate(versionNum, subsectionName)
	assert.NotNil(t, subsection)
	assert.Equal(t, subsectionName, subsection.Name)
	assert.Exactly(t, version.Subsections[0], subsection)

	// Not nil if version and subsection exist.
	c = NewChangelog()
	version = c.GetVersionOrCreate(versionNum)
	subsection = NewSubsection(subsectionName)
	version.Subsections = append(version.Subsections, subsection)
	subsection = c.GetSubsectionOrCreate(versionNum, subsectionName)
	assert.NotNil(t, subsection)
	assert.Equal(t, subsectionName, subsection.Name)
	assert.Exactly(t, version.Subsections[0], subsection)
}

func TestAddLineToVersion(t *testing.T) {
	versionNum := "HEAD"
	line := ChangeLine{Summary: "Added a fun method.", Reference: "#23"}

	c := NewChangelog()
	c.AddLineToVersion(versionNum, &line)
	assert.Equal(t, "## HEAD\n\n  * Added a fun method. (#23)\n", c.String())
}

func TestAddLineToSubsection(t *testing.T) {
	versionNum := "HEAD"
	subsectionName := "Bug Fixes"
	line := ChangeLine{Summary: "Added a fun method.", Reference: "#23"}

	c := NewChangelog()
	c.AddLineToSubsection(versionNum, subsectionName, &line)
	assert.Equal(t, "## HEAD\n\n### Bug Fixes\n\n  * Added a fun method. (#23)\n", c.String())
}
