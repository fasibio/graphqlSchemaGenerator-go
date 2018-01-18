# Project name
Graphql Schema Generator


# Description
Generate Go code from a GraphQL schema.
For example:
```gql
type test {
  id: Int!
  name: [String]
}

```

will generate:

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
## Generate Go Code
### Simple
 - Use the [webpage](http://gql2go.fasibio.de)
### Advanced
 - Use the [graphql Api](http://gql2go.fasibio.de/graphql) to get Go code from a schema string
 ```gql
 {
  generateGoCode(schemaStr: "your schema file at string the \n(breaks are importent) ")
}
 ```
### Expert
```go get github.com/fasibio/graphqlSchemaGenerator-go```

```golang
import (
	"github.com/fasibio/graphqlSchemaGenerator-go/goCodeGenerator"
	"github.com/fasibio/graphqlSchemaGenerator-go/schemaInterpretations"
)
  //return the Golang Code (String)
  goCode := goCodeGenerator.GetGenerateFile(schemaInterpretations.GetSchemaList(schemaStr),   schemaInterpretations.GetEnumList(schemaStr))
```

## Use the generated code in your Go application
See [https://github.com/graphql-go/graphql](https://github.com/graphql-go/graphql)
to learn how to use GraphQL in Go
- Generate a folder/package schema in your project
- Copy the generated Go file to a folder schema (for example a schema.go)
- Create GraphQL queries and set the type as the generated Graphql Schema you want to use. As result you can fill the generated struct
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

# Notes
  - The generator is a work in progress
  - I'm happy if you want to help to make it better
  - So feel free to open an issue or send pull requests

# TODO
- [x] understand schemas
- [x] generate code for schemas
  - [x] understand deprecated
  - [] generate code for deprecated 
- [x] understand enums
- [] generate code for enums

# License 
GNU General Public License v3.0
