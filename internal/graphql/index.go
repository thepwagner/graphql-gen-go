package graphql

import (
	"github.com/graphql-go/graphql/language/ast"
)

type IndexedSchema struct {
	TypeIndex map[string]Type
	Query     Type
}

func IndexSchema(schema *ast.Document) IndexedSchema {
	// Index types of the schema, and track the type of the root query:
	index := map[string]Type{}
	var rootQueryName string
	for _, def := range schema.Definitions {
		switch t := def.(type) {
		case *ast.InterfaceDefinition:
			index[t.Name.Value] = indexInterface(t)
		case *ast.ObjectDefinition:
			index[t.Name.Value] = indexObject(t)
		case *ast.SchemaDefinition:
			for _, op := range t.OperationTypes {
				if op.Operation == "query" {
					rootQueryName = op.Type.Name.Value
				}
				// TODO: mutation?
			}
		}
	}

	// Make a second pass to xref indexed types
	for _, indexed := range index {
		for _, f := range indexed.FieldIndex {
			t := index[f.FieldType.Go]
			f.Type = &t
		}
	}
	return IndexedSchema{
		Query:     index[rootQueryName],
		TypeIndex: index,
	}
}

func indexInterface(iface *ast.InterfaceDefinition) Type {
	t := Type{
		Name:       iface.Name.Value,
		Interface:  true,
		FieldIndex: map[string]*Field{},
	}
	for _, f := range iface.Fields {
		field := NewField(f)
		t.Fields = append(t.Fields, field)
		t.FieldIndex[f.Name.Value] = field
	}
	return t
}

func indexObject(obj *ast.ObjectDefinition) Type {
	t := Type{
		Name:       obj.Name.Value,
		FieldIndex: map[string]*Field{},
	}
	for _, f := range obj.Fields {
		field := NewField(f)
		t.Fields = append(t.Fields, field)
		t.FieldIndex[f.Name.Value] = field
	}
	return t
}
