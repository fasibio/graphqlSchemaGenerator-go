package schemaInterpretations

import (
	"regexp"
	"strings"

	"fasibio.de/graphqlSchemaGenerator-go/helper"
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

func GetSchemaList(schemaStr string) []Schema {
	var re = regexp.MustCompile(`(?ms)type.*?\}`)
	schemas := re.FindAllString(schemaStr, -1)

	var schemaList []Schema
	for index := range schemas {
		schemaList = append(schemaList, getSchemaObj(schemas[index]))
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
		fields[index] = helper.TrimEmpty(value)
	}
	return Enum{
		Name:   name,
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
		Name:       name,
		Fields:     getFields(fields),
		Implements: implement,
	}
}
