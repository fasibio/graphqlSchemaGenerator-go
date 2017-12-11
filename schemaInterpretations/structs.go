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
