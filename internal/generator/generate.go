package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	log "github.com/sirupsen/logrus"

	"github.com/thepwagner/graphql-gen-go/internal/graphql"
)

func Generate(outputDir string, schema *ast.Document, queries ...string) error {
	index := graphql.IndexSchema(schema)

	// Generate model types, indexed by name:
	types := map[string]*goTypeTemplateParams{}
	for _, gqlType := range index.TypeIndex {
		types[gqlType.Name] = newModelType(gqlType)
	}

	// Generate queries, updating types as necessary
	queryTypes, err := generateQueries(types, index.Query, queries)
	if err != nil {
		return err
	}

	// Write all generated files out:
	for name, params := range types {
		if err := writeGoFile(outputDir, name, "", params); err != nil {
			return err
		}
	}
	for name, params := range queryTypes {
		readFunc := fmt.Sprintf(`import ("io"
"encoding/json"
)

func Read%s(r io.Reader) (%s, error) {
	var ret %s
	if err := json.NewDecoder(r).Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, nil
}`, name, name, strings.ToLower(name[0:1])+name[1:])
		if err := writeGoFile(outputDir, name, readFunc, params...); err != nil {
			return err
		}
	}
	log.WithFields(log.Fields{
		"types":   len(types),
		"queries": len(queries),
	}).Debug("Generated go code.")

	return nil
}

// newModelType creates a Go type for a GraphQL object
func newModelType(gqlType graphql.Type) *goTypeTemplateParams {
	var fields []*goField
	var funcs []*goFunc
	var maxFieldNameLen int
	for _, gqlField := range gqlType.Fields {
		fieldName := gqlField.Name
		gqlField := gqlType.FieldIndex[fieldName]

		if nameLen := len(fieldName); nameLen > maxFieldNameLen {
			maxFieldNameLen = nameLen
		}

		fieldName = strings.Title(fieldName)
		if !gqlType.Interface {
			funcs = append(funcs, &goFunc{
				Signature: fmt.Sprintf("Get%s() %s", fieldName, gqlField.FieldType.GoType()),
				Body:      fmt.Sprintf(`	return r.%s`, fieldName),
			})
		} else {
			fieldName += "()"
		}

		fields = append(fields, &goField{
			Name: fieldName,
			Type: gqlField.FieldType.GoType(),
		})
	}

	// Pad field names to the longest name's length, for alignment:
	for _, field := range fields {
		if pad := maxFieldNameLen - len(field.Name); pad > 0 {
			field.Name = field.Name + strings.Repeat(" ", pad)
		}
	}

	return &goTypeTemplateParams{
		Name:      gqlType.Name,
		Interface: gqlType.Interface,
		Fields:    fields,
		Funcs:     funcs,
	}
}

func writeGoFile(dir, name string, code string, params ...*goTypeTemplateParams) error {
	// Write code to buffer:
	var buf bytes.Buffer
	_, pkg := filepath.Split(dir)
	if err := writeHeader(&buf, pkg); err != nil {
		return err
	}
	buf.WriteString(code)
	for _, param := range params {
		if err := goTypeTemplate.Execute(&buf, param); err != nil {
			return err
		}
	}

	// Format, so that templates can be readable:
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	path := filepath.Join(dir, fmt.Sprintf("%s.gql.go", strings.ToLower(name)))
	return ioutil.WriteFile(path, formatted, 0660)
}
