# go-jsonvisitor
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tfaller/go-jsonvisitor)](https://pkg.go.dev/github.com/tfaller/go-jsonvisitor)

jsonvisitor is here to simply visit all entries of a json
document that was parsed like:

```go
var doc interface{}
json.Unmarshal([]byte(`{"a": 1, "b": [0]}`), &doc)
```

Visit the entries:
```go
jsonvisitor.Visit(doc, func(path []string, value interface{}) bool {
    fmt.Printf("%v: %v", strings.Join(path, "."), value)
    return true
})
```

The result:
```
: map[a:1 b:[0]]
a: 1
b: [0]
b.0: 0
```