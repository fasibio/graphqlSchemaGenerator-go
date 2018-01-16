# Project name

Graphql Schema Generator


# Description
Generate graphql-go code from a Schemafile.
```gql
type test {
  id: Int!
  name: [String]
}

```

generate

```gql
package schema

import graphql "github.com/graphql-go/graphql"

type Test struct {
        Id   int      `json:"id"`
        Name []string `json:"name"`
}

func GetTest() *graphql.Object {
        return graphql.NewObject(graphql.ObjectConfig{Name: "test", Fields: graphql.Fields{"id": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)}, "name": &graphql.Field{Type: graphql.NewList(graphql.String)}}})
}

```
# Usage
## Generate the Golang Code
### Simple
 - Use the [webpage](http://gql2go.fasibio.de)
### Complexer
 - Use the [graphql Api](http://gql2go.fasibio.de/graphql) to get the Go Code String at API Result
 ```gql
 {
  generateGoCode(schemaStr: "your schema file at string the \n(breaks are importent) ")
}
 ```
### Complex
```go get github.com/fasibio/graphqlSchemaGenerator-go```

```golang
import (
	"github.com/fasibio/graphqlSchemaGenerator-go/goCodeGenerator"
	"github.com/fasibio/graphqlSchemaGenerator-go/schemaInterpretations"
)
  //return the Golang Code (String)
  goCode := goCodeGenerator.GetGenerateFile(schemaInterpretations.GetSchemaList(schemaStr),   schemaInterpretations.GetEnumList(schemaStr))
```

## Use the generated Code in your Go-Application
See [https://github.com/graphql-go/graphql](https://github.com/graphql-go/graphql)
to learn how to Use graphql in Go
- Generate a Folder/Package schema in you Project
- Copy the generated Go File to a folder schema in a for example schema.go
- Create Graphql Queries and set as Type the generated Graphql Schema you want to use. As result you can fill the Generated struct
```golang

	fields := graphql.Fields{
		"test": &graphql.Field{
			Type: schema.GetTest(),// <- this is the generated Funktion
			Args: graphql.FieldConfigArgument{
				"type": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "developer",
				},
			},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var result schema.Test // <- this is the generated Struct
					result = schema.Test{
						Id:   "MUHAHAH",
						Name: "Fasibio",
					}

				return result, nil
			},
		},
```

# By the way 
  - The generator is in Work 
  - Im Happy if you want to help to make it better
  - So feel free to open issu or send Pull Requests

# TODO
- [x] understand schemas
- [x] generate Code for Schemas
  - [x] understand Deprecated
  - [] generate Code for Deprecated 
- [x] understand Enums
- [] generate Code for Enums

# License 
GNU General Public License v3.0
