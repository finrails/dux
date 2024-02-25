package main

import (
	"dux/evaluator"
	"dux/lexer"
	"dux/object"
	"dux/parser"
	"dux/repl"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		user, err := user.Current()

		if err != nil {
			panic(err)
		}

		fmt.Printf("Welcome %s. Dux Language Interpreter!\n", user.Username)
		fmt.Printf("You can evalute Dux commands here\n");

		repl.Start(os.Stdin, os.Stdout)
	} else {
		absp, err := filepath.Abs(args[0])
		if err != nil { fmt.Println("Error:", err); return }

		content, err := os.ReadFile(absp)
		if err != nil { fmt.Println("Error:", err); return }

		l := lexer.New(string(content))
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()

		evaluated := evaluator.Eval(program, env)

		fmt.Print(evaluated.Inspect())
	}
}
