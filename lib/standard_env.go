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
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	//	"runtime"
	"strconv"
	"strings"
	//	"sync"
	"time"
	"unicode"
)

func StandardEnv() Environment {
	e := Environment{map[string]Expr{
		"#f":                  Boolean(false),
		"#t":                  Boolean(true),
		"+":                   BuiltIn{"+", 0, -1, add},
		"-":                   BuiltIn{"-", 1, -1, sub},
		"*":                   BuiltIn{"*", 0, -1, mul},
		"/":                   BuiltIn{"/", 1, -1, div},
		">":                   BuiltIn{">", 2, -1, gt},
		"<":                   BuiltIn{"<", 2, -1, lt},
		">=":                  BuiltIn{">=", 2, -1, ge},
		"<=":                  BuiltIn{"<=", 2, -1, le},
		"=":                   BuiltIn{"=", 2, -1, eq},
		"<-":                  BuiltIn{"<-", 2, 2, send},
		"->":                  BuiltIn{"->", 1, 1, receive},
		"acos":                BuiltIn{"acos", 1, 1, acos},
		"append":              BuiltIn{"append", 2, -1, sappend},
		"apply":               BuiltIn{"apply", 2, -1, apply},
		"asin":                BuiltIn{"asin", 1, 1, asin},
		"atan":                BuiltIn{"atan", 1, 1, atan},
		"begin":               BuiltIn{"begin", 0, -1, begin},
		"char?":               BuiltIn{"char?", 1, 1, char_},
		"close":               BuiltIn{"close", 1, 1, sclose},
		"car":                 BuiltIn{"car", 1, 1, car},
		"cdr":                 BuiltIn{"cdr", 1, 1, cdr},
		"chan":                BuiltIn{"chan", 0, 0, schan},
		"char-alphabetic?":    BuiltIn{"char-alphabetic?", 1, 1, charalpha_},
		"char-downcase":       BuiltIn{"char-downcase", 1, 1, chardown},
		"char-lower-case?":    BuiltIn{"char-lower-case?", 1, 1, charlower_},
		"char-numeric?":       BuiltIn{"char-numeric?", 1, 1, charnumeric_},
		"char-upcase":         BuiltIn{"char-upcase", 1, 1, charup},
		"char-upper-case?":    BuiltIn{"char-upper-case?", 1, 1, charupper_},
		"char-whitespace?":    BuiltIn{"char-whitespace?", 1, 1, charwhitespace_},
		"close-input-port":    BuiltIn{"close-input-port", 1, 1, closeinport},
		"close-output-port":   BuiltIn{"close-output-port", 1, 1, closeoutport},
		"cons":                BuiltIn{"cons", 2, 2, cons},
		"cos":                 BuiltIn{"cos", 1, 1, cos},
		"current-input-port":  Port{os.Stdin, nil, bufio.NewReader(os.Stdin), nil},
		"current-output-port": Port{nil, os.Stdout, nil, bufio.NewWriter(os.Stdout)},
		"exp":              BuiltIn{"exp", 1, 1, exp},
		"eq?":              BuiltIn{"eq?", 2, 2, eqv},
		"equal?":           BuiltIn{"equal?", 2, 2, eq},
		"eqv?":             BuiltIn{"eqv?", 2, 2, eqv},
		"error":            BuiltIn{"error", 1, 1, serror},
		"error?":           BuiltIn{"error?", 1, 1, error_},
		"flush":            BuiltIn{"flush", 0, 1, flush},
		"input-port?":      BuiltIn{"input-port?", 1, 1, inputport_},
		"length":           BuiltIn{"length", 1, 1, length},
		"list":             BuiltIn{"list", 0, -1, list},
		"list?":            BuiltIn{"list?", 1, 1, list_},
		"list->string":     BuiltIn{"list->string", 1, 1, listtostr},
		"load":             BuiltIn{"load", 1, 1, load},
		"log":              BuiltIn{"log", 1, 1, log},
		"max":              BuiltIn{"max", 2, -1, max},
		"min":              BuiltIn{"min", 2, -1, min},
		"modulo":           BuiltIn{"modulo", 2, 2, modulo},
		"newline":          BuiltIn{"newline", 1, 2, newline},
		"not":              BuiltIn{"not", 1, 1, not},
		"null?":            BuiltIn{"null?", 1, 1, null_},
		"number?":          BuiltIn{"number?", 1, 1, number_},
		"open-input-file":  BuiltIn{"open-input-file", 1, 1, openinfile},
		"open-output-file": BuiltIn{"open-output-file", 1, 1, openoutfile},
		"output-port?":     BuiltIn{"output-port?", 1, 1, outputport_},
		"peek-char":        BuiltIn{"peek-char", 0, 1, peekchar},
		//"pmap": BuiltIn{"pmap", 2, -1, pmap},
		"procedure?":   BuiltIn{"procedure?", 1, 1, procedure_},
		"read-char":    BuiltIn{"read-char", 0, 1, readchar},
		"remainder":    BuiltIn{"remainder", 2, 2, remainder},
		"round":        BuiltIn{"round", 1, 1, round},
		"sin":          BuiltIn{"sin", 1, 1, sin},
		"sleep":        BuiltIn{"sleep", 1, 1, sleep},
		"string->list": BuiltIn{"string->list", 1, 1, strtolist},
		"string?":      BuiltIn{"string?", 1, 1, string_},
		"symbol?":      BuiltIn{"symbol?", 1, 1, symbol_},
		"tan":          BuiltIn{"tan", 1, 1, tan},
		"write":        BuiltIn{"write", 1, 2, write},
		"write-char":   BuiltIn{"write-char", 1, 2, writechar},
		//TODO: eq?
	}, nil}
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
			r := Eval(Parse(&t, true), e)
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
	if len(args) == 1 {
		return Number(-ret)
	}
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
		return Number(1 / unwrapNumber(args[0]))
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
	proc := args[0].(Proc)
	argn := args[1 : len(args)-1]
	argl := args[len(args)-1].(ExprList)
	return proc.eval(e, append(argn, argl...)...)
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

func char_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(Character)
	return Boolean(ok)
}

func charalpha_(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Character); !ok {
		return Error{"char-alphabetic?: Argument 1 is not a char."}
	} else {
		return Boolean(unicode.IsLetter(rune(v)))
	}
}

func chardown(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Character); !ok {
		return Error{"char-downcase?: Argument 1 is not a char."}
	} else {
		return Character(unicode.ToLower(rune(v)))
	}
}

func charlower_(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Character); !ok {
		return Error{"char-lower-case?: Argument 1 is not a char."}
	} else {
		return Boolean(unicode.IsLower(rune(v)))
	}
}

func charnumeric_(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Character); !ok {
		return Error{"char-numeric?: Argument 1 is not a char."}
	} else {
		return Boolean(unicode.IsNumber(rune(v)))
	}
}

func charup(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Character); !ok {
		return Error{"char-upcase?: Argument 1 is not a char."}
	} else {
		return Character(unicode.ToUpper(rune(v)))
	}
}

func charupper_(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Character); !ok {
		return Error{"char-upper-case?: Argument 1 is not a char."}
	} else {
		return Boolean(unicode.IsUpper(rune(v)))
	}
}

func charwhitespace_(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Character); !ok {
		return Error{"char-whitespace?: Argument 1 is not a char."}
	} else {
		return Boolean(unicode.IsSpace(rune(v)))
	}
}

func cons(e Environment, args ...Expr) Expr {
	if l, ok := args[1].(ExprList); ok {
		return append(ExprList{args[0]}, l...)
	}
	return ExprList{args[0], args[1]}
}

func exp(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"exp: Argument 1 is not a number."}
	} else {
		return Number(math.Exp(float64(v)))
	}
}

func serror(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(String); !ok {
		return Error{"error: Argument 1 is not a string."}
	} else {
		return Error{string(v)}
	}
}

func error_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(Error)
	return Boolean(ok)
}

func flush(e Environment, args ...Expr) Expr {
	var ep Expr
	if len(args) == 1 {
		ep = args[0]
	} else {
		ep = e.Local["current-output-port"]
	}
	p, ok := ep.(Port)
	if !ok && p.w == nil {
		return Error{"flush: Not an output port."}
	}
	p.w.Flush()
	return Boolean(true)
}

func inputport_(e Environment, args ...Expr) Expr {
	p, ok := args[0].(Port)
	return Boolean(ok && p.r != nil)
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
	_, ok := args[0].(ExprList)
	return Boolean(ok)
}

func listtostr(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(ExprList); !ok {
		return Error{"list->string: Argument 1 is not a list."}
	}
	l := args[0].(ExprList)
	s := make([]rune, len(l))
	for i, v := range l {
		if _, ok2 := v.(Character); !ok2 {
			return Error{"list->string: All members of list must be characters."}
		}
		//Icky, but cannot cast directly to rune because it's not an Expr.
		c := v.(Character)
		s[i] = rune(c)
	}
	return String(s)
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

func newline(e Environment, args ...Expr) Expr {
	var ep Expr
	if len(args) == 1 {
		ep = args[0]
	} else {
		ep = e.Local["current-output-port"]
	}
	p, ok := ep.(Port)
	if !ok && p.w == nil {
		return Error{"newline: Not an output port."}
	}
	fmt.Fprintln(p.w, "")
	return Boolean(true)
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
	eList, ok := args[0].(ExprList)
	return Boolean(ok && len(eList) == 0)
}

func number_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(Number)
	return Boolean(ok)
}

func openinfile(e Environment, args ...Expr) Expr {
	if s, ok := args[0].(String); !ok {
		return Error{"open-input-file: Argument 1 is not a string."}
	} else {
		f, err := os.Open(string(s))
		if err != nil {
			return Error{err.Error()}
		}
		return Port{f, nil, bufio.NewReader(f), nil}
	}
}

func openoutfile(e Environment, args ...Expr) Expr {
	if s, ok := args[0].(String); !ok {
		return Error{"open-output-file: Argument 1 is not a string."}
	} else {
		f, err := os.Create(string(s))
		if err != nil {
			return Error{err.Error()}
		}
		return Port{nil, f, nil, bufio.NewWriter(f)}
	}
}

func outputport_(e Environment, args ...Expr) Expr {
	p, ok := args[0].(Port)
	return Boolean(ok && p.w != nil)
}

func peekchar(e Environment, args ...Expr) Expr {
	var p Port
	if len(args) == 0 {
		p, _ = e.Local["current-input-port"].(Port)
	} else if p2, ok := args[0].(Port); !ok {
		return Error{"peek-char: Argument 1 is not a port."}
	} else {
		p = p2
	}
	if p.r == nil {
		return Error{"peek-char: Not an input port."}
	}
	r, _, err := p.r.ReadRune()
	if err != nil {
		return Error{err.Error()}
	}
	p.r.UnreadRune()
	return Character(r)
}

func procedure_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(Proc)
	return Boolean(ok)
}

func readchar(e Environment, args ...Expr) Expr {
	var ep Expr
	if len(args) == 1 {
		ep = args[0]
	} else {
		ep = e.Local["current-input-port"]
	}
	p, ok := ep.(Port)
	if !ok && p.w == nil {
		return Error{"read-char: Not an input port."}
	}
	r, _, err := p.r.ReadRune()
	if err != nil {
		return Error{err.Error()}
	}
	return Character(r)
}

func round(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"round: Argument 1 is not a number."}
	}
	s := strconv.FormatFloat(unwrapNumber(args[0]), 'f', 0, 64)
	r, _ := strconv.ParseFloat(s, 64)
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

func closeinport(e Environment, args ...Expr) Expr {
	if p, ok := args[0].(Port); !ok {
		return Error{"close-input-port: Argument 1 is not a port."}
	} else {
		if p.r == nil {
			return Error{"close-input-port: Not an input port."}
		}
		p.r = nil
		p.rclose.Close()
		return Boolean(true)
	}
}

func closeoutport(e Environment, args ...Expr) Expr {
	if p, ok := args[0].(Port); !ok {
		return Error{"close-output-port: Argument 1 is not a port."}
	} else {
		if p.w == nil {
			return Error{"close-output-port: Not an output port."}
		}
		p.w.Flush()
		p.w = nil
		p.wclose.Close()
		return Boolean(true)
	}
}

//TODO: Could allow loading multiple files in one call.
func load(e Environment, args ...Expr) Expr {
	var s string
	if v, ok := args[0].(Symbol); !ok {
		if v2, ok2 := args[0].(String); !ok2 {
			return Error{"load: Argument 1 is not a string"}
		} else {
			s = string(v2)
		}
	} else {
		s = string(v)
	}
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
		r := Eval(Parse(&t, true), GlobalEnv)
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

func receive(e Environment, args ...Expr) Expr {
	if c, ok := args[0].(Channel); !ok {
		return Error{"->: Argument 1 is not a channel."}
	} else {
		return <-c
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
	t := time.Duration(unwrapNumber(args[0])) * time.Millisecond
	<-time.After(t)
	return Boolean(true)
}

func strtolist(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(String); !ok {
		return Error{"string->list: Argument 1 is not a string."}
	} else {
		r := make([]Expr, len(v))
		for i, c := range v {
			r[i] = Character(c)
		}
		return ExprList(r)
	}
}

func string_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(String)
	return Boolean(ok)
}

func write(e Environment, args ...Expr) Expr {
	var ep Expr
	if len(args) == 2 {
		ep = args[1]
	} else {
		ep = e.Local["current-output-port"]
	}
	p, ok := ep.(Port)
	if !ok && p.w == nil {
		return Error{"write: Not an output port."}
	}
	fmt.Fprint(p.w, args[0])
	return Boolean(true)
}

func writechar(e Environment, args ...Expr) Expr {
	var ep Expr
	if len(args) == 2 {
		ep = args[1]
	} else {
		ep = e.Local["current-output-port"]
	}
	p, ok := ep.(Port)
	if !ok && p.w == nil {
		return Error{"write-char: Not an output port."}
	}
	if c, ok2 := args[0].(Character); !ok2 {
		return Error{"write-char: Argument 1 is not a character."}
	} else {
		_, err := p.w.WriteRune(rune(c))
		if err != nil {
			return Error{err.Error()}
		}
	}
	return Boolean(true)
}

/*
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
*/
