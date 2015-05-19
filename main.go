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
	schemec.GlobalEnv = schemec.Environment{schemec.StandardEnv(), nil}
	readLoop()
}

func readLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>")
		in, _ := reader.ReadString('\n')
		in = strings.Trim(in, " \r\n")
		eval(in)
	}
}

func eval(s string) {
	t := schemec.Tokenize(s)
	p := schemec.Parse(&t)
	r := schemec.Eval(p, schemec.GlobalEnv)
	if f, ok := r.(schemec.Number); ok {
		fs := strconv.FormatFloat(float64(f), 'f', -1, 64)
		r = schemec.Symbol(fs)
	}
	if _, ok := r.(schemec.Error); ok {
		fmt.Print("Error: ")
	}
	fmt.Println(r)
}
