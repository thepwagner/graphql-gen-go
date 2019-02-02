package generator

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/thepwagner/magenny/internal/graphql"
)

func (g Generator) GenerateTypes() error {
	for _, gqlType := range g.indexed.TypeIndex {
		if err := g.writeGoType(gqlType); err != nil {
			return err
		}
	}
	return nil
}

var goTypeTemplate = template.Must(template.New("type").Parse(`type {{.TypeName}} {{ if .Interface }}interface{{ else }}struct{{end}} {
{{ range $field := .Fields }}	{{ $field.Name }} {{ $field.Type }}
{{ end }}}
`))

type goTypeTemplateParams struct {
	TypeName  string
	Interface bool
	Fields    []*goField
}

type goField struct {
	Name string
	Type string
}

func (g Generator) writeGoType(gqlType graphql.Type) error {
	f, err := g.NewFile(gqlType.Name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = fmt.Fprintln(f, "")

	params := goTypeTemplateParams{
		TypeName:  gqlType.Name,
		Interface: gqlType.Interface,
		Fields:    goFields(gqlType),
	}
	return goTypeTemplate.Execute(f, params)
}

func goFields(gqlType graphql.Type) []*goField {
	var fields []*goField
	var maxFieldNameLen int
	for _, gqlField := range gqlType.Fields {
		fieldName := gqlField.Name
		gqlField := gqlType.FieldIndex[fieldName]
		if gqlType.Interface {
			fieldName += "()"
		}
		if nameLen := len(fieldName); nameLen > maxFieldNameLen {
			maxFieldNameLen = nameLen
		}

		fields = append(fields, &goField{
			Name: strings.Title(fieldName),
			Type: gqlField.FieldType.GoType(),
		})
	}

	// Pad field names to the longest name's length, for alignment:
	for _, field := range fields {
		if pad := maxFieldNameLen - len(field.Name); pad > 0 {
			field.Name = field.Name + strings.Repeat(" ", pad)
		}
	}
	return fields
}
