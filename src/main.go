package main

import (
	"fmt"
	"os/user"
	"os"
	"dux/src/repl"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome %s. Dux Language Interpreter!\n", user.Username)
	fmt.Printf("You can evalute Dux commands here\n");

	repl.Start(os.Stdin, os.Stdout)
}
