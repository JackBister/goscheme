package goscheme

import (
	"bytes"
	"fmt"
)

type Vector []Expr

func (v Vector) isExpr() {}

func (v Vector) String() string {
	var b bytes.Buffer
	fmt.Fprint(&b, "#(")
	for i, e := range []Expr(v) {
		fmt.Fprint(&b, e)
		if i != len([]Expr(v))-1 {
			fmt.Fprint(&b, " ")
		}
	}
	fmt.Fprint(&b, ")")
	return b.String()
}

func makevec(e Environment, args ...Expr) Expr {
	i, ok := args[0].(Number)
	if !ok {
		return Error{"make-vector: Argument 1 is not a number."}
	}
	ret := make([]Expr, int(i))
	if len(args) == 2 {
		for i := range ret {
			ret[i] = args[1]
		}
	}
	return Vector(ret)
}

func vector(e Environment, args ...Expr) Expr {
	return Vector(args)
}

func vector_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(Vector)
	return Boolean(ok)
}

func vectorfill(e Environment, args ...Expr) Expr {
	v, ok := args[0].(Vector)
	if !ok {
		return Error{"vector-fill!: Argument 1 is not a vector."}
	}
	vl := []Expr(v)
	for i := range vl {
		vl[i] = args[1]
	}
	return Vector(vl)
}

func vectorlen(e Environment, args ...Expr) Expr {
	v, ok := args[0].(Vector)
	if !ok {
		return Error{"vector-length: Argument 1 is not a vector."}
	}
	return Number(len([]Expr(v)))
}

func vectorref(e Environment, args ...Expr) Expr {
	v, ok := args[0].(Vector)
	if !ok {
		return Error{"vector-ref: Argument 1 is not a vector."}
	}
	i, ok := args[1].(Number)
	if !ok {
		return Error{"vector-ref: Argument 2 is not a number."}
	}
	return []Expr(v)[int(i)]
}

func vectorset(e Environment, args ...Expr) Expr {
	v, ok := args[0].(Vector)
	if !ok {
		return Error{"vector-set!: Argument 1 is not a vector."}
	}
	i, ok := args[1].(Number)
	if !ok {
		return Error{"vector-set!: Argument 2 is not a number."}
	}
	[]Expr(v)[int(i)] = args[2]
	return v
}


func vectolist(e Environment, args ...Expr) Expr {
	v, ok := args[0].(Vector)
	if !ok {
		return Error{"vector->list: Argument 1 is not a vector."}
	}
	return ExprList([]Expr(v))
}
