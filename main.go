package main
import (
	"bufio"
	"fmt"
	"github.com/jackbister/goscheme/lib"
	"os"
	"strconv"
	"strings"
)

func main() {
	goscheme.GlobalEnv = goscheme.Environment{goscheme.StandardEnv(), nil}
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
	t := goscheme.Tokenize(s)
	p := goscheme.Parse(&t)
	r := goscheme.Eval(p, goscheme.GlobalEnv)
	if f, ok := r.(goscheme.Number); ok {
		fs := strconv.FormatFloat(float64(f), 'f', -1, 64)
		r = goscheme.Symbol(fs)
	}
	if _, ok := r.(goscheme.Error); ok {
		fmt.Print("Error: ")
	}
	if s, ok := r.(goscheme.Symbol); !ok || string(s) != "" {
		fmt.Println(r)
	}
}
