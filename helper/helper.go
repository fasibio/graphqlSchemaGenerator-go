package helper

import (
	"regexp"
	"strings"
)

func TrimEmpty(value string) string {
	return strings.Trim(value, " ")
}

func MatchString(rexexp string, value string) bool {
	var re = regexp.MustCompile(rexexp)
	return re.MatchString(value)
}
