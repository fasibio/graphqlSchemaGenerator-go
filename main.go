package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"fasibio.de/graphqlSchemaGenerator-go/goCodeGenerator"
	"fasibio.de/graphqlSchemaGenerator-go/schemaInterpretations"
)

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
	var enumList []schemaInterpretations.Enum
	for index := range schemas {
		schemaList = append(schemaList, schemaInterpretations.GetSchemaObj(schemas[index]))
	}

	var reEnum = regexp.MustCompile(`(?ms)enum.*?\}`)
	enums := reEnum.FindAllString(schemaStr, -1)
	for index := range enums {
		enumList = append(enumList, schemaInterpretations.GetEnumObj(enums[index]))
	}
	//fmt.Printf("%+v\n", enumList)
	goCodeGenerator.GenerateFile(schemaList, enumList)

}
