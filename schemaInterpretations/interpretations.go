package schemaInterpretations

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fasibio/graphqlSchemaGenerator-go/helper"
)

func splitByKeyword(value string, keyWord string) (string, string) {
	index := strings.Index(value, keyWord)
	secend := helper.TrimEmpty(string(value[index+len(keyWord):]))
	first := string(value[0:index])
	return first, secend
}

func getFields(fields []string) []Field {
	var result []Field
	for _, v := range fields {
		v = removeCommentsFromLine(v)
		indexOfSeperator := strings.Index(v, ":")
		if indexOfSeperator != -1 {
			vStrArray := []rune(v)

			key := string(vStrArray[0:indexOfSeperator])
			value := string(vStrArray[indexOfSeperator+1:])
			value = strings.Trim(value, " ")
			deprecated := false
			required := false
			isArray := false
			if helper.MatchString("@deprecated", value) {
				value, _ = splitByKeyword(value, "@deprecated")
				deprecated = true
			}
			value = helper.TrimEmpty(value)
			if helper.MatchString(`^\[.*\]$`, value) {
				isArray = true
				value = string(value[1 : len(value)-1])
			}
			key = helper.TrimEmpty(key)

			if helper.MatchString("!$", value) {
				required = true
				value, _ = splitByKeyword(value, "!")
			}
			result = append(result, Field{
				DataType:     value,
				Name:         key,
				IsDeprecated: deprecated,
				IsArray:      isArray,
				Required:     required,
			})
		}

	}
	return result
}

func GetEnumList(schemaStr string) []Enum {
	var enumList []Enum

	var reEnum = regexp.MustCompile(`(?ms)enum.*?\}`)
	enums := reEnum.FindAllString(schemaStr, -1)
	for index := range enums {
		enumList = append(enumList, getEnumObj(enums[index]))
	}
	return enumList
}

func removeCommentsFromLine(line string) string {
	commentBeginIndex := strings.Index(line, "#")
	lineStrArray := []rune(line)
	result := line
	if commentBeginIndex != -1 {
		result = string(lineStrArray[0:commentBeginIndex])
	}
	return result
}

func getQueryParams(paramsStr string) []QueryParam {
	lines := strings.Split(paramsStr, ",")
	var result []QueryParam
	for index := range lines {
		oneParamStr := lines[index]
		oneParamStr = helper.TrimEmpty(oneParamStr)
		paramLines := strings.Split(oneParamStr, ":")
		result = append(result, QueryParam{
			Name: helper.TrimEmpty(paramLines[0]),
			Type: helper.TrimEmpty(paramLines[1]),
		})
	}
	return result
}

func getQuery(line string) Query {
	paramBeginIndex := strings.Index(line, "(")
	paramEndIndex := strings.Index(line, ")")
	line = helper.TrimEmpty(line)
	schemaStrArray := []rune(line)
	name := string(schemaStrArray[0 : paramBeginIndex-2])
	name = strings.Trim(name, "\n")
	name = helper.TrimEmpty(name)
	return Query{
		Name:   name,
		Params: getQueryParams(string(schemaStrArray[paramBeginIndex-1 : paramEndIndex-2])),
	}
}

func GetQuerys(schemaStr string) []Query {
	var re = regexp.MustCompile(`(?ms)type query.*?\}`)
	schemas := re.FindAllString(schemaStr, -1)
	if len(schemas) > 1 {
		panic("find more than one query Type!")
	}
	querySchema := schemas[0]
	lines := strings.Split(querySchema, "\n")
	fmt.Printf("%+v\n", getQuery(lines[1]))
	return nil
}

func GetSchemaList(schemaStr string) []Schema {
	GetQuerys(schemaStr)
	var re = regexp.MustCompile(`(?ms)type.*?\}`)
	schemas := re.FindAllString(schemaStr, -1)

	var schemaList []Schema
	for index := range schemas {
		schemaObj := getSchemaObj(schemas[index])
		if schemaObj.Name != "query" {
			schemaList = append(schemaList, schemaObj)
		}

	}
	return schemaList
}

func getEnumObj(enumStr string) Enum {
	index := strings.Index(enumStr, "{")
	enumStrArray := []rune(enumStr)
	name := string(enumStrArray[5:index])
	name = strings.Trim(name, "\n")
	name = helper.TrimEmpty(name)
	rest := string(enumStrArray[index+1:])
	rest = strings.TrimRight(rest, "}")
	fields := strings.Split(rest, "\n")
	for index, value := range fields {
		value = removeCommentsFromLine(value)
		fields[index] = helper.TrimEmpty(value)
	}
	return Enum{
		Name:   removeCommentsFromLine(name),
		Values: fields,
	}
}

func getSchemaObj(schemaStr string) Schema {
	index := strings.Index(schemaStr, "{")
	schemaStrArray := []rune(schemaStr)
	name := string(schemaStrArray[5:index])
	name = strings.Trim(name, "\n")
	name = helper.TrimEmpty(name)
	var implement string
	if helper.MatchString("implements", name) {
		name, implement = splitByKeyword(name, "implements")
	}
	rest := string(schemaStrArray[index+1:])
	rest = strings.TrimRight(rest, "}")
	fields := strings.Split(rest, "\n")
	return Schema{
		Name:       removeCommentsFromLine(name),
		Fields:     getFields(fields),
		Implements: removeCommentsFromLine(implement),
	}
}
