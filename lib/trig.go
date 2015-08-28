package goscheme

import (
	"math"
)

func acos(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"acos: Argument 1 is not a number."}
	} else {
		return Number(math.Acos(float64(v)))
	}
}

func asin(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"asin: Argument 1 is not a number."}
	} else {
		return Number(math.Asin(float64(v)))
	}
}

func atan(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"atan: Argument 1 is not a number."}
	} else {
		return Number(math.Atan(float64(v)))
	}
}

func cos(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"cos: Argument 1 is not a number."}
	} else {
		return Number(math.Cos(float64(v)))
	}
}

func sin(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"sin: Argument 1 is not a number."}
	} else {
		return Number(math.Sin(float64(v)))
	}
}

func tan(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"tan: Argument 1 is not a number."}
	} else {
		return Number(math.Tan(float64(v)))
	}
}
