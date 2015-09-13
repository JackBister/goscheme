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
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Environment struct {
	Local  map[string]Expr
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
var wsReplacer, _ = regexp.Compile("\t|\n|\r|;.*?\n")

func Tokenize(s string) []string {
	s = wsReplacer.ReplaceAllString(s, "")
	ss := strings.Split(strings.Replace(strings.Replace(strings.Replace(strings.Replace(s, ")", " ) ", -1), "(", " ( ", -1), "'", " ' ", -1), "\"", " \" ", -1), " ")
	r := make([]string, 0)
	for _, e := range ss {
		if e != " " && e != "" {
			r = append(r, e)
		}
	}
	return r
}

func Parse(s *[]string, allowblock bool) Expr {
	if len(*s) == 0 {
		return Error{"Unexpected EOF."}
	}
	t := (*s)[0]
	*s = (*s)[1:]
	if t == "'" {
		if allowblock {
			return EvalBlock{Parse(s, false)}
		} else {
			return Parse(s, allowblock)
		}
	}
	if t == "\"" {
		ss := (*s)[0]
		*s = (*s)[1:]
		if ss == "\"" {
			return String("")
		}
		if len(*s) == 0 {
			return Error{"Missing end quote"}
		}
		for (*s)[0] != "\"" {
			if len(*s) == 0 {
				return Error{"Missing end quote"}
			}
			ss = strings.Join([]string{ss, (*s)[0]}, " ")
			*s = (*s)[1:]
		}
		(*s) = (*s)[1:]
		return String(ss)
	}
	if t == "(" {
		l := make(ExprList, 0)
		for (*s)[0] != ")" {
			l = append(l, Parse(s, allowblock))
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
		return env.find(unwrapSymbol(e))[unwrapSymbol(e)]
	} else if eb, ok := e.(EvalBlock); ok {
		return eb.e
	} else if v, ok := e.(String); ok {
		return v
	} else if !bool(list_(env, e).(Boolean)) {
		return e
	} else if el := e.(ExprList); bool(symbol_(env, el[0]).(Boolean)) {
		if s0 := unwrapSymbol(e.(ExprList)[0]); s0 == "quote" {
			return el[1]
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
			if l, ok := el[1].(ExprList); ok {
				for i, v := range l {
					if v == Symbol(".") {
						if i != len(l)-2 {
							return Error{"Multiple variables after '.' not allowed!"}
						}
						//append(...) removes the . from the list of params
						return UserProc{true, append(l[:i], l[i+1:]...), el[2]}
					}
				}
				return UserProc{false, el[1].(ExprList), el[2]}
			} else if v, ok := el[1].(Symbol); ok {
				return UserProc{true, ExprList{v}, el[2]}
			}
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
	if strings.HasPrefix(s, "#\\") {
		return decodeCharacter(s[2:])
	}
	if i, err := strconv.Atoi(s); err == nil {
		return Number((i))
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return Number(f)
	}
	return Symbol(s)
}
