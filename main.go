package main

import (
	"fmt"
	"io/ioutil"
	"log"

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

	goCode := goCodeGenerator.GetGenerateFile(schemaInterpretations.GetSchemaList(schemaStr), schemaInterpretations.GetEnumList(schemaStr))
	fmt.Printf(goCode)
}

