package graphql

import "github.com/graphql-go/graphql/language/ast"

type Type struct {
	Name       string
	Interface  bool
	Fields     []*Field
	FieldIndex map[string]*Field
}

type Field struct {
	Name      string
	FieldType FieldType
	Type      *Type

	Definition *ast.FieldDefinition
}

func NewField(f *ast.FieldDefinition) *Field {
	return &Field{
		Name:       f.Name.Value,
		Definition: f,
		FieldType:  NewFieldType(f),
	}
}
