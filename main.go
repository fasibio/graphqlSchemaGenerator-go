package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"fasibio.de/graphqlSchemaGenerator-go/SchemaFileType"
	. "github.com/dave/jennifer/jen"
)

type schema struct {
	name       string
	fields     []field
	implements string
}

type field struct {
	name         string
	dataType     string
	isDeprecated bool
	required     bool
	isArray      bool
}

type enum struct {
	name   string
	values []string
}

func main() {

	fmt.Println("Lets go")
	buf, err := ioutil.ReadFile("/home/fasibio/git/goWorkspace/src/fasibio.de/graphqlSchemaGenerator-go/schema.schema")
	if err != nil {
		log.Fatalln(err)
	}
	schemaStr := string(buf)

	var re = regexp.MustCompile(`(?ms)type.*?\}`)
	schemas := re.FindAllString(schemaStr, -1)

	var schemaList []schema
	var enumList []enum
	for index := range schemas {
		schemaList = append(schemaList, getSchemaObj(schemas[index]))
	}

	var reEnum = regexp.MustCompile(`(?ms)enum.*?\}`)
	enums := reEnum.FindAllString(schemaStr, -1)
	for index := range enums {
		enumList = append(enumList, getEnumObj(enums[index]))
	}
	//fmt.Printf("%+v\n", enumList)
	generateFile(schemaList, enumList)

}

func isSimpleDataType(typeStr string) bool {
	isSimpleDataType := false
	switch typeStr {
	case "String":
		isSimpleDataType = true
	case "Float":
		isSimpleDataType = true
	case "Int":
		isSimpleDataType = true
	case "Boolean":
		isSimpleDataType = true
	case "ID":
		isSimpleDataType = true
	}
	return isSimpleDataType
}

func getGoDataType(typeStr string) (string, bool) {
	switch typeStr {
	case "String":
		return "string", true
	case "Float":
		return "float64", true
	case "Int":
		return "int", true
	case "Boolean":
		return "bool", true
	}
	return strings.Title(typeStr), false
}

func getGraphQlType(fieldValue field) *Statement {

	var result *Statement
	if isSimpleDataType(fieldValue.dataType) {
		result = Qual("github.com/graphql-go/graphql", fieldValue.dataType)
	} else {
		result = Id("Get" + strings.Title(fieldValue.dataType) + "()")
	}
	if fieldValue.isArray {
		result = Qual("github.com/graphql-go/graphql", "NewList").Call(result)
	}

	if fieldValue.required {
		result = Qual("github.com/graphql-go/graphql", "NewNonNull").Call(result)
	}
	return result

}
func generateFile(schemaList []schema, enumList []enum) {
	f := NewFile("schema")
	for _, value := range schemaList {
		var typeStructValues []Code
		var graphQLFields []Code
		for _, fieldValue := range value.fields {
			val := Id(strings.Title(fieldValue.name))
			if fieldValue.isArray {
				val = val.Id("[]")
			}
			valStr, _ := getGoDataType(fieldValue.dataType)
			val = val.Id(valStr).Tag(map[string]string{"json": fieldValue.name})
			typeStructValues = append(typeStructValues, val)

			fieldval := Lit(fieldValue.name).Id(":").Op("&").Qual("github.com/graphql-go/graphql", "Field").Values(
				Id("Type:").Add(getGraphQlType(fieldValue)),
			)
			graphQLFields = append(graphQLFields, fieldval)
		}

		f.Type().Id(strings.Title(value.name)).Struct(
			typeStructValues...,
		)

		f.Func().Id("Get" + strings.Title(value.name)).Params().Id("*graphql.Object").Block(
			Return(Qual("github.com/graphql-go/graphql", "NewObject").Call(
				Qual("github.com/graphql-go/graphql", "ObjectConfig").Values(
					Id("Name:").Lit(TrimEmpty(value.name)),
					Id("Fields:").Qual("github.com/graphql-go/graphql", "Fields").Values(graphQLFields...),
				),
			)),
		)

	}
	fmt.Printf("%#v", f)
}

func getType(schema string) string {
	if strings.HasPrefix(schema, "type") {
		return SchemaFileType.TYPE
	}
	return "UNKNOWN"
}

func TrimEmpty(value string) string {
	return strings.Trim(value, " ")
}

func matchString(rexexp string, value string) bool {
	var re = regexp.MustCompile(rexexp)
	return re.MatchString(value)
}

func getFields(fields []string) []field {
	var result []field
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
			if matchString("@deprecated", value) {
				value, _ = splitByKeyword(value, "@deprecated")
				deprecated = true
			}
			value = TrimEmpty(value)
			if matchString(`^\[.*\]$`, value) {
				isArray = true
				value = string(value[1 : len(value)-1])
			}
			key = TrimEmpty(key)
			if matchString("!$", value) {
				required = true
				value, _ = splitByKeyword(value, "!")
			}
			result = append(result, field{
				dataType:     value,
				name:         key,
				isDeprecated: deprecated,
				isArray:      isArray,
				required:     required,
			})
		}

	}
	return result
}

func splitByKeyword(value string, keyWord string) (string, string) {
	index := strings.Index(value, keyWord)
	secend := TrimEmpty(string(value[index+len(keyWord):]))
	first := string(value[0:index])
	return first, secend
}

func getEnumObj(enumStr string) enum {
	index := strings.Index(enumStr, "{")
	enumStrArray := []rune(enumStr)
	name := string(enumStrArray[5:index])
	name = strings.Trim(name, "\n")
	name = TrimEmpty(name)
	rest := string(enumStrArray[index+1:])
	rest = strings.TrimRight(rest, "}")
	fields := strings.Split(rest, "\n")
	for index, value := range fields {
		fields[index] = TrimEmpty(value)
	}
	return enum{
		name:   name,
		values: fields,
	}
}

func getSchemaObj(schemaStr string) schema {
	index := strings.Index(schemaStr, "{")
	schemaStrArray := []rune(schemaStr)
	name := string(schemaStrArray[5:index])
	name = strings.Trim(name, "\n")
	name = TrimEmpty(name)
	var implement string
	if matchString("implements", name) {
		name, implement = splitByKeyword(name, "implements")
	}
	rest := string(schemaStrArray[index+1:])
	rest = strings.TrimRight(rest, "}")
	fields := strings.Split(rest, "\n")
	return schema{
		name:       name,
		fields:     getFields(fields),
		implements: implement,
	}
}
