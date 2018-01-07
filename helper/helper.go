package helper

import (
	"regexp"
	"strings"
)

func TrimEmpty(value string) string {
	result := strings.Trim(value, " ")
	result = strings.Trim(result, "\t")
	return result
}

func MatchString(rexexp string, value string) bool {
	var re = regexp.MustCompile(rexexp)
	return re.MatchString(value)
}
