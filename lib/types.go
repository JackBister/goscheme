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

/*
	TypedExpr represents an expression which has a type.
	The endgame is that Expr could be removed entirely because all expressions will be typed.
*/
type TypedExpr interface {
	getType(e Environment) Type
}

type Environment struct {
	Local       map[string]Expr
	LocalSyntax map[Symbol]transformer
	//TODO: Types could be put in the regular Local map.
	LocalTypes map[Symbol]Type
	Parent     *Environment
}

func (e Environment) getType(env Environment) Type {
	return env.findType(Symbol("Environment"))
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

/*
	findType looks up a symbol first in the LocalTypes map, then recurses on the parent.
*/
func (e *Environment) findType(s Symbol) Type {
	if t, ok := e.LocalTypes[s]; ok {
		return t
	}
	if e.Parent == nil {
		return e.findType("Undefined")
	}
	return e.Parent.findType(s)
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
	ntm := map[Symbol]Type{}
	for k, v := range e.LocalTypes {
		ntm[k] = v
	}
	return Environment{nm, nsm, ntm, e.Parent}
}

/*
	A Type is attached to an expression.
	A Type can either be "just a type", or a composition of other types.
	If a Type is a function, the parameters field will contain a slice of the types of each parameter
	and the returnType field will hold a pointer to the type of the returned value of the function.
	If a Type is a list, the parameters field will contain a slice of the types of each element in the list.

	TODO:
	returnType isn't used anywhere.
	Vector isn't typed.
*/
type Type struct {
	name   string
	isFunc bool
	isList bool
	//Parameters when isFunc == true, list content when isList == true
	parameters []Type
	returnType *Type
}

/*
	The type of a Type value is the Type "Type". Wew.
*/
func (t Type) getType(env Environment) Type {
	return env.findType(Symbol("Type"))
}
func (t Type) isExpr() {}

/*
	Compares two types.
	Two types are equal if:
		(1) Both types are functions and their slices of parameters contains the same types and their returnType is the same.
		(2) Both types are lists and their slices contain the same types.
		(3) Their names are equal.
*/
func (t Type) isEqual(rhs Type) bool {
	if t.isFunc {
		if !rhs.isFunc {
			return false
		}
		if len(t.parameters) != len(rhs.parameters) {
			return false
		}
		for i, p := range t.parameters {
			if !p.isEqual(rhs.parameters[i]) {
				return false
			}
		}
		fmt.Println(t.parameters, rhs.parameters)
		fmt.Println(t.returnType, rhs.returnType)
		if t.returnType != nil {
			if rhs.returnType == nil || !t.returnType.isEqual(*rhs.returnType) {
				return false
			}
		} else if rhs.returnType != nil {
			return false
		}
		return true
	}
	if t.isList {
		if !rhs.isList {
			return false
		}
		if len(t.parameters) != len(rhs.parameters) {
			return false
		}
		for i, p := range t.parameters {
			if !p.isEqual(rhs.parameters[i]) {
				return false
			}
		}
		return true
	}
	return t.name == rhs.name
}

/*
	Returns a string representation of the type. In the future this may be more fleshed out.
*/
func (t Type) String() string {
	return t.name
}

/*
	An ExprAndExpectedType is the result of an expression of the form <Expr> : <Type-Expr>.
	e represents <Expr>, t represents <Type-Expr>
*/
type ExprAndExpectedType struct {
	e Expr
	t Expr
}

func (e ExprAndExpectedType) isExpr() {}

/*
	getType does a limited evaluation of t in the environment env to determine the type of the expression, and then returns that type.

	TODO: Should it just be a straight up Eval call?
*/
func (e ExprAndExpectedType) getType(env Environment) Type {
	if ts, ok := e.t.(Symbol); ok {
		return env.findType(ts)
	}
	if tt, ok := e.t.(Type); ok {
		return tt
	}
	//TODO: Unsure about this. What if the type-expr has side effects?
	if res, ok := Eval(e.t, env).(Type); ok {
		return res
	}
	if tel, ok := e.t.(ExprList); ok {
		tels := ExprListToSlice(tel)
		//If the type is a function, it looks like ((ParametersType) => ReturnType)
		if len(tels) == 3 {
			if maybeArrow, ok := tels[1].(Symbol); ok && string(maybeArrow) == "=>" {	
				return parseFuncType(tels, env)
			}
		}
		ret := Type{"", false, true, []Type{}, nil}
		var b bytes.Buffer
		b.WriteString("(")
		ret, s, ok := getSliceType(tels, ret, &env)
		if !ok {
			return env.findType(Symbol("Undefined"))
		}
		b.WriteString(s + ")")
		ret.name = b.String()
		return ret
	}
	return env.findType(Symbol("Undefined"))
}

func parseFuncType(el []Expr, env Environment) Type {
	var b bytes.Buffer
	b.WriteString("(")
	ret := Type{"", true, false, []Type{}, nil}
	if parTypeList, ok := el[0].(ExprList); !ok {
		return env.findType(Symbol("Undefined"))
	} else {
		parTypeListType := ExprAndExpectedType{nil, parTypeList}.getType(env)
		ret.parameters = parTypeListType.parameters
		b.WriteString(parTypeListType.String() + " => ")
	}
	//This code isn't normal but on Go it is
	retType := ExprAndExpectedType{nil, el[2]}.getType(env)
	ret.returnType = &retType
	b.WriteString(retType.String() + ")")
	ret.name = b.String()
	return ret
}

/*
	getSliceType is used for retrieving the type of a Type-Expr when it is a list in ExprAndExpectedType.
*/
func getSliceType(typeExprSlice []Expr, ret Type, env *Environment) (Type, string, bool) {
	var b bytes.Buffer
	for i, expr := range typeExprSlice {
		if es, ok := expr.(Symbol); ok {
			t := env.findType(es)
			ret.parameters = append(ret.parameters, t)
			b.WriteString(t.String())
		} else if el, ok := expr.(ExprList); ok {
			b.WriteString("(")
			listType := Type{"", false, true, []Type{}, nil}
			listType, elString, ok := getSliceType(ExprListToSlice(el), listType, env)
			if !ok {
				return ret, "", false
			}
			ret.parameters = append(ret.parameters, listType)
			b.WriteString(elString + ")")
		} else {
			return ret, "", false
		}
		if i != len(typeExprSlice)-1 {
			b.WriteString(" ")
		}
	}
	ret.name = b.String()
	return ret, b.String(), true
}

/*
Number type
The language does not distinguish between precise (int) and imprecise (float)
values right now. 64 bit float should be sufficiently accurate for large ints.
*/
type Number float64

func (n Number) getType(env Environment) Type {
	return env.findType(Symbol("Number"))
}
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

func (s Symbol) getType(env Environment) Type {
	return env.findType(Symbol("Symbol"))
}
func (s Symbol) isExpr() {}
func unwrapSymbol(s Expr) string {
	return string(s.(Symbol))
}

type String string

func (s String) getType(env Environment) Type {
	return env.findType(Symbol("String"))
}
func (s String) isExpr() {}
func (s String) String() string {
	return "\"" + string(s) + "\""
}

type Boolean bool

func (b Boolean) getType(env Environment) Type {
	return env.findType(Symbol("Boolean"))
}
func (b Boolean) isExpr() {}
func (b Boolean) String() string {
	if bool(b) {
		return "#t"
	}
	return "#f"
}

type Byte byte

func (b Byte) getType(env Environment) Type {
	return env.findType(Symbol("Byte"))
}
func (b Byte) isExpr() {}
func (b Byte) String() string {
	return fmt.Sprintf("0x%x", byte(b))
}

type Character rune

func (c Character) getType(env Environment) Type {
	return env.findType(Symbol("Character"))
}
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

func (c Channel) getType(env Environment) Type {
	return env.findType(Symbol("Channel"))
}
func (c Channel) isExpr() {}

type Complex complex128

func (c Complex) getType(env Environment) Type {
	return env.findType(Symbol("Complex"))
}
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

/*
	Returns the type of a list. Unfortunately requires iterating through the whole list to get types of all elements.

	TODO:
	All ExprList operations could be optimized with a cached slice in each ExprList.
*/
func (el ExprList) getType(env Environment) Type {
	ret := Type{"", false, true, []Type{}, nil}
	var b bytes.Buffer
	b.WriteString("(")
	els := ExprListToSlice(el)
	for i, expr := range els {
		if et, ok := expr.(TypedExpr); ok {
			t := et.getType(env)
			ret.parameters = append(ret.parameters, t)
			b.WriteString(t.String())
		} else {
			ret.parameters = append(ret.parameters, Type{"Undefined", false, false, []Type{}, nil})
			b.WriteString("Undefined")
		}
		if i != len(els)-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString(")")
	ret.name = b.String()
	return ret
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

func (p Port) getType(env Environment) Type {
	return env.findType(Symbol("Port"))
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
	//A list of arguments which have already been passed (for currying)
	//e.g. if a 3-parameter function is called with 2 parameters partialArgs will
	//contain those args and the UserProc will be returned so the last arg can be fulfilled
	partialArgs []Expr
	body        Expr
}

/*
	TODO: Everything
*/
func (u UserProc) getType(env Environment) Type {
	var b bytes.Buffer
	b.WriteString("((")
	parList := ExprListToSlice(u.params)
	
	parTypes := make([]Type, 0, len(parList))
	for i := len(u.partialArgs); i < len(parList); i++ {
		par := parList[i]
		if part, ok := par.(ExprAndExpectedType); ok {
			t := part.getType(env)
			parTypes = append(parTypes, t)
			b.WriteString(t.String())
		} else {
			b.WriteString("Undefined")
			parTypes = append(parTypes, env.findType(Symbol("Undefined")))
		}
		if i != len(parList)-1 {
			b.WriteString(" ")
		}
	}
	//TODO:
	b.WriteString(") => Undefined)")
	retType := env.findType(Symbol("Undefined"))
	ret := Type{b.String(), true, false, parTypes, &retType}
	return ret
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
	if len(args)+len(u.partialArgs) < u.params.Length() {
		if !u.variadic || len(args) != u.params.Length()-1 {
			for _, arg := range args {
				u.partialArgs = append(u.partialArgs, arg)
			}
			return u
		}
	}
	if len(args)+len(u.partialArgs) > u.params.Length() && !u.variadic {
		return Error{"Too many arguments (need " + strconv.Itoa(u.params.Length()) + ")."}
	}
	for i, par := range ExprListToSlice(u.params) {
		isPartialArg := i < len(u.partialArgs)
		part, isTyped := par.(ExprAndExpectedType)
		var s string
		if isTyped {
			s = unwrapSymbol(part.e)
		} else {
			s = unwrapSymbol(par)
		}
		if i == u.params.Length()-1 && u.variadic {
			if isPartialArg {
				e.Local[s] = SliceToExprList(u.partialArgs[i:])
			} else {
				e.Local[s] = SliceToExprList(args[i-len(u.partialArgs):])
			}
		} else {
			if isPartialArg {
				e.Local[s] = u.partialArgs[i]
			} else {
				e.Local[s] = args[i-len(u.partialArgs)]
			}
		}
		if isTyped {
			parameterType := part.getType(e)
			if argt, ok := e.Local[s].(TypedExpr); ok {
				argType := argt.getType(e)
				if parameterType.name == "Undefined" {
					//If the type name is "Undefined" we don't want to modify the type map
					if parts, ok := part.t.(Symbol); ok && string(parts) != "Undefined" {
						e.LocalTypes[parts] = argType.getType(e)
					}
				} else if parameterType.isList {
					partel := part.t.(ExprList)
					partels := ExprListToSlice(partel)
					if !typedUserProcListParameter(partels, argType, &e) {
						return Error{"Parameter " + s + " expected type " + part.getType(e).String() + ", have " + argType.String() + "."}
					}
				} else if !argType.isEqual(parameterType) {
					return Error{"Parameter " + s + " expected type " + part.getType(e).String() + ", have " + argType.String() + "."}
				}
			} else {
				return Error{"Parameter " + s + " expected type " + part.getType(e).String() + " but has no type."}
			}
		}
	}
	return Eval(u.body, e)
}

/*
	typedUserProcListParameter is used to do type checking for UserProcs. This needed to be broken out of UserProc.eval for recursion purposes.
*/
func typedUserProcListParameter(parameterTypedExprSlice []Expr, argType Type, env *Environment) bool {
	if !argType.isList || len(parameterTypedExprSlice) != len(argType.parameters) {
		return false
	}
	ret := true
	for i := range parameterTypedExprSlice {
		if parSymbol, ok := parameterTypedExprSlice[i].(Symbol); ok {
			newParType := env.findType(parSymbol)
			if argType.parameters[i].isEqual(newParType) {
				continue
			}
			if newParType.name == "Undefined" {
				env.LocalTypes[parSymbol] = argType.parameters[i]
				continue
			}
			ret = false
		} else if parExprList, ok := parameterTypedExprSlice[i].(ExprList); ok {
			if !typedUserProcListParameter(ExprListToSlice(parExprList), argType.parameters[i], env) {
				ret = false
			}
		} else {
			ret = false
		}
	}
	//We want to wait with returning so that the whole list can become as defined as possible.
	//Otherwise some types will be listed as undefined even though they would be defined if the loop continued.
	return ret
}

type BuiltIn struct {
	//The name of the function. Used for error printouts.
	name string
	//if maxParams is -1, the function is variadic.
	minParams, maxParams int
	//The go function this struct represents.
	fn func(Environment, ...Expr) Expr
	//Same as UserProc.partialArgs
	partialArgs []Expr
}

func (b BuiltIn) isExpr() {}

func (b BuiltIn) eval(e Environment, args ...Expr) Expr {
	if len(args)+len(b.partialArgs) < b.minParams {
		for _, arg := range args {
			b.partialArgs = append(b.partialArgs, arg)
		}
		return b
	}
	if len(args)+len(b.partialArgs) > b.maxParams && b.maxParams != -1 {
		return Error{b.name + ": Too many arguments (max " + strconv.Itoa(b.maxParams) + ")"}
	}
	return b.fn(e, append(b.partialArgs, args...)...)
}

//NewBuiltIn exists to maintain one interface for creating new built ins even if the struct layout changes.
//Maybe a pimpl style thing would work in the future?
func NewBuiltIn(name string, minParams, maxParams int, fn func(Environment, ...Expr) Expr) BuiltIn {
	return BuiltIn{name, minParams, maxParams, fn, []Expr{}}
}

type Error struct {
	s string
}

func (e Error) Error() string { return e.s }
func (e Error) getType(env Environment) Type {
	return env.findType(Symbol("Error"))
}
func (e Error) isExpr() {}

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

func (s SyntaxRule) getType(env Environment) Type {
	return env.findType(Symbol("SyntaxRule"))
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
