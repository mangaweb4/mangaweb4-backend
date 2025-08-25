package tag

import (
	"os"
	"strings"
	"testing"

	"github.com/mangaweb4/mangaweb4-backend/configuration"
	"github.com/stretchr/testify/suite"
)

type ParseTagNoFirstLevelTestSuite struct {
	suite.Suite
}

func TestParseTagNoFirstLevelTestSuite(t *testing.T) {
	suite.Run(t, new(ParseTagNoFirstLevelTestSuite))
}

func (suite *ParseTagNoFirstLevelTestSuite) SetupSuite() {
	configuration.Init(configuration.Config{
		FirstLevelDirAsTag: false,
	})
}

func (s *ParseTagNoFirstLevelTestSuite) TestTagOne() {
	s.Assert().ElementsMatch(
		ParseTag("[Test]Some weird name"),
		[]string{"Test"})
}

func (s *ParseTagNoFirstLevelTestSuite) TestTagNoTag() {
	s.Assert().Empty(ParseTag("Hello World"))
}

func (s *ParseTagNoFirstLevelTestSuite) TestTagMultiple() {
	s.Assert().ElementsMatch(
		ParseTag("[Test]Some weird name [Download]"),
		[]string{"Test", "Download"})
}

func (s *ParseTagNoFirstLevelTestSuite) TestTagDuplicate() {
	s.Assert().ElementsMatch(
		ParseTag("[Test]something[Download]/[Test]Some weird name [Download]"),
		[]string{"Test", "Download"})
}

func (s *ParseTagNoFirstLevelTestSuite) TestEmptyTag() {
	s.Assert().ElementsMatch(
		ParseTag("[]Some weird name"),
		[]string{})
}

func (s *ParseTagNoFirstLevelTestSuite) TestEmptyAndOtherTag() {
	s.Assert().ElementsMatch(
		ParseTag("[]Some weird name[Test]"),
		[]string{"Test"})
}

type ParseTagFirstLevelTestSuite struct {
	suite.Suite
}

func TestParseTagFirstLevelTestSuite(t *testing.T) {
	suite.Run(t, new(ParseTagFirstLevelTestSuite))
}

func (suite *ParseTagFirstLevelTestSuite) SetupSuite() {
	configuration.Init(configuration.Config{
		FirstLevelDirAsTag: true,
	})
}

func (s *ParseTagFirstLevelTestSuite) TestNoFirstLevel() {
	s.Assert().ElementsMatch(ParseTag("[Test]Some weird name"), []string{"Test"})
}

func (s *ParseTagFirstLevelTestSuite) TestWithOneLevel() {
	s.Assert().ElementsMatch(
		ParseTag(strings.Join([]string{"dir1", "[Test]Some weird name"}, string(os.PathSeparator))),
		[]string{"dir1", "Test"})
}

func (s *ParseTagFirstLevelTestSuite) TestWithTwoLevel() {
	s.Assert().ElementsMatch(
		ParseTag(strings.Join([]string{"dir1", "dir2", "[Test]Some weird name"}, string(os.PathSeparator))),
		[]string{"dir1", "Test"})
}
