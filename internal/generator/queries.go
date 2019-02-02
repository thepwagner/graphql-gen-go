package generator

import (
	"fmt"
	"github.com/thepwagner/magenny/internal/graphql"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

func (g *Generator) GenerateQueries(queries ...string) error {
	for _, query := range queries {
		parsed, err := parser.Parse(parser.ParseParams{
			Source: query,
			Options: parser.ParseOptions{
				NoLocation: false,
				NoSource:   true,
			},
		})
		if err != nil {
			return err
		}

		for _, def := range parsed.Definitions {
			switch o := def.(type) {
			case *ast.OperationDefinition:
				switch o.Operation {
				case "query":
					if err := g.generateQuery(o); err != nil {
						return err
					}
				default:
					fmt.Sprintf("unknonwn %T %+v\n", o, o)
				}
			}
		}
	}

	return nil
}

func (g *Generator) generateQuery(op *ast.OperationDefinition) error {
	name := op.Name.Value
	types, err := typeParams(name, op.SelectionSet, &g.indexed.Query)
	if err != nil {
		return err
	}

	f, err := g.NewFile(name)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, t := range types {
		_, _ = fmt.Fprintln(f, "")
		if err := goTypeTemplate.Execute(f, t); err != nil {
			return err
		}
	}
	return nil
}

func typeParams(name string, selSet *ast.SelectionSet, typePtr *graphql.Type) ([]*goTypeTemplateParams, error) {
	params := &goTypeTemplateParams{
		Interface: true,
		TypeName:  name,
	}
	types := []*goTypeTemplateParams{params}

	for _, sel := range selSet.Selections {
		switch t := sel.(type) {
		case *ast.Field:
			selName := t.Name.Value
			indexedField, found := typePtr.FieldIndex[selName]
			if !found {
				return nil, fmt.Errorf("field %q not found", selName)
			}
			selNameTitle := strings.Title(selName)

			// If this is a complex type, generate a type with it's specific fields:
			var goType string
			if indexedField.FieldType.Primitive {
				goType = indexedField.FieldType.GoType()
			} else {

				goType = fmt.Sprintf("%s%s", name, selNameTitle)
			}

			if t.SelectionSet != nil {
				descendents, err := typeParams(goType, t.SelectionSet, indexedField.Type)
				if err != nil {
					return nil, err
				}
				types = append(types, descendents...)
			}

			params.Fields = append(params.Fields, &goField{
				Name: fmt.Sprintf("%s()", selNameTitle),
				Type: goType,
			})
		}
	}

	return types, nil
}
