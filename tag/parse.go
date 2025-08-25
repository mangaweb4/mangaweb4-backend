package tag

import (
	"os"
	"regexp"
	"strings"

	"github.com/mangaweb4/mangaweb4-backend/configuration"
)

var regex *regexp.Regexp

func init() {
	const pattern = "\\[(.*?)\\]"
	var err error
	regex, err = regexp.Compile(pattern)

	if err != nil {
		panic(err)
	}
}

func parseFirstLevelDirTag(name string) string {
	seprator := strings.Index(name, string(os.PathSeparator))

	if seprator == -1 {
		return ""
	}

	return name[0:seprator]
}

func ParseTag(name string) []string {
	c := configuration.Get()

	matches := regex.FindAllStringSubmatch(name, -1)
	tagSet := make(map[string]bool)
	output := make([]string, 0)

	if c.FirstLevelDirAsTag {
		if tag := parseFirstLevelDirTag(name); tag != "" {
			output = append(output, tag)
		}
	}

	for _, match := range matches {
		tag := match[1]
		if tag == "" {
			continue
		}

		if _, found := tagSet[tag]; !found {
			tagSet[tag] = true
			output = append(output, tag)
		}
	}

	return output
}
