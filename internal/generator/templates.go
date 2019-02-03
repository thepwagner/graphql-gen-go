package generator

import (
	"io"
	"text/template"
)

var goTypeTemplate = template.Must(template.New("type").Parse(`
type {{.Name}} {{if .Interface}}interface{{else}}struct{{end}} {
{{range $field := .Fields}}{{.Name}} {{.Type}} {{if ne .Tags ""}}` + "`{{.Tags}}`" + `{{end}}
{{end}}}

{{range $func := .Funcs}}
func (r {{$.Name}}) {{.Signature}} {
{{.Body}}
}
{{end}}
`))

type goTypeTemplateParams struct {
	Name      string
	Interface bool
	Fields    []*goField
	Funcs     []*goFunc
}

type goField struct {
	Name string
	Type string
	Tags string
}

type goFunc struct {
	Signature string
	Body      string
}

var headerTmpl = template.Must(template.New("header").Parse(`// Code generated by graphql-gen-go. DO NOT EDIT.
package {{.Name}}
`))

type headerTmplParams struct {
	Name string
}

func writeHeader(w io.Writer, pkg string) error {
	return headerTmpl.Execute(w, headerTmplParams{
		Name: pkg,
	})
}