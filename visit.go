// Package jsonvisitor is here to simply visit all entries of a json
// document that was parsed like:
//  var doc interface{}
//  json.Unmarshal([]byte(`{"a": 1, "b": [0]}`), &doc)
package jsonvisitor

import "fmt"

// VisitorFunc is the function which is called for
// each found json entry.
// path is the json pointer of the entry.
// value is the entry itself.
type VisitorFunc func(path []string, value interface{}) bool

// Visit searches in a depth-first approach all (sub) entries of
// a json tree. For each found entry the visitor function is called.
// The visitor function can decide with the return value whether the visitor
// should visit also the sub entries. Visit can only see the sub entries of
// map[string]interface{} and []interface{}. These are the results of json.Unmarshal() if
// unmarshalled into interface{}
func Visit(e interface{}, visitor VisitorFunc) {
	VisitWithPath([]string{}, e, visitor)
}

// VisitWithPath works like Visit, but a path prefix can be set.
// This is useful if the visitor decides first stop to go depper
// and than resumes the operation. To not lose the path information the
// path can than be passed here. So basically a known path can be set,
// if the here to visit entry is a part of a document.
func VisitWithPath(path []string, entry interface{}, visitor VisitorFunc) {
	if !visitor(path, entry) {
		return
	}

	// if it's an object ... traverse props
	obj, _ := entry.(map[string]interface{})
	for k, v := range obj {
		VisitWithPath(append(path, k), v, visitor)
	}

	// if it's an array ... traverse indices
	arr, _ := entry.([]interface{})
	for k, v := range arr {
		VisitWithPath(append(path, fmt.Sprint(k)), v, visitor)
	}
}
