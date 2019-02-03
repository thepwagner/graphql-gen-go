package graphql

import (
	"strings"

	"github.com/graphql-go/graphql/language/ast"
)

// TODO: this is too simple - i could have a list of lists of lists
type FieldType struct {
	List    bool
	NonNull bool
	Go      string
	Scalar  bool
}

func NewFieldType(def *ast.FieldDefinition) FieldType {
	var f FieldType
	for t := def.Type; t != nil; {
		switch o := t.(type) {
		case *ast.List:
			f.List = true
			t = o.Type
		case *ast.NonNull:
			f.NonNull = true
			t = o.Type
		case *ast.Named:
			switch o.Name.Value {
			case "String", "ID":
				f.Go = "string"
				f.Scalar = true
			case "Int":
				f.Go = "int64"
				f.Scalar = true
			case "Float":
				f.Go = "float64"
				f.Scalar = true
			default:
				f.Go = o.Name.Value
			}
			return f

		default:
			return f
		}
	}
	return f
}

func (f FieldType) GoType() string {
	// TODO: using pointers for nullable fields makes the API annoying
	var s strings.Builder
	if f.List {
		s.WriteString("[]")
	}
	if !f.NonNull {
		s.WriteString("*")
	}
	s.WriteString(f.Go)
	return s.String()
}
