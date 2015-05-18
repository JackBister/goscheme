package schemec 

import (
	"math"
	"reflect"
	"strconv"
)

func StandardEnv() map[string]Expr {
	return map[string]Expr {
		"+": BuiltIn{add},
		"-": BuiltIn{sub},
		"*": BuiltIn{mul},
		"/": BuiltIn{div},
		">": BuiltIn{gt},
		"<": BuiltIn{lt},
		">=": BuiltIn{ge},
		"<=": BuiltIn{le},
		"=": BuiltIn{eq},
		"abs": BuiltIn{abs},
		"append": BuiltIn{sappend},
		"apply": BuiltIn{apply},
		"begin": BuiltIn{begin},
		"car": BuiltIn{car},
		"cdr": BuiltIn{cdr},
		"equal?": BuiltIn{eq},
		"length": BuiltIn{length},
		"list": BuiltIn{list},
		"list?": BuiltIn{list_},
		"map": BuiltIn{smap},
		"max": BuiltIn{max},
		"min": BuiltIn{min},
		"null?": BuiltIn{null_},
		"number?": BuiltIn{number_},
		"procedure?": BuiltIn{procedure_},
		"round": BuiltIn{round},
		"symbol?": BuiltIn{symbol_},
		//TODO: cons, eq?, not
	}
}

func typeOf(e Expr) reflect.Kind {
	return reflect.TypeOf(e).Kind()
}

func add(e Environment, args ...Expr) Expr {
	ret := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		ret += unwrapNumber(args[i])
	}
	return Number(ret)
}

func sub(e Environment, args ...Expr) Expr {
	ret := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		ret -= unwrapNumber(args[i])
	}
	return Number(ret)
}

func mul(e Environment, args ...Expr) Expr {
	ret := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		ret *= unwrapNumber(args[i])
	}
	return Number(ret)

}

func div(e Environment, args ...Expr) Expr {
	ret := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		ret /= unwrapNumber(args[i])
	}
	return Number(ret)

}

func gt(e Environment, args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0].(Number) > args[1].(Number))
}

func lt(e Environment, args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0].(Number) < args[1].(Number)) 
}

func ge(e Environment, args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0].(Number) >= args[1].(Number))
}

func le(e Environment, args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0].(Number) <= args[1].(Number))
}

func eq(e Environment, args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	cmp := args[0]
	for _, arg := range args {
		if arg != cmp {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func abs(e Environment, args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 {
		//TODO: Error
	}
	xf := float64(args[0].(Number))
	return Expr(Number(math.Abs(xf)))
}

func sappend(e Environment, args ...Expr) Expr {
	ret := make(ExprList, 0)
	for _, arg := range args {
		argl := arg.(ExprList)
		ret = append(ret, argl...)
	} 
	return ret
}

func apply(e Environment, args ...Expr) Expr {
	proc := args[0].(Proc)
	argl := args[1].(ExprList)
	return proc.eval(e, argl...)
}

func begin(e Environment, args ...Expr) Expr {
	return ExprList(args)[len(args)-1]
}

func car(e Environment, args ...Expr) Expr {
	//TODO: Error check
	return args[0].(ExprList)[0]
}

func cdr(e Environment, args ...Expr) Expr {
	//TODO: Error check
	return args[0].(ExprList)[1:]
}

func length(e Environment, args ...Expr) Expr {
	//TODO: Error check
	return Number(len(args[0].(ExprList)))
}

func list(e Environment, args ...Expr) Expr {
	ret := make(ExprList, 0)
	for _, e := range args {
		ret = append(ret, e)
	}
	return ret
}

func list_(e Environment, args ...Expr) Expr {
	return Boolean(reflect.TypeOf(args[0]).Name() == "ExprList")
}

func smap(e Environment, args ...Expr) Expr {
	proc := args[0].(Proc)
	eList := args[1].(ExprList)
	ret := make(ExprList, 0)
	for _, exp := range eList {
		ret = append(ret, proc.eval(e, exp))
	}
	return ret
}

func max(e Environment, args ...Expr) Expr {
	//TODO: Error
	max := math.Inf(-1)
	eList := args[0].(ExprList)
	for _, arg := range eList {
		n := unwrapNumber(arg)
		if n > max {
			max = n
		}
	}
	return Number(max)
}

func min(e Environment, args ...Expr) Expr {
	min := math.Inf(1)
	eList := args[0].(ExprList)
	for _, arg := range eList {
		n := unwrapNumber(arg)
		if n < min {
			min = n
		}
	}
	return Number(min)
}

func null_(e Environment, args ...Expr) Expr {
	if len(args[0].(ExprList)) == 0 {
		return Number(1)
	}
	return Number(0)
}

func number_(e Environment, args ...Expr) Expr {
	if reflect.TypeOf(args[0]).Implements(reflect.TypeOf(Number(0))) {
		return Number(1)
	}
	return Number(0)
}

func procedure_(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Proc); ok {
		return Boolean(true)
	}
	return Boolean(false)
}

func round(e Environment, args ...Expr) Expr {
	s := strconv.FormatFloat(unwrapNumber(args[0]), 'f', 0, 64)
	r,_ := strconv.ParseFloat(s, 64)
	return Number(r)
}

func symbol_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(Symbol)
	return Boolean(ok)
}
