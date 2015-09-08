package goscheme

import (
	"strconv"
)

type Expr interface {
	isExpr()
}

/*
Number type
The language does not distinguish between precise (int) and imprecise (float)
values right now. 64 bit float should be sufficiently accurate for large ints.
*/
type Number float64
func (n Number) isExpr() {}
func unwrapNumber(n Expr) float64 {
	return float64(n.(Number))
}

/*
Symbol type
The type used for variables. Can also be used as strings as long as it doesn't
contain spaces...
*/
type Symbol string
func (s Symbol) isExpr() {}
func unwrapSymbol(s Expr) string {
	return string(s.(Symbol))
}

type String string
func (s String) isExpr() {}

type Boolean bool
func (b Boolean) isExpr() {}

type Channel chan Expr
func (c Channel) isExpr() {}

//An EvalBlock wraps an expression and delays evaluation of the expr.
//Primarily(only?) used for actions involving apostrophes.
type EvalBlock struct {
	e Expr
}
func (e EvalBlock) isExpr() {}

type ExprList []Expr
func (el ExprList) isExpr() {}

type Func func(...Expr) Expr
func (a Func) isExpr() {}

/*
A Proc is any type that has a function satisfying the function signature.
Built in functions and user defined functions are different, but both satisfy
this interface.
*/
type Proc interface {
	eval(Environment, ...Expr) Expr
}

type UserProc struct {
	//Does this function accept a variable amount of arguments?
	//If true, any excess arguments will be bound to a list stored in the
	//last parameter symbol.
	variadic bool
	//A list of symbols that the arguments to the function will be bound to.
	params ExprList
	body Expr
}
func (u UserProc) isExpr() {}

func (u UserProc) eval(e Environment, args ...Expr) Expr {
	if len(args) < len(u.params) {
		return Error{"Too few arguments (need " + strconv.Itoa(len(u.params)) + ")"}
	}
	if len(args) > len(u.params) && !u.variadic {
		return Error{"Too many arguments (need " + strconv.Itoa(len(u.params)) + ")"}
	}
	for i, par := range u.params {
		e.Local[unwrapSymbol(par)] = args[i]
		if i == len(u.params)-1 && u.variadic {
			if i != 0 && unwrapSymbol(u.params[i-1]) == "." {
				i -= 1
			}
			e.Local[unwrapSymbol(par)] = ExprList(args[i:])
		}
	}
	return Eval(u.body, e)
}

type BuiltIn struct {
	//The name of the function. Used for error printouts.
	name string
	//if maxParams is -1, the function is variadic.
	minParams, maxParams int
	//The go function this struct represents.
	fn func(Environment, ...Expr) Expr
}
func (b BuiltIn) isExpr() {}

func (b BuiltIn) eval(e Environment, args ...Expr) Expr {
	if len(args) < b.minParams {
		return Error{b.name + ": Too few arguments (need " + strconv.Itoa(b.minParams) + ")"}
	}
	if len(args) > b.maxParams && b.maxParams != -1 {
		return Error{b.name + ": Too many arguments (max " + strconv.Itoa(b.maxParams) + ")"}
	}
	return b.fn(e, args...)
}

type Error struct {
	s string
}
func (e Error) Error() string { return e.s }
func (e Error) isExpr() {}


