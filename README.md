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

## PairVisit
Visit two json documents at the same time - e.g. good for finding differences:

```go
var docA, docB interface{}
json.Unmarshal([]byte(`{"a": 1, "b": [0, 1, 3]}`), &docA)
json.Unmarshal([]byte(`{"a": 1, "b": [0], "c": 3}`), &docB)

jsonvisitor.PairVisit(docA, docB, func(path []string, a, b interface{}) bool {
	fmt.Printf("%v: %v <-> %v", strings.Join(path, "."), a, b)
	return true
})
```
The result:
```
: map[a:1 b:[0 1 3]] <-> map[a:1 b:[0] c:3]
a: 1 <-> 1
b: [0 1 3] <-> [0]
b.0: 0 <-> 0
b.1: 1 <-> Undefined
b.2: 3 <-> Undefined
c: Undefined <-> 3
```