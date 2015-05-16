package main
import (
	"bufio"
	"fmt"
	"github.com/jackbister/schemec/lib"
	"os"
	"strconv"
	"strings"
)

var env schemec.Environment

func main() {
	env = schemec.Environment{schemec.StandardEnv(), nil}
	readLoop()
}

func readLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>")
		in, _ := reader.ReadString('\n')
		if len(in) > 0 {
			eval(in[:len(in)-1])
		}
	}
}

func eval(s string) {
	t := schemec.Tokenize(s)
	p := schemec.Parse(&t)
	r := schemec.Eval(p, env)
	if f, ok := r.(schemec.Number); ok {
		fs := strconv.FormatFloat(float64(f), 'f', -1, 64)
		doti := strings.Index(s, ".")
		if doti > 0 {
			fs = strings.Trim(fs[doti:], "0")
		}
		r = schemec.Symbol(fs)
	}
	fmt.Println(r)
}
