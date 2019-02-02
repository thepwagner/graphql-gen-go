package generator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTypes(t *testing.T) {
	gen := NewTestGenerator(t)
	err := gen.GenerateTypes()
	assert.NoError(t, err)
}
