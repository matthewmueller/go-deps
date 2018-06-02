package deps

import (
	"go/build"
	"go/parser"
	"path/filepath"
	"strings"

	"github.com/matthewmueller/deps/internal/mains"
	"github.com/matthewmueller/deps/internal/std"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
)

// Find dependencies
func Find(pkgs ...string) (deps []string, err error) {
	if len(pkgs) == 0 {
		return deps, nil
	}

	var conf loader.Config

	// only parse imports
	conf.ParserMode = parser.ImportsOnly

	// ignore typechecking function bodies
	conf.TypeCheckFuncBodies = func(path string) bool {
		return false
	}

	// go source
	gosrc := filepath.Join(build.Default.GOPATH, "src")

	// tap into package resolution
	conf.FindPackage = func(ctxt *build.Context, importPath, fromDir string, mode build.ImportMode) (*build.Package, error) {
		// ignore stdlib files
		if std.In(importPath) {
			return ctxt.Import(importPath, gosrc, mode)
		}

		// HACK to handle go source's own vendoring
		// here: $GOROOT/src/vendor/
		// e.g. /usr/local/go/src/vendor/
		if strings.HasPrefix(importPath, "golang_org/") {
			importPath = "vendor/" + importPath
			gosrc = filepath.Join(build.Default.GOROOT, "src")
		}

		pkg, err := ctxt.Import(importPath, gosrc, mode)
		if err != nil {
			return pkg, err
		}

		for _, file := range pkg.GoFiles {
			deps = append(deps, filepath.Join(gosrc, importPath, file))
		}

		return pkg, nil
	}

	files, err := mains.Find(pkgs...)
	if err != nil {
		return nil, err
	}

	// import all the packages
	for _, file := range files {
		rel, err := filepath.Rel(gosrc, file)
		if err != nil {
			return nil, err
		}
		conf.Import(rel)
	}

	// load all the packages
	if _, err := conf.Load(); err != nil {
		return nil, errors.Wrap(err, "unable to load the go package")
	}

	return deps, nil
}

// FindWithTests finds dependencies with tests
func FindWithTests(pkgs ...string) (deps []string, err error) {
	return deps, err
}
