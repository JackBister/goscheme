package main

import (
	"bufio"
	"os"
)

var reader *bufio.Reader

func exit() {
	os.Exit(0)
}

func replStart() {
	reader = bufio.NewReader(os.Stdin)
}

func readLine() string {
	in, _ := reader.ReadString('\n')
	return in
}
