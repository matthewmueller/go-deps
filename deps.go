package deps

import (
	"go/build"
	"go/parser"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/matthewmueller/go-deps/internal/mains"
	"github.com/matthewmueller/go-deps/internal/std"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
)

// Find dependencies
func Find(pkgs ...string) (deps []string, err error) {
	if len(pkgs) == 0 {
		return deps, nil
	}

	files, err := mains.Find(pkgs...)
	if err != nil {
		return nil, err
	}

	resolve := func(pkg *build.Package) (files []string) {
		for _, file := range pkg.GoFiles {
			files = append(files, file)
		}
		return files
	}

	return find(resolve, files...)
}

// FindWithTests finds dependencies with tests
func FindWithTests(pkgs ...string) (deps []string, err error) {
	if len(pkgs) == 0 {
		return deps, nil
	}

	files, err := mains.FindTests(pkgs...)
	if err != nil {
		return nil, err
	}

	resolve := func(pkg *build.Package) (files []string) {
		for _, file := range pkg.GoFiles {
			files = append(files, file)
		}
		for _, file := range pkg.TestGoFiles {
			files = append(files, file)
		}
		for _, file := range pkg.XTestGoFiles {
			files = append(files, file)
		}
		return files
	}

	return find(resolve, files...)
}

func find(resolve func(pkg *build.Package) []string, files ...string) (deps []string, err error) {
	var conf loader.Config

	// only parse imports
	conf.ParserMode = parser.ImportsOnly

	// ignore typechecking function bodies
	conf.TypeCheckFuncBodies = func(path string) bool {
		return false
	}

	// go source
	gopath := filepath.Join(build.Default.GOPATH, "src")

	// tap into package resolution
	conf.FindPackage = func(context *build.Context, path, srcDir string, mode build.ImportMode) (*build.Package, error) {
		gosrc := gopath

		if strings.HasPrefix(path, "golang_org/") {
			path = "vendor/" + path
			gosrc = filepath.Join(build.Default.GOROOT, "src")
		}

		pkg, err := context.Import(path, srcDir, mode)
		if err != nil {
			return pkg, err
		}

		// ignore stdlib files
		if std.In(path) {
			return pkg, nil
		}

		// use our resolve function to know what to watch
		for _, file := range resolve(pkg) {
			deps = append(deps, filepath.Join(gosrc, path, file))
		}

		return pkg, nil
	}

	// import all the packages
	for _, file := range files {
		rel, err := filepath.Rel(gopath, file)
		if err != nil {
			return nil, err
		}
		conf.Import(rel)
	}

	// load all the packages
	if _, err := conf.Load(); err != nil {
		return nil, errors.Wrap(err, "unable to load the go package")
	}

	return exists(deps), nil
}

func exists(files []string) (exists []string) {
	l := len(files)

	arr := make([]string, l)
	wg := &sync.WaitGroup{}
	wg.Add(l)

	fn := func(i int, file string) {
		if _, err := os.Stat(file); err == nil {
			arr[i] = file
		}
		wg.Done()
	}

	for i, file := range files {
		go fn(i, file)
	}
	wg.Wait()

	cache := map[string]bool{}
	for _, file := range arr {
		if file == "" || cache[file] {
			continue
		}
		exists = append(exists, file)
		cache[file] = true
	}

	return exists
}
