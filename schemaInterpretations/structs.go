package schemaInterpretations

type Schema struct {
	Name       string
	Fields     []Field
	Implements string
}

type Field struct {
	Name         string
	DataType     string
	IsDeprecated bool
	Required     bool
	IsArray      bool
}

type Enum struct {
	Name   string
	Values []string
}

type Query struct {
	Name       string
	Params     []QueryParam
	RetrunType string
}

type QueryParam struct {
	Name string
	Type string
}
