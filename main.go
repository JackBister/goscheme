/*
    goscheme - a Lisp interpreter in Go
    Copyright (C) 2015 Jack Bister

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
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
