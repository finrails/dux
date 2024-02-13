package repl

import (
	"bufio"
	"dux/evaluator"
	"dux/lexer"
	"dux/object"
	"dux/parser"
	"fmt"
	"io"
)

const ARROW = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(ARROW)
		scanned := scanner.Scan()

		if !scanned { return }

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		evaluated := evaluator.Eval(program, env)

		if evaluated == nil { continue }
		
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
