package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"fasibio.de/graphqlSchemaGenerator-go/SchemaFileType"
	"fasibio.de/graphqlSchemaGenerator-go/helper"
	"fasibio.de/graphqlSchemaGenerator-go/schemaInterpretations"
	. "github.com/dave/jennifer/jen"
)

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

	var schemaList []schemaInterpretations.Schema
	var enumList []enum
	for index := range schemas {
		schemaList = append(schemaList, schemaInterpretations.GetSchemaObj(schemas[index]))
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

func getGraphQlType(fieldValue schemaInterpretations.Field) *Statement {

	var result *Statement
	if isSimpleDataType(fieldValue.DataType) {
		result = Qual("github.com/graphql-go/graphql", fieldValue.DataType)
	} else {
		result = Id("Get" + strings.Title(fieldValue.DataType) + "()")
	}
	if fieldValue.IsArray {
		result = Qual("github.com/graphql-go/graphql", "NewList").Call(result)
	}

	if fieldValue.Required {
		result = Qual("github.com/graphql-go/graphql", "NewNonNull").Call(result)
	}
	return result

}
func generateFile(schemaList []schemaInterpretations.Schema, enumList []enum) {
	f := NewFile("schema")
	for _, value := range schemaList {
		var typeStructValues []Code
		var graphQLFields []Code
		for _, fieldValue := range value.Fields {
			val := Id(strings.Title(fieldValue.Name))
			if fieldValue.IsArray {
				val = val.Id("[]")
			}
			valStr, _ := getGoDataType(fieldValue.DataType)
			val = val.Id(valStr).Tag(map[string]string{"json": fieldValue.Name})
			typeStructValues = append(typeStructValues, val)

			fieldval := Lit(fieldValue.Name).Id(":").Op("&").Qual("github.com/graphql-go/graphql", "Field").Values(
				Id("Type:").Add(getGraphQlType(fieldValue)),
			)
			graphQLFields = append(graphQLFields, fieldval)
		}

		f.Type().Id(strings.Title(value.Name)).Struct(
			typeStructValues...,
		)

		f.Func().Id("Get" + strings.Title(value.Name)).Params().Id("*graphql.Object").Block(
			Return(Qual("github.com/graphql-go/graphql", "NewObject").Call(
				Qual("github.com/graphql-go/graphql", "ObjectConfig").Values(
					Id("Name:").Lit(helper.TrimEmpty(value.Name)),
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

func getEnumObj(enumStr string) enum {
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
	return enum{
		name:   name,
		values: fields,
	}
}
