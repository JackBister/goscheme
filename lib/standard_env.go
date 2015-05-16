package schemec 

import (
	"math"
	"reflect"
)

func StandardEnv() map[string]Expr {
	return map[string]Expr {
		"+": Func(add),
		"-": Func(sub),
		"*": Func(mul),
		"/": Func(div),
		">": Func(gt),
		"<": Func(lt),
		">=": Func(ge),
		"<=": Func(le),
		"=": Func(eq),
		"abs": Func(abs),
		//TODO:append,apply,begin,car,cdr,cons,eq?
		"car": Func(car),
		"cdr": Func(cdr),
		"equal?": Func(eq),
		"length": Func(length),
		"list": Func(list),
		"list?": Func(list_),
		"null?": Func(null_),
		"number?": Func(number_),
		"symbol?": Func(symbol_),
	}
}

func typeOf(e Expr) reflect.Kind {
	return reflect.TypeOf(e).Kind()
}

func add(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Expr(Number(args[0].(Number)+args[1].(Number)))
}

func sub(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Expr(Number(args[0].(Number)-args[1].(Number)))
}

func mul(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Expr(Number(args[0].(Number)*args[1].(Number)))
}

func div(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Expr(Number(args[0].(Number)/args[1].(Number)))
}

func gt(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0].(Number) > args[1].(Number))
}

func lt(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0].(Number) < args[1].(Number)) 
}

func ge(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0].(Number) >= args[1].(Number))
}

func le(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0].(Number) <= args[1].(Number))
}

func eq(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 || typeOf(args[1]) != reflect.Float64 {
		//TODO: Error
	}
	return Boolean(args[0] == args[1])
}

func abs(args ...Expr) Expr {
	if typeOf(args[0]) != reflect.Float64 {
		//TODO: Error
	}
	xf := float64(args[0].(Number))
	return Expr(Number(math.Abs(xf)))
}

func car(args ...Expr) Expr {
	//TODO: Error check
	return args[0].(ExprList)[0]
}

func cdr(args ...Expr) Expr {
	//TODO: Error check
	return args[0].(ExprList)[1:]
}

func length(args ...Expr) Expr {
	//TODO: Error check
	return Number(len(args[0].(ExprList)))
}

func list(args ...Expr) Expr {
	ret := make(ExprList, 0)
	for _, e := range args {
		ret = append(ret, e)
	}
	return ret
}

func list_(args ...Expr) Expr {
	return Boolean(reflect.TypeOf(args[0]).Name() == "ExprList")
}

func null_(args ...Expr) Expr {
	if len(args[0].(ExprList)) == 0 {
		return Number(1)
	}
	return Number(0)
}

func number_(args ...Expr) Expr {
	if reflect.TypeOf(args[0]).Implements(reflect.TypeOf(Number(0))) {
		return Number(1)
	}
	return Number(0)
}

func symbol_(args ...Expr) Expr {
	if _, ok := args[0].(Symbol); ok {
		return Number(1)
	}
	return Number(0)
}
