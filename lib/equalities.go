package goscheme

import (
	"reflect"
	"strconv"
)

func gt(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{">: Argument 1 is not a number."}
	}
	last := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		if _, ok := args[i].(Number); !ok {
			return Error{">: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		if unwrapNumber(args[i]) >= last {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func lt(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"<: Argument 1 is not a number."}
	}
	last := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		if _, ok := args[i].(Number); !ok {
			return Error{"<: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		if unwrapNumber(args[i]) <= last {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func ge(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{">=: Argument 1 is not a number."}
	}
	last := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		if _, ok := args[i].(Number); !ok {
			return Error{">=: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		if unwrapNumber(args[i]) > last {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func le(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"<=: Argument 1 is not a number."}
	}
	last := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		if _, ok := args[i].(Number); !ok {
			return Error{"<=: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		if unwrapNumber(args[i]) < last {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func eq(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"=: Argument 1 is not a number."}
	}
	last := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		if _, ok := args[i].(Number); !ok {
			return Error{"=: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		if unwrapNumber(args[i]) != last {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

//TODO: char=?, exact/inexact
func eqv(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(ExprList); ok {
		if v2, ok2 := args[1].(ExprList); ok2 {
			return Boolean(listeqv(v, v2))
		}
	}

	if v, ok := args[0].(BuiltIn); ok {
		if v2, ok2 := args[1].(BuiltIn); ok2 {
			return Boolean(builtineqv(v, v2))
		}
	}

	if v, ok := args[0].(UserProc); ok {
		if v2, ok2 := args[1].(UserProc); ok2 {
			return eqv(e, v.body, v2.body).(Boolean)
		}
	}

	return Boolean(args[0] == args[1])
}

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

