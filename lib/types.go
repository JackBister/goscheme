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
	//A list of symbols that the arguments to the function will be bound to.
	params ExprList
	body   Expr
}

func (u UserProc) isExpr() {}

func (u UserProc) eval(e Environment, args ...Expr) Expr {
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
