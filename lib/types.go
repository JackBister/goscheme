package goscheme

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/cmplx"
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

func (e Environment) isExpr() {}

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

type Complex complex128

func (c Complex) isExpr() {}

func complex_(e Environment, args ...Expr) Expr {
	_, ok := args[0].(Complex)
	return Boolean(ok)
}

func makepolar(e Environment, args ...Expr) Expr {
	f1, ok := args[0].(Number)
	if !ok {
		return Error{"make-polar: Argument 1 is not a number."}
	}
	f2, ok2 := args[1].(Number)
	if !ok2 {
		return Error{"make-polar: Argument 2 is not a number."}
	}
	//For some reason "Rect" returns the polar form?
	return Complex(cmplx.Rect(float64(f1), float64(f2)))
}

func makerect(e Environment, args ...Expr) Expr {
	f1, ok := args[0].(Number)
	if !ok {
		return Error{"make-rectangular: Argument 1 is not a number."}
	}
	f2, ok2 := args[1].(Number)
	if !ok2 {
		return Error{"make-rectangular: Argument 2 is not a number."}
	}
	return Complex(complex(float64(f1), float64(f2)))
}

func angle(e Environment, args ...Expr) Expr {
	c, ok := args[0].(Complex)
	if !ok {
		return Error{"angle: Argument 1 is not a complex number."}
	}
	_, theta := cmplx.Polar(complex128(c))
	return Number(theta)
}

func imagpart(e Environment, args ...Expr) Expr {
	c, ok := args[0].(Complex)
	if !ok {
		return Error{"imag-part: Argument 1 is not a complex number."}
	}
	return Number(imag(complex128(c)))
}

func magnitude(e Environment, args ...Expr) Expr {
	c, ok := args[0].(Complex)
	if !ok {
		return Error{"magnitude: Argument 1 is not a complex number."}
	}
	r, _ := cmplx.Polar(complex128(c))
	return Number(r)
}

func realpart(e Environment, args ...Expr) Expr {
	c, ok := args[0].(Complex)
	if !ok {
		return Error{"real-part: Argument 1 is not a complex number."}
	}
	return Number(real(complex128(c)))
}

//An EvalBlock wraps an expression and delays evaluation of the expr.
//Primarily(only?) used for actions involving apostrophes.
type EvalBlock struct {
	e Expr
}

func (e EvalBlock) isExpr() {}

type ExprList struct {
	car *Expr
	cdr *ExprList
}

func (el ExprList) isExpr() {}

func (el ExprList) Length() int {
	it := &el
	length := 0
	for {
		if it == nil || it.car == nil {
			break
		}
		length++
		it = it.cdr
	}
	return length
}

func (el ExprList) String() string {
	var b bytes.Buffer
	b.WriteString("(")
	it := &el
	for {
		if it == nil || it.car == nil {
			break
		}
		fmt.Fprint(&b, *it.car)
		if it.cdr != nil && it.cdr.car != nil {
			b.WriteString(" ")
		}
		it = it.cdr
	}
	b.WriteString(")")
	return b.String()
}

func ExprListToSlice(el ExprList) []Expr {
	ret := make([]Expr, 0, el.Length())
	it := &el
	for {
		if it.car == nil || it == nil {
			break
		}
		ret = append(ret, *it.car)
		it = it.cdr
	}
	return ret
}

func SliceToExprList(el []Expr) ExprList {
	if len(el) == 0 {
		return ExprList{nil, nil}
	}
	car := el[0]
	cdr := SliceToExprList(el[1:])
	return ExprList{&car, &cdr}
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
	isExpr()
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
	if len(args) < u.params.Length() {
		if u.variadic {
			if len(args) != u.params.Length()-1 {
				return Error{"Too few arguments (need " + strconv.Itoa(u.params.Length()-1) + ")"}
			}
		} else {
			return Error{"Too few arguments (need " + strconv.Itoa(u.params.Length()) + ")"}
		}
	}
	if len(args) > u.params.Length() && !u.variadic {
		return Error{"Too many arguments (need " + strconv.Itoa(u.params.Length()) + ")"}
	}
	for i, par := range ExprListToSlice(u.params) {
		if i == u.params.Length()-1 && u.variadic {
			e.Local[unwrapSymbol(par)] = SliceToExprList(args[i:])
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

// type Pattern []Expr

type Match struct {
	direct   Expr
	ellipsis []Expr
}

type Pattern Expr

func replace(curr, last Expr, m map[Symbol]Match) (Expr, bool) {
	if s, ok := curr.(Symbol); ok {
		if m[s].direct != nil {
			return m[s].direct, true
		}
		if string(s) == "..." {
			lastSymbol, ok2 := last.(Symbol)
			if !ok2 {
				return nil, false
			}
			return Vector(m[lastSymbol].ellipsis), true
		}
		return curr, true
	}
	var elSlice []Expr
	var isExprList, isVector bool
	if _, isExprList = curr.(ExprList); isExprList {
		elSlice = ExprListToSlice(curr.(ExprList))
	} else {
		_, isVector = curr.(Vector)
		if isVector {
			elSlice = []Expr(curr.(Vector))
		}
	}
	if isExprList || isVector {
		var lastExpr Expr
		rel := make([]Expr, 0, len(elSlice))
		for _, v := range elSlice {
			e, ok2 := replace(v, lastExpr, m)
			if !ok2 {
				return nil, false
			}
			lastExpr = v
			//if v is an ellipsis the result will be a vector that should be expanded in place instead of replacing the variable
			if s, ok2 := v.(Symbol); ok2 && string(s) == "..." {
				rel = append(rel, []Expr(e.(Vector))...)
			} else {
				rel = append(rel, e)
			}
		}
		if isExprList {
			return SliceToExprList(rel), true
		}
		return Vector(rel), true
	}
	//TODO:
	return curr, true
}

/*
	literals: if literals[k] == true, k is a literal
	env: The matched pattern variables so far, used for recursion. Should be returned at the end of the function.
	e: The expression to match against the pattern
	ellipsis: Whether the matching is being done for an ellipsis at the end of a list/vector.

	Returns a map with bindings from symbols in patterns to expressions in the input and whether the pattern matched the input.
*/
func match(p Pattern, literals map[Symbol]bool, env map[Symbol]Match, e Expr, ellipsis bool) (map[Symbol]Match, bool) {
	var elSlice, pelSlice []Expr
	var isExprList, isVector bool

	if _, isExprList = p.(ExprList); isExprList {
		if _, ok := e.(ExprList); !ok {
			return env, false
		}
		elSlice = ExprListToSlice(e.(ExprList))
		pelSlice = ExprListToSlice(p.(ExprList))
	}
	if _, isVector = p.(Vector); isVector {
		if _, ok := e.(Vector); !ok {
			return env, false
		}
		elSlice = []Expr(e.(Vector))
		pelSlice = []Expr(p.(Vector))
	}

	if !isExprList && !isVector {
		if ps, ok2 := p.(Symbol); ok2 {
			es, ok3 := e.(Symbol)
			if literals[ps] || (ok3 && literals[es]) {
				//The pattern or input is a literal
				return env, ps == es
			}
			//The pattern is a pattern variable
			pv, exists := env[ps]
			if !exists {
				env[ps] = Match{}
			}

			if ellipsis {
				pv.ellipsis = append(env[ps].ellipsis, e)
			} else {
				pv.direct = e
			}
			env[ps] = pv
			return env, true
		}
		//TODO: else if vector, else return (equal? p i)
		return env, bool(eq(Environment{}, p, e).(Boolean))
	} else {
		if len(elSlice) < len(pelSlice) {
			if ps, ok3 := pelSlice[len(pelSlice)-1].(Symbol); !ok3 || string(ps) != "..." {
				//There is no case where the input matches the pattern while also being shorter than the pattern (unless there is an ellipsis that matches nothing)
				//From here on we can be a bit reckless with indexing into elSlice since we know it's at least as big as pelSlice, so any valid index in pelSlice is valid in elSlice.
				return env, false
			}
		}
		if ps, ok3 := pelSlice[len(pelSlice)-2].(Symbol); ok3 {
			//"P is of the form (P1 P2 ... Pn . Px) and F is a list or improper list of n or more elements whose first n elements match P1 through Pn and whose nth cdr matches Px,"
			if string(ps) == "." {
				//n = len(pelSlice)-2
				for i, v := range pelSlice[:len(pelSlice)-2] {
					nEnv, match := match(v.(Pattern), literals, env, elSlice[i], false)
					if !match {
						return nEnv, false
					}
					env = nEnv
				}
				//matches "the nth cdr" (the elements following n) in the input with the element after the dot
				return match(pelSlice[len(pelSlice)-1].(Pattern), literals, env, SliceToExprList(elSlice[len(pelSlice)-1:]), false)
			}
		}
		if ps, ok3 := pelSlice[len(pelSlice)-1].(Symbol); ok3 {
			//"P is of the form (P1 ... Pn Px ...) and F is a proper list of n or more elements whose first n elements match P1 through Pn and whose remaining elements each match Px,"
			if string(ps) == "..." {
				//n = len(pelSlice)-2
				for i, v := range pelSlice[:len(pelSlice)-1] {
					nEnv, match := match(v.(Pattern), literals, env, elSlice[i], false)
					if !match {
						return nEnv, false
					}
					env = nEnv
				}
				//Elements after n
				for _, v := range elSlice[len(pelSlice)-1:] {
					nEnv, match := match(pelSlice[len(pelSlice)-2].(Pattern), literals, env, v, true)
					if !match {
						return nEnv, false
					}
					env = nEnv
				}
				return env, true
			}
		}
		if len(pelSlice) == len(elSlice) {
			//"P is of the form (P1 ... Pn) and F is a list of n elements that match P1 through Pn,"
			for i, v := range pelSlice {
				nEnv, match := match(v.(Pattern), literals, env, elSlice[i], ellipsis)
				if !match {
					return nEnv, false
				}
				env = nEnv
			}
			return env, true
		}
	}
	return env, false
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
		if e, ok := match(p, s.keywords, map[Symbol]Match{}, SliceToExprList(el), false); !ok {
			continue
		} else {
			ret, ok2 := replace(s.replacements[i], nil, e)
			if !ok2 {
				//TODO: More info
				return Error{"syntax-rule: Error encountered while transforming input."}
			}
			return ret
		}
	}
	return Error{"syntax-rule: Input does not match any pattern."}
}
