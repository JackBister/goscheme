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
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var GlobalEnv Environment

func Tokenize(s string) []string {
	ignoreString := "\t\n\r"

	afterComment := false
	inQuotes := false
	var b bytes.Buffer
	for i, r := range s {
		if r == '"' {
			inQuotes = !inQuotes
			b.WriteRune(r)
			continue
		}
		if inQuotes {
			b.WriteRune(r)
			continue
		}
		if r == ';' {
			if i != 0 && s[i-1] == '\\' {
				b.WriteRune(r)
				continue
			}
			afterComment = true
			continue
		}
		if r == '\n' && afterComment {
			afterComment = false
		}
		if afterComment {
			continue
		}
		if strings.ContainsRune(ignoreString, r) {
			continue
		}
		if r == '\'' || r == '(' || r == ')' {
			b.WriteString(" " + string(r) + " ")
			continue
		}
		b.WriteRune(r)
	}
	ss := strings.Split(b.String(), " ")
	r := make([]string, 0)
	for i := 0; i < len(ss); i++ {
		if strings.HasPrefix(ss[i], "\"") {
			toJoin := make([]string, 0)
			for j := 0; i+j < len(ss); j++ {
				toJoin = append(toJoin, ss[i+j])
				if strings.HasSuffix(ss[i+j], "\"") {
					break
				}
			}
			r = append(r, strings.Join(toJoin, " "))
			i += len(toJoin) - 1
			continue
		}
		if ss[i] != " " && ss[i] != "" {
			r = append(r, ss[i])
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
	if t == "#" && (*s)[0] == "(" {
		p := Parse(s, allowblock)
		if e, ok := p.(Error); ok {
			return e
		}
		//wew
		return vector(Environment{}, ExprListToSlice(p.(ExprList))...)
	}
	if strings.HasPrefix(t, "\"") {
		if t[len(t)-1] != '"' {
			return Error{"Missing end quote"}
		}
		return String(t[1 : len(t)-1])
	}
	if t == "(" {
		l := make([]Expr, 0)
		for (*s)[0] != ")" {
			l = append(l, Parse(s, allowblock))
			if len(*s) == 0 {
				return Error{"Missing ')'"}
			}
		}
		*s = (*s)[1:]
		return SliceToExprList(l)
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
	} else if _, ok := e.(ExprList); !ok {
		return e
	} else if el := ExprListToSlice(e.(ExprList)); len(el) != 0 && bool(symbol_(env, el[0]).(Boolean)) {
		if s0 := unwrapSymbol(*e.(ExprList).car); s0 == "quote" {
			if len(el) != 2 {
				return Error{"quote: Must be of form '(quote <datum>)'"}
			}
			return el[1]
		} else if s0 == "syntax-rules" {
			if len(el) != 3 {
				return Error{"syntax-rules: Must be of form '(syntax-rules (<keywords>) ((<pattern>) (<template>) ... (<pattern>) (<template>)))'."}
			}
			sr := SyntaxRule{keywords: map[Symbol]bool{}}
			kw, ok := el[1].(ExprList)
			if !ok {
				return Error{"syntax-rules: Must be of form '(syntax-rules (<keywords>) ((<pattern>) (<template>) ... (<pattern>) (<template>)))'."}
			}
			for _, exp := range ExprListToSlice(kw) {
				exps, ok := exp.(Symbol)
				if !ok {
					return Error{"syntax-rules: All keywords must be symbols."}
				}
				sr.keywords[exps] = true
			}
			tr, ok := el[2].(ExprList)
			if !ok {
				return Error{"syntax-rules: Must be of form '(syntax-rules (<keywords>) ((<pattern>) (<template>) ... (<pattern>) (<template>)))'."}
			}
			trl := ExprListToSlice(tr)
			nextIsPattern := true
			for _, exp := range trl {
				if nextIsPattern {
					expl, ok := exp.(ExprList)
					if !ok {
						return Error{"syntax-rules: Must be of form '(syntax-rules (<keywords>) ((<pattern>) (<template>) ... (<pattern>) (<template>)))'."}
					}
					sr.patterns = append(sr.patterns, Pattern(expl))
					nextIsPattern = false
				} else {
					sr.replacements = append(sr.replacements, exp)
					nextIsPattern = true
				}
			}
			if !nextIsPattern {
				return Error{"syntax-rules: More patterns than templates given."}
			}
			return sr
		} else if s0 == "if" {
			if len(el) < 3 || len(el) > 4 {
				return Error{"if: Must be of form '(if <test> <consequent> <alternate>)' or '(if <test> <consequent>)'"}
			}
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
			if len(el) != 3 {
				return Error{"define: Must be of form '(define <variable> <expression>)'"}
			}
			er := Eval(el[2], env)
			env.Local[unwrapSymbol(el[1])] = er
			if _, ok := er.(Proc); !ok {
				return er
			}
		} else if s0 == "set!" {
			if len(el) != 3 {
				return Error{"set!: Must be of form '(set! <variable> <expression>)'"}
			}
			if env.find(unwrapSymbol(el[1])) != nil {
				er := Eval(el[2], env)
				env.find(unwrapSymbol(el[1]))[unwrapSymbol(el[1])] = er
				return er
			}
			//TODO: Error?
		} else if s0 == "define-syntax" {
			if len(el) != 3 {
				return Error{"define-syntax: Must be of form '(define-syntax <name> <syntax transformer>)'."}
			}
			if _, ok := el[1].(Symbol); !ok {
				return Error{"define-syntax: Must be of form '(define-syntax <name> <syntax transformer>)'."}
			}
			if t, ok := Eval(el[2], env).(transformer); !ok {
				return Error{"define-syntax: Must be of form '(define-syntax <name> <syntax transformer>)'."}
			} else {
				env.LocalSyntax[el[1].(Symbol)] = t
			}
		} else if s0 == "lambda" {
			if len(el) != 3 {
				return Error{"lambda: Must be of form '(lambda <formals> <body>)'"}
			}
			//newenv := env.copy()
			newenv := Environment{map[string]Expr{}, map[Symbol]transformer{}, &env}
			if l, ok := el[1].(ExprList); ok {
				expl := ExprListToSlice(l)
				for i, v := range expl {
					if v == Symbol(".") {
						if i != l.Length()-2 {
							return Error{"Multiple variables after '.' not allowed!"}
						}
						//append(...) removes the . from the list of params
						return UserProc{true, newenv, SliceToExprList(append(expl[:i], expl[i+1:]...)), []Expr{}, el[2]}
					}
				}
				return UserProc{false, newenv, el[1].(ExprList), []Expr{}, el[2]}
			} else if v, ok := el[1].(Symbol); ok {
				return UserProc{true, newenv, SliceToExprList([]Expr{v}), []Expr{}, el[2]}
			}
		} else if s0 == "go" {
			if len(el) != 2 {
				return Error{"go: Must be of form '(go <expression>)'"}
			}
			c := make(Channel)
			go func(c chan Expr) {
				c <- Eval(el[1], env)
			}(c)
			return c
		} else if s0 == "time" {
			if len(el) != 2 {
				return Error{"time: Must be of form '(time <expression>)'"}
			}
			t := time.Now()
			ret := Eval(el[1], env)
			fmt.Println("time:", time.Now().Sub(t))
			return ret
		} else if env.LocalSyntax[Symbol(s0)] != nil {
			exp := env.LocalSyntax[Symbol(s0)].transform(el)
			if _, ok := exp.(Error); ok {
				return exp
			}
			return Eval(exp, env)
		} else {
			//Kinda sucks having to repeat this here and below. Problem is the s0 cast is too convenient.
			proc := Eval(el[0], env)
			if procp, ok2 := proc.(Proc); ok2 {
				elcopy := make([]Expr, len(el))
				copy(elcopy, el)
				elcopy[0] = procp
				return Eval(SliceToExprList(elcopy), env)
			} else {
				//TODO
				fmt.Println("Error: Expected procedure, have", proc)
			}
		}
	} else if p, ok := el[0].(Proc); ok {
		args := make([]Expr, 0, len(el[1:]))
		for _, arg := range el[1:] {
			args = append(args, Eval(arg, env))
		}
		nEnv := Environment{map[string]Expr{}, map[Symbol]transformer{}, &env}
		return p.eval(nEnv, args...)
	} else {
		proc := Eval(el[0], env)
		if procp, ok2 := proc.(Proc); ok2 {
			elcopy := make([]Expr, len(el))
			copy(elcopy, el)
			elcopy[0] = procp
			return Eval(SliceToExprList(elcopy), env)
		} else {
			//TODO
			fmt.Println("Error: Expected procedure, have", proc)
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
