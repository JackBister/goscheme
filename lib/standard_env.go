/*
    goscheme - a Lisp interpreter in Go
    Copyright (C) 2015 Jack Bister

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package goscheme

import (
	"fmt"
	"io/ioutil"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func StandardEnv() Environment {
	e := Environment{ map[string]Expr {
		"#f": Boolean(false),
		"#t": Boolean(true),
		"+": BuiltIn{"+", 0, -1, add},
		"-": BuiltIn{"-", 1, -1, sub},
		"*": BuiltIn{"*", 0, -1, mul},
		"/": BuiltIn{"/", 1, -1, div},
		">": BuiltIn{">", 2, -1, gt},
		"<": BuiltIn{"<", 2, -1, lt},
		">=": BuiltIn{">=", 2, -1, ge},
		"<=": BuiltIn{"<=", 2, -1, le},
		"=": BuiltIn{"=", 2, -1, eq},
		"<-": BuiltIn{"<-", 2, 2, send},
		"->": BuiltIn{"->", 1, 1, receive},
		"abs": BuiltIn{"abs", 1, 1, abs},
		"acos": BuiltIn{"acos", 1, 1, acos},
		"append": BuiltIn{"append", 2, -1, sappend},
		"apply": BuiltIn{"apply", 2, -1, apply},
		"asin": BuiltIn{"asin", 1, 1, asin},
		"atan": BuiltIn{"atan", 1, 1, atan},
		"begin": BuiltIn{"begin", 0, -1, begin},
		"close": BuiltIn{"close", 1, 1, sclose},
		"car": BuiltIn{"car", 1, 1, car},
		"cdr": BuiltIn{"cdr", 1, 1, cdr},
		"chan": BuiltIn{"chan", 0, 0, schan},
		"cons": BuiltIn{"cons", 2, 2, cons},
		"cos": BuiltIn{"cos", 1, 1, cos},
		"exp": BuiltIn{"exp", 1, 1, exp},
		"eq?": BuiltIn{"eq?", 2, 2, eqv},
		"equal?": BuiltIn{"equal?", 2, 2, eq},
		"eqv?": BuiltIn{"eqv?", 2, 2, eqv},
		"length": BuiltIn{"length", 1, 1, length},
		"list": BuiltIn{"list", 0, -1, list},
		"list?": BuiltIn{"list?", 1, 1, list_},
		"load": BuiltIn{"load", 1, 1, load},
		"log": BuiltIn{"log", 1, 1, log},
		"map": BuiltIn{"map", 2, -1, smap},
		"max": BuiltIn{"max", 2, -1, max},
		"min": BuiltIn{"min", 2, -1, min},
		"remainder": BuiltIn{"remainder", 2, 2, remainder},
		"not": BuiltIn{"not", 1, 1, not},
		"null?": BuiltIn{"null?", 1, 1, null_},
		"number?": BuiltIn{"number?", 1, 1, number_},
		"pmap": BuiltIn{"pmap", 2, -1, pmap},
		"procedure?": BuiltIn{"procedure?", 1, 1, procedure_},
		"round": BuiltIn{"round", 1, 1, round},
		"sin": BuiltIn{"sin", 1, 1, sin},
		"sleep": BuiltIn{"sleep", 1, 1, sleep},
		"symbol?": BuiltIn{"symbol?", 1, 1, symbol_},
		"tan": BuiltIn{"tan", 1, 1, tan},
		//TODO: eq?
	}, nil }
	dirc, err := ioutil.ReadDir("std")
	if err != nil {
		panic("Error while loading standard library")
	}
	for _, fi := range dirc {
		if !strings.HasSuffix(fi.Name(), ".scm") {
			continue
		}
		in, err := ioutil.ReadFile("std/" + fi.Name())
		if err != nil {
			panic("Error while loading standard library")
		}
		ins := string(in)
		t := Tokenize(ins)
		for len(t) != 0 {
			r := Eval(Parse(&t), e)
			if s, ok := r.(Symbol); !ok || string(s) != "" {
				fmt.Println(r)
			}
		}
	}
	return e
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

func abs(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"abs: Argument 1 is not a number."}
	}
	xf := unwrapNumber(args[0])
	return Number(math.Abs(xf))
}

func sappend(e Environment, args ...Expr) Expr {
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
	//TODO: Errors
	//For example. (apply + 1 (list 1 2)) is a valid call
	proc := args[0].(Proc)
	argl := args[1].(ExprList)
	return proc.eval(e, argl...)
}

//TODO: 0 args => return the begin proc
func begin(e Environment, args ...Expr) Expr {
	return ExprList(args)[len(args)-1]
}

func car(e Environment, args ...Expr) Expr {
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
	if _, ok := args[0].(ExprList); !ok {
		return Error{"cdr: Argument 1 is not a list."}
	}
	eList := args[0].(ExprList)
	if len(eList) < 2 {
		return ExprList{}
	}
	return eList[1:]
}

func cons(e Environment, args ...Expr) Expr {
	return ExprList{args[0], args[1]}
}

func exp(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"exp: Argument 1 is not a number."}
	} else {
		return Number(math.Exp(float64(v)))
	}
}

func length(e Environment, args ...Expr) Expr {
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
	if _, ok := args[0].(ExprList); !ok {
		return Boolean(false)
	}
	return Boolean(true)
}

func smap(e Environment, args ...Expr) Expr {
	//TODO: Handle multiple list parameters, e.g.
	//(map + (list 1 2) (list 3 4)) => [4 6]
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
		nEnv := e.copy()
		nEnv.Parent = &e
		if uprocp, ok := proc.(UserProc); ok {
			nEnv.Local[unwrapSymbol(uprocp.params[0])] = exp
		}
		ret = append(ret, proc.eval(nEnv, exp))
	}
	return ret
}

func max(e Environment, args ...Expr) Expr {
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
	if _, ok := args[0].(Boolean); !ok {
		if args[0] != nil {
			return Boolean(false)
		}
		return Boolean(true)
	}
	return Boolean(!(bool(args[0].(Boolean))))
}

func null_(e Environment, args ...Expr) Expr {
	if eList, ok := args[0].(ExprList); ok {
		if len(eList) == 0 {
			return Number(1)
		}
	}
	return Number(0)
}

func number_(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Boolean(false)
	}
	return Boolean(true)
}

func procedure_(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Proc); ok {
		return Boolean(true)
	}
	return Boolean(false)
}

func round(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"round: Argument 1 is not a number."}
	}
	s := strconv.FormatFloat(unwrapNumber(args[0]), 'f', 0, 64)
	r,_ := strconv.ParseFloat(s, 64)
	return Number(r)
}

func symbol_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(Symbol)
	return Boolean(ok)
}

func schan(e Environment, args ...Expr) Expr {
	return Channel(make(chan Expr))
}

func sclose(e Environment, args ...Expr) Expr {
	if c, ok := args[0].(Channel); !ok {
		return Error{"close: Argument 1 is not a channel."}
	} else {
		close(c)
	}
	return Boolean(true)
}

//TODO: Could allow loading multiple files in one call.
func load(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Symbol); !ok {
		return Error{"load: Argument 1 is not a symbol"}
	}
	s := unwrapSymbol(args[0])
	fmt.Println("Reading file " + s + "...")
	in, err := ioutil.ReadFile(s)
	if err != nil && !strings.HasSuffix(s, ".scm") {
		s += ".scm"
		fmt.Println("Not found, reading file " + s + "...")
		in, err = ioutil.ReadFile(s)
		if err != nil {
			return Boolean(false)
		}
	}
	ins := string(in)
	t := Tokenize(ins)
	for len(t) != 0 {
		r := Eval(Parse(&t), GlobalEnv)
		if s, ok := r.(Symbol); !ok || string(s) != "" {
			fmt.Println(r)
		}
	}
	return Boolean(true)
}

func log(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"log: Argument 1 is not a number"}
	} else {
		return Number(math.Log(float64(v)))
	}
}

func receive(e Environment, args ... Expr) Expr {
	if c, ok := args[0].(Channel); !ok {
		return Error{"->: Argument 1 is not a channel."}
	} else {
		return <- c
	}
}

func send(e Environment, args ...Expr) Expr {
	if c, ok := args[0].(Channel); !ok {
		return Error{"<-: Argument 1 is not a channel."}
	} else {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("<-: Attempt to send on a closed channel")
			}
		}()
		c <- args[1]
	}
	return args[1]
}

func sleep(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"sleep: Argument 1 is not a number."}
	}
	t := time.Duration(unwrapNumber(args[0]))*time.Millisecond
	<-time.After(t)
	return Boolean(true)
}

func pmap(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Proc); !ok {
		return Error{"pmap: Argument 1 is not a function."}
	}
	if _, ok := args[1].(ExprList); !ok {
		return Error{"pmap: Argument 2 is not a list."}
	}
	proc := args[0].(Proc)
	eList := args[1].(ExprList)
	maxprocs := runtime.GOMAXPROCS(0)
	procs := maxprocs
	if len(eList) < maxprocs {
		procs = len(eList)
	}
	var wg sync.WaitGroup
	wg.Add(procs)
	ret := make(ExprList, len(eList))
	for i := 0; i < procs; i++ {
		go func(j int) {
			//length := int(math.Ceil(float64(len(eList))/float64(procs)))
			length := len(eList)/procs
			start := j*length
			end := start + length
			if j == procs-1 && end < len(ret) {
				end = len(ret)
			}
			r := smap(e, proc.(Expr), eList[start:end])
			rlist := r.(ExprList)
			for k := 0; k < len(rlist); k++ {
				ret[start+k] = rlist[k]
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return ret
}
