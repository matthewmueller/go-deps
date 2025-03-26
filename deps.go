package deps

import (
	"fmt"
	"go/build"
	"path/filepath"

	"github.com/livebud/mod"
)

// Walk a directory calling fn for each Go source file. This function doesn't
// walk into directories that are not part of the module.
func Walk(dir string, fn func(path string) error) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("deps: unable to get absolute path %q. %w", dir, err)
	}
	module, err := mod.Find(absDir)
	if err != nil {
		return fmt.Errorf("deps: unable to find module %q. %w", dir, err)
	}
	seen := map[string]bool{}
	return walk(absDir, fn, module, seen)
}

func walk(dir string, fn func(path string) error, module *mod.Module, seen map[string]bool) error {
	pkg, err := build.ImportDir(dir, build.ImportMode(0))
	if err != nil {
		return fmt.Errorf("deps: unable to import directory %q. %w", dir, err)
	}

	for _, path := range pkg.GoFiles {
		if err := fn(filepath.Join(dir, path)); err != nil {
			return err
		}
	}
	for _, path := range pkg.TestGoFiles {
		if err := fn(filepath.Join(dir, path)); err != nil {
			return err
		}
	}
	for _, path := range pkg.XTestGoFiles {
		if err := fn(filepath.Join(dir, path)); err != nil {
			return err
		}
	}

	seen[module.Import()] = true

	for _, importPath := range pkg.Imports {
		if seen[importPath] || !module.Contains(importPath) {
			continue
		}
		dir, err := module.ResolveDir(importPath)
		if err != nil {
			return fmt.Errorf("deps: unable to resolve directory from import path %q. %w", importPath, err)
		}
		if err := walk(dir, fn, module, seen); err != nil {
			return err
		}
	}

	return nil
}
