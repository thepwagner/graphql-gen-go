package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
)

var goPackage string
var gqlSchemaPath string

var generateCmd = &cobra.Command{
	Use: "generate",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called", gqlSchemaPath, goPackage)

		// TODO: what about dynamic queries?
		// i.e. i want to look up N nodes using aliases returned in a map[string]Node

		// Map golang package to directory:
		pkg, err := build.Import(goPackage, build.Default.SrcDirs()[0], build.IgnoreVendor)
		if err != nil {
			panic(err)
		}

		// Parse files in directory:
		fset := token.NewFileSet()
		parsedPackages, err := parser.ParseDir(fset, pkg.Dir, nil, parser.AllErrors)
		if err != nil {
			panic(err)
		}
		// Crawl top level constants:
		for _, pkg := range parsedPackages {
			for _, file := range pkg.Files {
				for _, decl := range file.Decls {
					if genDecl, ok := decl.(*ast.GenDecl); ok {
						if genDecl.Tok != token.CONST {
							continue
						}
						for _, spec := range genDecl.Specs {
							if value, ok := spec.(*ast.ValueSpec); ok {
								if lit, ok := value.Values[0].(*ast.BasicLit); ok {
									if lit.Kind != token.STRING {
										continue
									}

									// We found a top level string constant, is it a GraphQL query?
									fmt.Println(lit.Value)
								}
							}
						}
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVar(&goPackage, "pkg", "", "Go package to scan")
	_ = generateCmd.MarkFlagRequired("pkg")
	_ = viper.BindPFlag("pkg", generateCmd.Flags().Lookup("pkg"))

	generateCmd.Flags().StringVar(&gqlSchemaPath, "schema", "./schema.graphql", "GraphQL schema in IDL form")
	_ = generateCmd.MarkFlagRequired("schema")
	_ = viper.BindPFlag("schema", generateCmd.Flags().Lookup("schema"))
}
