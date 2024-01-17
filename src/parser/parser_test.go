package parser

import (
	"testing"
	"dux/src/ast"
	"dux/src/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) < 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d instead", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]

		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" { 
		t.Errorf("statement.TokenLiteral() not 'let'. got=%q instead", statement.TokenLiteral())
		return false
	}

	letStmt, ok := statement.(*ast.LetStatement)

	// case not a LetStatement
	if !ok {
		t.Errorf("statement is not *ast.LetStatement. got=%T instead", statement)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s intead", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}
