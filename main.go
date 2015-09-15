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
	"flag"
	"fmt"
	"github.com/jackbister/goscheme/lib"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	maxp := flag.Int("cores", runtime.NumCPU(), "Sets the number of CPU cores that the interpreter may use. If not given, all available cores will be used.")
	interactive := flag.Bool("i", false, "Enters interactive mode after executing the given files. If no files are given this is the default.")
	flag.Parse()
	runtime.GOMAXPROCS(*maxp)
	goscheme.GlobalEnv = goscheme.StandardEnv()
	for _, a := range flag.Args() {
		eval("(load (quote " + a + "))")
	}
	if *interactive || len(flag.Args()) == 0 {
		readLoop()
	}
}

func readLoop() {
	replStart()
	for {
		fmt.Print(">>")
		in := readLine()
		in = strings.Trim(in, " \r\n")
		if replFuncs[in] != nil {
			replFuncs[in]()
		} else {
			eval(in)
		}
	}
}

func eval(s string) {
	t := goscheme.Tokenize(s)
	p := goscheme.Parse(&t, true)
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
