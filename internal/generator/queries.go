package generator

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"

	"github.com/thepwagner/graphql-gen-go/internal/graphql"
)

func generateQueries(types map[string]*goTypeTemplateParams, rootQuery graphql.Type, queries []string) (map[string][]*goTypeTemplateParams, error) {
	queryTypes := map[string][]*goTypeTemplateParams{}
	for _, query := range queries {
		parsed, err := parser.Parse(parser.ParseParams{
			Source: query,
			Options: parser.ParseOptions{
				NoLocation: false,
				NoSource:   true,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, def := range parsed.Definitions {
			switch o := def.(type) {
			case *ast.OperationDefinition:
				switch o.Operation {
				case "query":
					opTypes, err := generateQuery(types, &rootQuery, o)
					if err != nil {
						return nil, err
					}
					queryTypes[o.Name.Value] = opTypes
				default:
					return nil, fmt.Errorf("unknown operation: %s", def)
				}
			}
		}
	}
	return queryTypes, nil
}

func generateQuery(types map[string]*goTypeTemplateParams, rootQuery *graphql.Type, op *ast.OperationDefinition) ([]*goTypeTemplateParams, error) {
	name := op.Name.Value
	views, err := viewInterfaces(types, "", name, op.SelectionSet, rootQuery)
	if err != nil {
		return nil, err
	}
	impl, err := implStruct(name, op.SelectionSet, rootQuery)
	if err != nil {
		return nil, err
	}

	return append(views, impl), nil
}

func viewInterfaces(types map[string]*goTypeTemplateParams, prefix, name string, selSet *ast.SelectionSet, typePtr *graphql.Type) ([]*goTypeTemplateParams, error) {
	params := &goTypeTemplateParams{
		Interface: true,
		Name:      prefix + name,
	}
	ret := []*goTypeTemplateParams{params}

	for _, sel := range selSet.Selections {
		switch t := sel.(type) {
		case *ast.Field:
			selName := t.Name.Value
			indexedField, found := typePtr.FieldIndex[selName]
			if !found {
				return nil, fmt.Errorf("field %q not found", selName)
			}
			selNameTitle := strings.Title(selName)

			var methodName string
			var goType string
			if indexedField.FieldType.Scalar {
				// Scalars are returned directly:
				methodName = fmt.Sprintf("Get%s()", selNameTitle)
				goType = indexedField.FieldType.GoType()
			} else {
				// Object types have a custom view generated to match the query's scope:
				goType = fmt.Sprintf("%s%s%s", prefix, name, selNameTitle)
				methodName = fmt.Sprintf("Get%s%s%s()", prefix, name, selNameTitle)

				if t.SelectionSet != nil {
					descendents, err := viewInterfaces(types, prefix+name, selNameTitle, t.SelectionSet, indexedField.Type)
					if err != nil {
						return nil, err
					}
					ret = append(ret, descendents...)
				}

				if indexedField.FieldType.List {
					goType = fmt.Sprintf("[]%s", goType)

					// Add a function to the model type to satisfy the invariant list
					hostType := types[name]
					hostType.Funcs = append(hostType.Funcs, &goFunc{
						Signature: fmt.Sprintf("%s %s", methodName, goType),
						Body: fmt.Sprintf(`ret := make(%s, len(r.%s))
	for i, o := range r.%s {
		ret[i] = o
	} 
	return ret`, goType, selNameTitle, selNameTitle),
					})
				}
			}

			params.Fields = append(params.Fields, &goField{
				Name: methodName,
				Type: goType,
			})
		}
	}

	return ret, nil
}

func implStruct(name string, selSet *ast.SelectionSet, typePtr *graphql.Type) (*goTypeTemplateParams, error) {
	structName := strings.ToLower(name[0:1]) + name[1:]
	params := &goTypeTemplateParams{
		Name: structName,
	}

	for _, sel := range selSet.Selections {
		switch t := sel.(type) {
		case *ast.Field:
			selName := t.Name.Value
			selNameTitle := strings.Title(selName)
			indexedField, _ := typePtr.FieldIndex[selName]
			params.Fields = append(params.Fields, &goField{
				Name: selNameTitle,
				Type: indexedField.FieldType.GoType(),
				Tags: fmt.Sprintf(`json:"%s"`, selName),
			})

			// If this is a complex type, generate a type with it's specific fields:
			var goType string
			if indexedField.FieldType.Scalar {
				goType = indexedField.FieldType.GoType()
			} else {
				goType = fmt.Sprintf("%s%s", name, selNameTitle)
			}

			params.Funcs = append(params.Funcs, &goFunc{
				Signature: fmt.Sprintf("Get%s%s() %s", name, selNameTitle, goType),
				Body:      fmt.Sprintf(`	return r.%s`, selNameTitle),
			})
		}
	}

	return params, nil
}
