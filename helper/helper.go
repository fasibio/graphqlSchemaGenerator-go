package helper

import (
	"bytes"
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

func MakeFirstLowerCase(s string) string {

	if len(s) < 2 {
		return strings.ToLower(s)
	}

	bts := []byte(s)

	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}
