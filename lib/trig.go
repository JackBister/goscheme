package goscheme

import (
	"math"
	"math/cmplx"
)

func acos(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		if v2, ok2 := args[0].(Complex); ok2 {
			return Complex(cmplx.Acos(complex128(v2)))
		}
		return Error{"acos: Argument 1 is not a number."}
	} else {
		return Number(math.Acos(float64(v)))
	}
}

func asin(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		if v2, ok2 := args[0].(Complex); ok2 {
			return Complex(cmplx.Asin(complex128(v2)))
		}
		return Error{"asin: Argument 1 is not a number."}
	} else {
		return Number(math.Asin(float64(v)))
	}
}

func atan(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		if v2, ok2 := args[0].(Complex); ok2 {
			return Complex(cmplx.Atan(complex128(v2)))
		}
		return Error{"atan: Argument 1 is not a number."}
	} else {
		return Number(math.Atan(float64(v)))
	}
}

func cos(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		if v2, ok2 := args[0].(Complex); ok2 {
			return Complex(cmplx.Cos(complex128(v2)))
		}
		return Error{"cos: Argument 1 is not a number."}
	} else {
		return Number(math.Cos(float64(v)))
	}
}

func sin(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		if v2, ok2 := args[0].(Complex); ok2 {
			return Complex(cmplx.Sin(complex128(v2)))
		}
		return Error{"sin: Argument 1 is not a number."}
	} else {
		return Number(math.Sin(float64(v)))
	}
}

func tan(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		if v2, ok2 := args[0].(Complex); ok2 {
			return Complex(cmplx.Tan(complex128(v2)))
		}
		return Error{"tan: Argument 1 is not a number."}
	} else {
		return Number(math.Tan(float64(v)))
	}
}
