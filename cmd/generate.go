package cmd

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thepwagner/graphql-gen-go/internal/generator"
)

var gqlSchemaPath string
var gqlQueryDir string
var goPackage string

var generateCmd = &cobra.Command{
	Use: "generate",
	RunE: func(cmd *cobra.Command, args []string) error {
		schema, err := loadSchema(gqlSchemaPath)
		if err != nil {
			return errors.Wrap(err, "loading schema")
		}
		queries, err := loadQueries(gqlQueryDir)
		if err != nil {
			return errors.Wrap(err, "loading queries")
		}
		dir, err := setupOutputDir(goPackage, "")
		if err != nil {
			return errors.Wrap(err, "setting up output dir")
		}
		if err := generator.Generate(dir, schema, queries...); err != nil {
			return errors.Wrap(err, "generating code")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVar(&gqlSchemaPath, "schema", "./schema.graphql", "GraphQL schema in IDL form")
	_ = viper.BindPFlag("schema", generateCmd.Flags().Lookup("schema"))

	generateCmd.Flags().StringVar(&gqlQueryDir, "queries", "./", "Directory containing GraphQL query files")
	_ = viper.BindPFlag("queries", generateCmd.Flags().Lookup("queries"))

	generateCmd.Flags().StringVar(&goPackage, "pkg", "", "Go package for generated code")
	_ = generateCmd.MarkFlagRequired("pkg")
	_ = viper.BindPFlag("pkg", generateCmd.Flags().Lookup("pkg"))
}

func loadSchema(path string) (*ast.Document, error) {
	log.WithField("schema", path).Debug("Loading schema.")
	schemaBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	schema, err := parser.Parse(parser.ParseParams{
		Source: string(schemaBytes),
		Options: parser.ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})
	if err != nil {
		return nil, err
	}
	log.WithField("definitions", len(schema.Definitions)).Debug("Loaded schema.")

	return schema, err
}

func loadQueries(path string) ([]string, error) {
	// TODO: what about dynamic queries?
	// i.e. i want to look up N nodes using aliases returned in a map[string]Node

	log.WithField("dir", path).Debug("Loading queries.")
	ls, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var queries []string
	for _, f := range ls {
		n := f.Name()
		if f.IsDir() || n == "schema.graphql" || !strings.HasSuffix(n, ".graphql") {
			continue
		}
		d, err := ioutil.ReadFile(filepath.Join(path, n))
		if err != nil {
			return nil, err
		}
		queries = append(queries, string(d))
	}

	log.WithField("queries", len(queries)).Debug("Loaded queries.")
	return queries, err
}

func setupOutputDir(path, suffix string) (string, error) {
	// The package may not exist yet, recurse to find the nearest parent that exists
	if path == "" {
		return "", fmt.Errorf("could not find package %s", suffix)
	}
	pkg, err := build.Import(path, build.Default.SrcDirs()[0], build.IgnoreVendor)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find package") {
			parent, name := filepath.Split(path)
			return setupOutputDir(strings.TrimSuffix(parent, "/"), filepath.Join(name, suffix))
		}
		return "", err
	}

	// If we had to recurse, create directories so the package exists
	if suffix != "" {
		toCreate := filepath.Join(pkg.Dir, suffix)
		if err := os.MkdirAll(toCreate, 0770); err != nil {
			return "", err
		}
		return toCreate, nil
	}
	return pkg.Dir, nil
}
