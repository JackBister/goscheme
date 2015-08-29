package goscheme

import (
	"math"
)

func modulo(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"modulo: Argument 1 is not a number"}
	}
	if _, ok := args[1].(Number); !ok {
		return Error{"modulo: Argument 2 is not a number"}
	}
	a0n := float64(args[0].(Number))
	a1n := float64(args[1].(Number))
	return Number(a0n - math.Floor(a0n/a1n)*a1n)
}

func remainder(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"remainder: Argument 1 is not a number"}
	}
	if _, ok := args[1].(Number); !ok {
		return Error{"remainder: Argument 2 is not a number"}
	}
	a0n := float64(args[0].(Number))
	a1n := float64(args[1].(Number))
	return Number(math.Mod(a0n, a1n))
}

