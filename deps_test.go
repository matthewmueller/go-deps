package deps_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/go-deps"
)

func TestInner(t *testing.T) {
	is := is.New(t)
	paths := []string{}
	err := deps.Walk(filepath.Join("testdata", "inner"), func(path string) error {
		paths = append(paths, path)
		return nil
	})
	is.NoErr(err)
	is.Equal(len(paths), 3)
	wd, err := os.Getwd()
	is.NoErr(err)
	is.Equal(paths[0], filepath.Join(wd, "testdata", "inner", "inner.go"))
	is.Equal(paths[1], filepath.Join(wd, "testdata", "inner", "innerdep", "another.go"))
	is.Equal(paths[2], filepath.Join(wd, "testdata", "inner", "innerdep", "innerdep.go"))
}
