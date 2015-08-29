package goscheme

func remainder(e Environment, args ...Expr) Expr {
	if _, ok := args[0].(Number); !ok {
		return Error{"modulo: Argument 1 is not a number"}
	}
	if _, ok := args[1].(Number); !ok {
		return Error{"modulo: Argument 2 is not a number"}
	}
	//Yuck!
	a0n := int64(float64(args[0].(Number)))
	a1n := int64(float64(args[1].(Number)))
	return Number(a0n % a1n)
}

