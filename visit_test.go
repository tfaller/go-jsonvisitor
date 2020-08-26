package jsonvisitor

import (
	"reflect"
	"strings"
	"testing"
)

var testJSON = map[string]interface{}{
	"a": "first",
	"b": []interface{}{
		0,
		"1",
		map[string]interface{}{
			"2": 2,
		},
	},
	"ignore-subs": []interface{}{0, 1},
}

func TestVisit(t *testing.T) {
	// track which pathes got visited
	called := map[string]struct{}{}

	expected := map[string]interface{}{
		"":            testJSON,
		"a":           testJSON["a"],
		"b":           testJSON["b"],
		"b.0":         0,
		"b.1":         "1",
		"b.2":         testJSON["b"].([]interface{})[2],
		"b.2.2":       2,
		"ignore-subs": testJSON["ignore-subs"],
	}

	Visit(testJSON, func(path []string, value interface{}) bool {
		pathStr := strings.Join(path, ".")
		called[pathStr] = struct{}{}

		if !reflect.DeepEqual(value, expected[pathStr]) {
			t.Errorf("Wrong value at %v", pathStr)
		}

		if pathStr == "ignore-subs" {
			// don't visit sub entries
			return false
		}

		return true
	})

	if len(called) != len(expected) {
		t.Error("Missing or too many visits")
	}
}
