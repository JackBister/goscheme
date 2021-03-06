// +build !windows

package main

import (
	"github.com/jackbister/goscheme/lib/terminal"
	"os"
)

var t *terminal.Terminal

func exit() {
	t.ReleaseFromStdInOut()
	os.Exit(0)
}

func replStart() {
	t, _ = terminal.NewWithStdInOut()
}

func readLine() string {
	in, _ := t.ReadLine()
	return in
}
