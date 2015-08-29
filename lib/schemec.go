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
	"strconv"
	"strings"
	"time"
)

type Expr interface {
	isExpr()
}

type Number float64
func (n Number) isExpr() {}
func unwrapNumber(n Expr) float64 {
	return float64(n.(Number))
}

type Symbol string
func (s Symbol) isExpr() {}
func unwrapSymbol(s Expr) string {
	return string(s.(Symbol))
}

type Boolean bool
func (b Boolean) isExpr() {}

type Channel chan Expr
func (c Channel) isExpr() {}

type ExprList []Expr
func (el ExprList) isExpr() {}

type Func func(...Expr) Expr
func (a Func) isExpr() {}

type Proc interface {
	eval(Environment, ...Expr) Expr
}

type UserProc struct {
	params ExprList
	body Expr
}
func (u UserProc) eval(e Environment, args ...Expr) Expr {
	if len(args) < len(u.params) {
		return Error{"Too few arguments (need " + strconv.Itoa(len(u.params)) + ")"}
	}
	if len(args) > len(u.params) {
		return Error{"Too many arguments (need " + strconv.Itoa(len(u.params)) + ")"}
	}
	for i, par := range u.params {
		e.Local[unwrapSymbol(par)] = args[i]
	}
	return Eval(u.body, e)
}
func (u UserProc) isExpr() {}

type BuiltIn struct {
	name string
	minParams, maxParams int
	fn func(Environment, ...Expr) Expr
}
func (b BuiltIn) eval(e Environment, args ...Expr) Expr {
	if len(args) < b.minParams {
		return Error{b.name + ": Too few arguments (need " + strconv.Itoa(b.minParams) + ")"}
	}
	if len(args) > b.maxParams && b.maxParams != -1 {
		return Error{b.name + ": Too many arguments (max " + strconv.Itoa(b.maxParams) + ")"}
	}
	return b.fn(e, args...)
}
func (b BuiltIn) isExpr() {}

type Error struct {
	s string
}
func (e Error) Error() string { return e.s }
func (e Error) isExpr() {}

type Environment struct {
	Local map[string]Expr
	Parent *Environment
}
func (e *Environment) find(s string) map[string]Expr {
	if e.Local[s] != nil {
		return e.Local
	}
	if e.Parent == nil {
		return nil
	}
	return e.Parent.find(s)
}
func (e *Environment) copy() Environment {
	nm := map[string]Expr{}
	for k, v := range e.Local {
		nm[k] = v
	}
	return Environment{nm, e.Parent}
}

var GlobalEnv Environment

func Tokenize(s string) []string {
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	ss := strings.Split(strings.Replace(strings.Replace(s, ")", " ) ", -1), "(", " ( ", -1), " ")
	r := make([]string, 0)
	for _, e := range ss {
		if e != " " && e != "" {
			r = append(r, e)
		}
	}
	return r
}

func Parse(s *[]string) Expr {
	if len(*s) == 0 {
		return Error{"Unexpected EOF."}
	}
	t := (*s)[0]
	*s = (*s)[1:]
	if t == "(" {
		l := make(ExprList, 0)
		for (*s)[0] != ")" {
			l = append(l, Parse(s))
			if len(*s) == 0 {
				return Error{"Missing ')'"}
			}
		}
		*s = (*s)[1:]
		return l
	} else if t == ")" {
		return Error{"Unexpected ')'."}
	}
	return atom(t)
}

func Eval(e Expr, env Environment) Expr {
	if bool(symbol_(env, e).(Boolean)) {
		if unwrapSymbol(e)[0] == '\'' {
			return Symbol(unwrapSymbol(e)[1:])
		}
		return env.find(unwrapSymbol(e))[unwrapSymbol(e)]
	} else if !bool(list_(env, e).(Boolean)) {
		return e
	} else if el := e.(ExprList); bool(symbol_(env, el[0]).(Boolean)) {
		if s0 := unwrapSymbol(e.(ExprList)[0]); s0 == "quote" {
			return el[1]
		} else if s0[0] == '\'' {
			return Symbol(s0[1:])
		} else if s0 == "if" {
			r := Eval(el[1], env)
			if _, ok := r.(Boolean); !ok {
				if r != nil {
					r = Boolean(true)
				} else {
					r = Boolean(false)
				}
			}
			if bool(r.(Boolean)) {
				return Eval(el[2], env)
			} else if len(el) > 3 {
				return Eval(el[3], env)
			}
		} else if s0 == "define" {
			er := Eval(el[2], env)
			env.Local[unwrapSymbol(el[1])] = er
			if _, ok := er.(Proc); !ok {
				return er
			}
		} else if s0 == "set!" {
			if env.find(unwrapSymbol(el[1])) != nil {
				er := Eval(el[2], env)
				env.find(unwrapSymbol(el[1]))[unwrapSymbol(el[1])] = er
				return er
			}
			//TODO: Error?
		} else if s0 == "lambda" {
			newenv := Environment{map[string]Expr{}, &env}
			for k, v := range env.Local {
				newenv.Local[k] = v
			}
			return UserProc{el[1].(ExprList), el[2]}
		} else if s0 == "go" {
			c := make(Channel)
			go func(c chan Expr) {
				c <- Eval(el[1], env)
			}(c)
			return c
		} else if s0 == "time" {
			t := time.Now()
			ret := Eval(el[1], env)
			fmt.Println("time:", time.Now().Sub(t))
			return ret
		} else {
			proc := Eval(el[0], env)
			args := make(ExprList, 0)
			for _, arg := range el[1:] {
				args = append(args, Eval(arg, env))
			}
			if procp, ok2 := proc.(Proc); ok2 {
				nEnv := env.copy()
				nEnv.Parent = &env
				return procp.eval(nEnv, args...)
			} else {
				//TODO
				fmt.Println("Error: Expected procedure")
			}
		}
	}
	return Symbol("")
}

func atom(s string) Expr {
	if i, err := strconv.Atoi(s); err == nil {
		return Number((i))
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return Number(f)
	}
	return Symbol(s)
}
