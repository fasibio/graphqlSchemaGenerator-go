package goCodeGenerator

import (
	"fmt"
	"strings"

	"fasibio.de/graphqlSchemaGenerator-go/helper"
	"fasibio.de/graphqlSchemaGenerator-go/schemaInterpretations"
	. "github.com/dave/jennifer/jen"
)

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
func GenerateFile(schemaList []schemaInterpretations.Schema, enumList []schemaInterpretations.Enum) {
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
