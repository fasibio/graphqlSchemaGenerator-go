# GOlang Code Generator for GRAPHQL Schemafiles
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

So you can use Schemafiles ...
##  How to Use


Check out the Repo: 
Write your schema to the schema.schema file an run the main.go

TODO:
- [x] understand schemas
- [x] generate Code for Schemas
  - [x] understand Deprecated
  - [] generate Code for Deprecated 
- [x] understand Enums
- [] generate Code for Enums

