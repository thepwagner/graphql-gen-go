package swapi_test

import (
	"fmt"
	"testing"

	"github.com/thepwagner/magenny/internal/test/swapi/generated"
)

func TestAPI(t *testing.T) {
	f := ""
	r2 := generated.Person{
		Name: &f,
	}
	fmt.Printf("%+v\n", r2)
}
