package main

import (
	"fmt"
	"os/user"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome %s. Dux Language Interpreter!\n", user.Username)
	fmt.Printf("You can evalute Dux commands here\n");
}
