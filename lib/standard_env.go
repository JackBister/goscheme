package goscheme

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
)

func StandardEnv() map[string]Expr {
	return map[string]Expr {
		"#f": Boolean(false),
		"#t": Boolean(true),
		"+": BuiltIn{add},
		"-": BuiltIn{sub},
		"*": BuiltIn{mul},
		"/": BuiltIn{div},
		">": BuiltIn{gt},
		"<": BuiltIn{lt},
		">=": BuiltIn{ge},
		"<=": BuiltIn{le},
		"=": BuiltIn{eq},
		"<-": BuiltIn{send},
		"->": BuiltIn{receive},
		"abs": BuiltIn{abs},
		"append": BuiltIn{sappend},
		"apply": BuiltIn{apply},
		"begin": BuiltIn{begin},
		"car": BuiltIn{car},
		"cdr": BuiltIn{cdr},
		"chan": BuiltIn{schan},
		"equal?": BuiltIn{eq},
		"length": BuiltIn{length},
		"list": BuiltIn{list},
		"list?": BuiltIn{list_},
		"load": BuiltIn{load},
		"map": BuiltIn{smap},
		"max": BuiltIn{max},
		"min": BuiltIn{min},
		"not": BuiltIn{not},
		"null?": BuiltIn{null_},
		"number?": BuiltIn{number_},
		"procedure?": BuiltIn{procedure_},
		"round": BuiltIn{round},
		"sleep": BuiltIn{sleep},
		"symbol?": BuiltIn{symbol_},
		//TODO: cons, eq?
	}
}

func add(e Environment, args ...Expr) Expr {
	ret := float64(0)
	for i, arg := range args {
		if _, ok := arg.(Number); !ok {
			return Error{"+: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		ret += unwrapNumber(arg)
	}
	return Number(ret)
}

func sub(e Environment, args ...Expr) Expr {
	if len(args) == 0 {
		return Error{"-: Too few arguments (at least 1)."}
	}
	if _, ok := args[0].(Number); !ok {
		return Error{"-: Argument 1 is not a number."}
	}
	ret := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		if _, ok := args[i].(Number); !ok {
			return Error{"-: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		ret -= unwrapNumber(args[i])
	}
	return Number(ret)

}

func mul(e Environment, args ...Expr) Expr {
	ret := float64(1)
	for i, arg := range args {
		if _, ok := arg.(Number); !ok {
			return Error{"*: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		ret *= unwrapNumber(arg)
	}
	return Number(ret)

}

func div(e Environment, args ...Expr) Expr {
	if len(args) == 0 {
		return Error{"/: Too few arguments (at least 1)."}
	}
	if len(args) == 1 {
		if _, ok := args[0].(Number); !ok {
			return Error{"/: Argument 1 is not a number."}
		}
		return Number(1/unwrapNumber(args[0]))
	}
	ret := unwrapNumber(args[0])
	for i := 1; i < len(args); i++ {
		if _, ok := args[i].(Number); !ok {
			return Error{"/: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		ret /= unwrapNumber(args[i])
	}
	return Number(ret)

}

func gt(e Environment, args ...Expr) Expr {
	if len(args) < 2 {
		return Error{">: Too few arguments (at least 2)."}
	}
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
	if len(args) < 2 {
		return Error{"<: Too few arguments (at least 2)."}
	}
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
	if len(args) < 2 {
		return Error{">=: Too few arguments (at least 2)."}
	}
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
	if len(args) < 2 {
		return Error{"<=: Too few arguments (at least 2)."}
	}
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
	if len(args) < 2 {
		return Error{"=: Too few arguments (at least 2)."}
	}
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

func abs(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"abs: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"abs: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(Number); !ok {
		return Error{"abs: Argument 1 is not a number."}
	}
	xf := unwrapNumber(args[0])
	return Number(math.Abs(xf))
}

func sappend(e Environment, args ...Expr) Expr {
	if len(args) < 2 {
		return Error{"append: Too few arguments (at least 2)."}
	}
	ret := make(ExprList, 0)
	for i, arg := range args {
		if _, ok := arg.(ExprList); !ok {
			return Error{"append: Argument " + strconv.Itoa(i+1) + " is not a list."}
		}
		argl := arg.(ExprList)
		ret = append(ret, argl...)
	}
	return ret
}

func apply(e Environment, args ...Expr) Expr {
	if len(args) < 2 {
		return Error{"apply: Too few arguments (at least 2)."}
	}
	//TODO: Errors
	//For example. (apply + 1 (list 1 2)) is a valid call
	proc := args[0].(Proc)
	argl := args[1].(ExprList)
	return proc.eval(e, argl...)
}

func begin(e Environment, args ...Expr) Expr {
	return ExprList(args)[len(args)-1]
}

func car(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"car: Too many arguments (need 1)."}
	}
	if len(args) == 0 {
		return Error{"car: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(ExprList); !ok {
		return Error{"car: Argument 1 is not a list."}
	}
	eList := args[0].(ExprList)
	if len(eList) == 0 {
		return Error{"car: List has length 0"}
	}
	return eList[0]
}

func cdr(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"cdr: Too many arguments (need 1)."}
	}
	if len(args) == 0 {
		return Error{"cdr: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(ExprList); !ok {
		return Error{"cdr: Argument 1 is not a list."}
	}
	eList := args[0].(ExprList)
	if len(eList) < 2 {
		return ExprList{}
	}
	return eList[1:]
}

func length(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"length: Too many arguments (need 1)."}
	}
	if len(args) == 0 {
		return Error{"length: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(ExprList); !ok {
		return Error{"length: Argument 1 is not a list."}
	}
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
	if len(args) > 1 {
		return Error{"list?: Too many arguments (need 1)."}
	}
	if len(args) == 0 {
		return Error{"list?: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(ExprList); !ok {
		return Boolean(false)
	}
	return Boolean(true)
}

func smap(e Environment, args ...Expr) Expr {
	if len(args) > 2 {
		return Error{"map: Too many arguments (need 2)."}
	}
	if len(args) < 0 {
		return Error{"map: Too few arguments (need 2)."}
	}
	if _, ok := args[0].(Proc); !ok {
		return Error{"map: Argument 1 is not a function."}
	}
	if _, ok := args[1].(ExprList); !ok {
		return Error{"map: Argument 2 is not a list."}
	}
	proc := args[0].(Proc)
	eList := args[1].(ExprList)
	ret := make(ExprList, 0)
	for _, exp := range eList {
		ret = append(ret, proc.eval(e, exp))
	}
	return ret
}

func max(e Environment, args ...Expr) Expr {
	if len(args) < 2 {
		return Error{"max: Too few arguments (at least 2)."}
	}
	max := math.Inf(-1)
	eList := args[0].(ExprList)
	for i, arg := range eList {
		if _, ok := arg.(Number); !ok {
			return Error{"max: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		n := unwrapNumber(arg)
		if n > max {
			max = n
		}
	}
	return Number(max)
}

func min(e Environment, args ...Expr) Expr {
	if len(args) < 2 {
		return Error{"min: Too few arguments (at least 2)."}
	}
	min := math.Inf(1)
	eList := args[0].(ExprList)
	for i, arg := range eList {
		if _, ok := arg.(Number); !ok {
			return Error{"min: Argument " + strconv.Itoa(i+1) + " is not a number."}
		}
		n := unwrapNumber(arg)
		if n < min {
			min = n
		}
	}
	return Number(min)
}

func not(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"not: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"not: Too few arguments (need 1)."}
	}

	if _, ok := args[0].(Boolean); !ok {
		if args[0] != nil {
			return Boolean(false)
		}
		return Boolean(true)
	}
	return Boolean(!(bool(args[0].(Boolean))))
}

func null_(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"null?: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"null?: Too few arguments (need 1)."}
	}
	if eList, ok := args[0].(ExprList); ok {
		if len(eList) == 0 {
			return Number(1)
		}
	}
	return Number(0)
}

func number_(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"number?: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"number?: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(Number); !ok {
		return Boolean(false)
	}
	return Boolean(true)
}

func procedure_(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"procedure?: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"procedure?: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(Proc); ok {
		return Boolean(true)
	}
	return Boolean(false)
}

func round(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"round: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"round: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(Number); !ok {
		return Error{"round: Argument 1 is not a number."}
	}
	s := strconv.FormatFloat(unwrapNumber(args[0]), 'f', 0, 64)
	r,_ := strconv.ParseFloat(s, 64)
	return Number(r)
}

func symbol_(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"symbol?: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"symbol?: Too few arguments (need 1)."}
	}
	_, ok := args[0].(Symbol)
	return Boolean(ok)
}

func schan(e Environment, args ...Expr) Expr {
	return Channel(make(chan Expr))
}

func sclose(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"close: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"close: Too few arguments (need 1)."}
	}
	if c, ok := args[0].(Channel); !ok {
		return Error{"close: Argument 1 is not a channel."}
	} else {
		close(c)
	}
	return Boolean(true)
}

func load(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"load: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"load: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(Symbol); !ok {
		return Error{"load: Argument 1 is not a symbol"}
	}
	s := unwrapSymbol(args[0])
	in, err := ioutil.ReadFile(s)
	if err != nil && !strings.HasSuffix(s, ".scm") {
		s += ".scm"
		in, err = ioutil.ReadFile(s)
		if err != nil {
			return Boolean(false)
		}
	}
	ins := string(in)
	ins = strings.Replace(ins, "\t", "", -1)
	ins = strings.Replace(ins, "\n", "", -1)
	ins = strings.Replace(ins, "\r", "", -1)
	t := Tokenize(ins)
	for len(t) != 0 {
		r := Eval(Parse(&t), GlobalEnv)
		if s, ok := r.(Symbol); !ok || string(s) != "" {
			fmt.Println(r)
		}
	}
	return Boolean(true)
}

func receive(e Environment, args ... Expr) Expr {
	if len(args) > 1 {
		return Error{"->: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"->: Too few arguments (need 1)."}
	}
	if c, ok := args[0].(Channel); !ok {
		return Error{"->: Argument 1 is not a channel."}
	} else {
		return <- c
	}
}

func send(e Environment, args ...Expr) Expr {
	if len(args) > 2 {
		return Error{"<-: Too many arguments (max 2)."}
	}
	if len(args) < 2 {
		return Error{"<-: Too few arguments (need 2)."}
	}
	if c, ok := args[0].(Channel); !ok {
		return Error{"<-: Argument 1 is not a channel."}
	} else {
		c <- args[1]
	}
	return args[1]
}

func sleep(e Environment, args ...Expr) Expr {
	if len(args) > 1 {
		return Error{"sleep: Too many arguments (max 1)."}
	}
	if len(args) == 0 {
		return Error{"sleep: Too few arguments (need 1)."}
	}
	if _, ok := args[0].(Number); !ok {
		return Error{"sleep: Argument 1 is not a number."}
	}
	t := time.Duration(unwrapNumber(args[0]))*time.Millisecond
	<-time.After(t)
	return Boolean(true)
}
