package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fasibio/graphqlSchemaGenerator-go/goCodeGenerator"
	"github.com/fasibio/graphqlSchemaGenerator-go/schemaInterpretations"
)

func main() {

	fmt.Println("Lets go")
	buf, err := ioutil.ReadFile("/home/fasibio/git/goWorkspace/src/github.com/fasibio/graphqlSchemaGenerator-go/schema.schema")
	if err != nil {
		log.Fatalln(err)
	}
	schemaStr := string(buf)

	goCodeGenerator.GetGenerateFile(schemaInterpretations.GetSchemaList(schemaStr), schemaInterpretations.GetEnumList(schemaStr))
	// fmt.Printf(goCode)
}
