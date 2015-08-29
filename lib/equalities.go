package goscheme

import (
	"reflect"
)

//This is a bad approach, but as far as I can tell there is no better way.
func builtineqv(a, b BuiltIn) bool {
	aa := reflect.ValueOf(a.fn)
	bb := reflect.ValueOf(b.fn)
	return aa.Pointer() == bb.Pointer()
}

func listeqv(a, b ExprList) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

