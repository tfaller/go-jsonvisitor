package jsonvisitor

import "fmt"

// PairVisitorFunc is the function which is called for
// each found json value pair.
// path is the json pointer of the pair.
// a and b are the values of the one and the other json document.
// If a value exists only in one document, the other value is
// the special value "Undefined"
type PairVisitorFunc func(path []string, a, b interface{}) bool

// undefined is a helper struct to create a unique
// undefined value with a string value of "Undefined"
type undefined struct{ preventCompare func() }

func (undefined) String() string {
	return "Undefined"
}

// Undefined represents a missing value. Only undefined == undefined is true.
// No other value will ever be equal with undefined.
var Undefined interface{} = &undefined{func() {}}

// PairVisit walks over two json trees at the same time in a depth first approach.
// A tree structure is defined by map[string]interface{} and []interface{}
// If a slice index or prop does not exists in the other tree, the walked value
// is replaced with "Undefined". For each found tree element the visitor function
// is called with the values from both trees and the current path.
// The visitor can decide with the return value whether to go deeper in the tree or not.
func PairVisit(a, b interface{}, visitor PairVisitorFunc) {
	PairVisitWithPath([]string{}, a, b, visitor)
}

// PairVisitWithPath works like PairVisit, but a path prefix can be set.
// This is useful if the visitor decides first stop to go depper
// and than resumes the operation. To not lose the path information the
// path can than be passed here. So basically a known path can be set,
// if the here to visit pair is a part of a document.
func PairVisitWithPath(path []string, a, b interface{}, visitor PairVisitorFunc) {
	if a == Undefined && b == Undefined {
		return
	}

	if !visitor(path, a, b) {
		return
	}

	pairVisitArr(path, a, b, visitor)
	pairVisitObj(path, a, b, visitor)
}

func pairVisitObj(path []string, a, b interface{}, visitor PairVisitorFunc) {
	aObj, _ := a.(map[string]interface{})
	bObj, _ := b.(map[string]interface{})

	for k, a := range aObj {
		b, inB := bObj[k]
		if !inB {
			b = Undefined
		}
		PairVisitWithPath(append(path, k), a, b, visitor)
	}

	// props that are in b, but not in a
	for k, b := range bObj {
		if _, inA := aObj[k]; !inA {
			PairVisitWithPath(append(path, k), Undefined, b, visitor)
		}
	}
}

func pairVisitArr(path []string, a, b interface{}, visitor PairVisitorFunc) {
	aArr, _ := a.([]interface{})
	bArr, _ := b.([]interface{})

	aLen, bLen := len(aArr), len(bArr)

	idx := 0
	for ; idx < aLen; idx++ {
		b := Undefined
		if idx < bLen {
			b = bArr[idx]
		}
		PairVisitWithPath(append(path, fmt.Sprint(idx)), aArr[idx], b, visitor)
	}

	// remaining values if b is larger than a
	for ; idx < bLen; idx++ {
		PairVisitWithPath(append(path, fmt.Sprint(idx)), Undefined, bArr[idx], visitor)
	}
}
