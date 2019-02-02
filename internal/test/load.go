package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/stretchr/testify/require"
)

type Schema string

const (
	StarWars Schema = "swapi"
)

func LoadTestSchema(t *testing.T, testSchema Schema) *ast.Document {
	schemaPath := filepath.Join(findTestDir(t), string(testSchema), "schema.graphql")
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	require.NoError(t, err)
	schema, err := parser.Parse(parser.ParseParams{
		Source: string(schemaBytes),
		Options: parser.ParseOptions{
			NoLocation: false,
			NoSource:   true,
		},
	})
	require.NoError(t, err)
	return schema
}

func SetupOutputDir(t *testing.T, testSchema Schema) string {
	outputDir := filepath.Join(findTestDir(t), string(testSchema), "generated")
	err := os.Mkdir(outputDir, 0700)
	if err != nil {
		require.True(t, os.IsExist(err))
	}
	return outputDir
}

func findTestDir(t *testing.T) string {
	dir, err := os.Getwd()
	require.NoError(t, err)
	projectPath := filepath.Join("github.com", "thepwagner", "magenny")
	split := strings.SplitN(dir, projectPath, 2)
	require.Len(t, split, 2, "could not find GOROOT")
	goRoot := split[0]
	return filepath.Join(goRoot, projectPath, "internal", "test")
}
