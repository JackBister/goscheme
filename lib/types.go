package goscheme

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

type Expr interface {
	isExpr()
}

type Environment struct {
	Local       map[string]Expr
	LocalSyntax map[Symbol]transformer
	Parent      *Environment
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
	nsm := map[Symbol]transformer{}
	for k, v := range e.LocalSyntax {
		nsm[k] = v
	}
	return Environment{nm, nsm, e.Parent}
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
func (s String) String() string {
	return "\"" + string(s) + "\""
}

type Boolean bool

func (b Boolean) isExpr() {}
func (b Boolean) String() string {
	if bool(b) {
		return "#t"
	}
	return "#f"
}

type Byte byte

func (b Byte) isExpr() {}
func (b Byte) String() string {
	return fmt.Sprintf("0x%x", byte(b))
}

type Character rune

func (c Character) isExpr() {}
func (c Character) String() string {
	return fmt.Sprintf("#\\%c", c)
}

func decodeCharacter(s string) Character {
	charMap := map[string]Character{
		"nul":       Character(0x00),
		"soh":       Character(0x01),
		"stx":       Character(0x02),
		"etx":       Character(0x03),
		"eot":       Character(0x04),
		"enq":       Character(0x05),
		"ack":       Character(0x06),
		"alarm":     Character(0x07),
		"bel":       Character(0x07),
		"backspace": Character(0x08),
		"bs":        Character(0x08),
		"tab":       Character(0x09),
		"ht":        Character(0x09),
		"linefeed":  Character(0x0A),
		"newline":   Character(0x0A),
		"lf":        Character(0x0A),
		"vtab":      Character(0x0B),
		"vt":        Character(0x0B),
		"page":      Character(0x0C),
		"ff":        Character(0x0C),
		"return":    Character(0x0D),
		"cr":        Character(0x0D),
		"so":        Character(0x0E),
		"si":        Character(0x0F),
		"dle":       Character(0x10),
		"dc1":       Character(0x11),
		"dc2":       Character(0x12),
		"dc3":       Character(0x13),
		"dc4":       Character(0x14),
		"nak":       Character(0x15),
		"syn":       Character(0x16),
		"etb":       Character(0x17),
		"can":       Character(0x18),
		"em":        Character(0x19),
		"sub":       Character(0x1A),
		"escape":    Character(0x1B),
		"esc":       Character(0x1B),
		"fs":        Character(0x1C),
		"gs":        Character(0x1D),
		"rs":        Character(0x1E),
		"us":        Character(0x1F),
		"space":     Character(0x20),
		"sp":        Character(0x20),
		"delete":    Character(0x7F),
		"del":       Character(0x7F),
	}
	if utf8.RuneCountInString(s) == 1 {
		r, _ := utf8.DecodeRuneInString(s)
		return Character(r)
	}
	return charMap[s]
}

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

func listtovec(e Environment, args ...Expr) Expr {
	l, ok := args[0].(ExprList)
	if !ok {
		return Error{"list->vector: Argument 1 is not a list."}
	}
	return Vector([]Expr(l))
}

/*
A port can either be written to or read from.
bufio.Reader/Writer is used for some conveniences, but ports also must be closeable.
The assumption is that for input-ports w and wclose will be nil, and vice versa.
If r != nil then rclose must also be != nil, same with w and wclose.
While R5RS seems to only deal with files this implementation is vague enough that
adding network support shouldn't be a big deal.
*/
type Port struct {
	rclose io.Closer
	wclose io.Closer
	r      *bufio.Reader
	w      *bufio.Writer
}

func (p Port) isExpr() {}

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
	//The environment the closure was created in
	env Environment
	//A list of symbols that the arguments to the function will be bound to.
	params ExprList
	body   Expr
}

func (u UserProc) isExpr() {}

func (u UserProc) String() string {
	return "<closure>"
}

func (u UserProc) eval(e Environment, args ...Expr) Expr {
	for k, v := range u.env.Local {
		if e.Local[k] == nil {
			e.Local[k] = v
		}
	}
	if len(args) < len(u.params) {
		if u.variadic {
			if len(args) != len(u.params)-1 {
				return Error{"Too few arguments (need " + strconv.Itoa(len(u.params)-1) + ")"}
			}
		} else {
			return Error{"Too few arguments (need " + strconv.Itoa(len(u.params)) + ")"}
		}
	}
	if len(args) > len(u.params) && !u.variadic {
		return Error{"Too many arguments (need " + strconv.Itoa(len(u.params)) + ")"}
	}
	for i, par := range u.params {
		if i == len(u.params)-1 && u.variadic {
			e.Local[unwrapSymbol(par)] = ExprList(args[i:])
		} else {
			e.Local[unwrapSymbol(par)] = args[i]
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
func (e Error) isExpr()       {}

type Pattern []Expr

//returns a map of bindings from symbol in the pattern to expression in the
//input expressionlist and a bool which is true if the match was successful
func (p Pattern) match(keywords map[Symbol]bool, env map[Symbol]Expr, el []Expr) (map[Symbol]Expr, bool) {
	pel := []Expr(p)
	if len(pel) < len(el) {
		if s, ok := pel[len(pel)-1].(Symbol); !ok || unwrapSymbol(s) != "..." {
			return env, false
		}
	} else if len(pel) > len(el) {
		return env, false
	}
	for i, e := range pel {
		if pelel, ok := e.(ExprList); ok {
			pp := Pattern([]Expr(pelel))
			if elel, ok2 := el[i].(ExprList); !ok2 {
				return env, false
			} else {
				if nEnv, ok3 := pp.match(keywords, env, []Expr(elel)); ok3 {
					for k, v := range nEnv {
						//TODO: Slow, copies entire env even though most of it could already be correct.
						//nEnv should only contain newly introduced bindings
						if env[k] != nil && !bool(eqv(Environment{}, env[k], v).(Boolean)) {
							return env, false
						}
						env[k] = v
					}
				} else {
					return env, false
				}
			}

		} else if ps, ok2 := e.(Symbol); ok2 {
			if keywords[ps] {
				els, ok3 := el[i].(Symbol)
				if !ok3 || els != ps {
					return env, false
				}
			}
			if unwrapSymbol(ps) == "..." && i == len(pel)-1 {
				env[ps] = ExprList(el[len(pel)-1:])
				continue
			}
			if env[ps] != nil && !bool(eqv(Environment{}, env[ps], el[i]).(Boolean)) {
				return env, false
			}
			env[ps] = el[i]
		} else {
			return env, false
		}
	}
	return env, true
}

type transformer interface {
	transform([]Expr) Expr
}

type SyntaxRule struct {
	patterns     []Pattern
	replacements []Expr
	//ech, didn't want to do a linear array search for each symbol in the input
	//keywords[s] will be true if s is a keyword
	keywords map[Symbol]bool
}

func (s SyntaxRule) isExpr() {}
func (s SyntaxRule) transform(el []Expr) Expr {
	for i, p := range s.patterns {
		if e, ok := p.match(s.keywords, map[Symbol]Expr{}, el); !ok {
			continue
		} else {
			ret, ok2 := replace(s.replacements[i], e)
			if !ok2 {
				//TODO: More info
				return Error{"syntax-rule: Error encountered while transforming input."}
			}
			return ret
		}
	}
	return Error{"Syntax-rule: Input does not match any pattern."}
}

func replace(e Expr, m map[Symbol]Expr) (Expr, bool) {
	if s, ok := e.(Symbol); ok {
		if m[s] != nil {
			return m[s], true
		}
		//TODO:
		return s, true
	}
	el, ok := e.(ExprList)
	if !ok {
		return e, false
	}
	ret := make([]Expr, len(el))
	for i, exp := range el {
		if exps, ok := exp.(Symbol); ok && m[exps] != nil {
			//TODO: If ... is used anywhere else than at the end of a list?
			if unwrapSymbol(exps) == "..." && i == len(el)-1 {
				rem, ok2 := m[exps].(ExprList)
				if !ok2 {
					return ExprList(ret), false
				}
				//need to make slice smaller to avoid <nil> after the append
				ret = ret[:len(ret)-1]
				ret = append(ret, []Expr(rem)...)
				continue
			}
			ret[i] = m[exps]
		} else if expel, ok := exp.(ExprList); ok {
			rexpel, ok2 := replace(expel, m)
			ret[i] = rexpel
			if !ok2 {
				return ExprList(ret), false
			}
		} else {
			ret[i] = exp
		}
	}
	return ExprList(ret), true
}
