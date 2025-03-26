# go-deps

[![Go Reference](https://pkg.go.dev/badge/github.com/matthewmueller/go-deps.svg)](https://pkg.go.dev/github.com/matthewmueller/go-deps)

Walk Go files that your package depends on.

## Features

- Resolves import paths into directories
- Ignores 3rd party and stdlib imports

## Install

```sh
go get github.com/matthewmueller/go-deps
```

## Example

```go
func main() {
	err := deps.Walk(".", func(path string) error {
		paths = append(paths, path)
		return nil
	})
}
```

## Contributors

- Matt Mueller ([@mattmueller](https://twitter.com/mattmueller))

## License

MIT
