# Path

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]

This package is a dependency free utility package.  Use it to construct uri paths from constant templates.

## Examples

```go
func SimpleExample() {
	const pathFoo = "/foo/:foo"
	uri := path.New(pathFoo).
		WithParam(":foo", "bar")
	fmt.Println(uri.String())
	// Output: /foo/bar
}

func AllFeaturesExample() {
	const pathFooBarBaz = "/foo/:foo/bar/:bar/:baz"
	uri := path.New(pathFooBarBaz).
		WithBaseURL("https://example.com").
		WithPrefix("v1").
		WithParam(":foo", "p1").
		WithParams(map[string]string{
			"bar": "p2",
			"baz": "p3",
		}).
		WithQuery("id", "1", "2").
		WithQuery("a", "A").
		WithQueryArgs(map[string]string{
			"b": "B",
			"c": "C",
		})
	fmt.Println(uri.String())
	// Output: https://example.com/v1/foo/p1/bar/p2/p3?a=A&b=B&c=C&id=1&id=2
}
```

## Constructors
```
  New(template string) Path
```

## Constructor Methods
```
  WithBaseURL(url string) Path
  WithPrefix(basePath string) Path
  WithParam(param, value string) Path
  WithParams(params map[string]string) Path
  WithQuery(key string, values ...string) Path
  WithQueryArgs(args map[string]string) Path
  WithQueryValues(query url.Values) Path
```



[build-img]: https://github.com/tempcke/path/actions/workflows/test.yml/badge.svg
[build-url]: https://github.com/tempcke/path/actions
[pkg-img]: https://pkg.go.dev/badge/tempcke/path
[pkg-url]: https://pkg.go.dev/github.com/tempcke/path
[reportcard-img]: https://goreportcard.com/badge/tempcke/path
[reportcard-url]: https://goreportcard.com/report/tempcke/path