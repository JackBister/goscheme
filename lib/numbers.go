package goscheme

import (
	"math"
)

func ceiling(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"ceiling: Argument 1 is not a number"}
	} else {
		return Number(math.Ceil(float64(v)))
	}

}

func floor(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"floor: Argument 1 is not a number"}
	} else {
		return Number(math.Floor(float64(v)))
	}
}

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

func sqrt(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"sqrt: Argument 1 is not a number"}
	} else {
		return Number(math.Sqrt(float64(v)))
	}
}

func truncate(e Environment, args ...Expr) Expr {
	if v, ok := args[0].(Number); !ok {
		return Error{"truncate: Argument 1 is not a number"}
	} else {
		return Number(math.Trunc(float64(v)))
	}
}
