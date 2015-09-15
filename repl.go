// +build !windows

package main

import (
	"github.com/carmark/pseudo-terminal-go/terminal"
	"os"
)

var t *terminal.Terminal

var replFuncs = map[string]func(){
	":q":    func() { t.ReleaseFromStdInOut(); os.Exit(0) },
	":quit": func() { t.ReleaseFromStdInOut(); os.Exit(0) },
}

func replStart() {
	t, _ = terminal.NewWithStdInOut()
}

func readLine() string {
	in, _ := t.ReadLine()
	return in
}
