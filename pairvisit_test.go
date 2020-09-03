package jsonvisitor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestPairVisit(t *testing.T) {
	jsonRaw1 := `{"a": 0, "b": [1], "c": "1", "e": [1], "ignore": [0]}`
	jsonRaw2 := `{"a": 0, "b": [2], "d": "1", "e": [1, 2], "ignore": [0]}`

	var docA, docB map[string]interface{}

	unmarshalJSONString(jsonRaw1, &docA, t)
	unmarshalJSONString(jsonRaw2, &docB, t)

	called := map[string]struct{}{}
	expected := map[string]struct{ a, b interface{} }{
		"":       {docA, docB},
		"a":      {0., 0.},
		"b":      {docA["b"], docB["b"]},
		"b.0":    {1., 2.},
		"c":      {docA["c"], Undefined},
		"d":      {Undefined, docB["d"]},
		"e":      {docA["e"], docB["e"]},
		"e.0":    {1., 1.},
		"e.1":    {Undefined, 2.},
		"ignore": {docA["ignore"], docB["ignore"]},
	}

	PairVisit(docA, docB, func(path []string, a, b interface{}) bool {
		pathStr := strings.Join(path, ".")
		called[pathStr] = struct{}{}

		expected := expected[pathStr]

		if !reflect.DeepEqual(a, expected.a) ||
			!reflect.DeepEqual(b, expected.b) {
			t.Errorf("Wrong value at %v", pathStr)
		}

		if pathStr == "ignore" {
			return false
		}

		return true
	})

	if len(called) != len(expected) {
		t.Error("Missing or too many visits")
	}
}

func TestUndefinedPair(t *testing.T) {
	PairVisit(Undefined, Undefined, func(path []string, a, b interface{}) bool {
		t.Fatal("If both inputs are undefined this func should not be called")
		return true
	})
}

func TestUndefined(t *testing.T) {
	if Undefined != Undefined {
		t.Error("Undefined must be equal to itself")
	}

	if !reflect.DeepEqual(Undefined, Undefined) {
		t.Error("DeepEqual must see undefined as equal to itself")
	}

	if reflect.DeepEqual(&undefined{func() {}}, Undefined) {
		t.Error("Undefined should not be equal to identical constructed value")
	}

	if fmt.Sprint(Undefined) != "Undefined" {
		t.Error("Wrong string representation of Undefined")
	}
}

// unmarshalJSONString is a simple helper to parse json
func unmarshalJSONString(str string, target interface{}, t *testing.T) {
	err := json.Unmarshal([]byte(str), target)
	if err != nil {
		t.Fatal(err)
	}
}
