package parser

import (
	"dux/src/ast"
	"dux/src/lexer"
	"fmt"
	"testing"
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
	checkParserErrors(t, p)

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

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d instead", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("statement not *ast.ReturnStatement. got=%T instead", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() not 'return', got=%q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, should have 1. got=%d instead", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement type. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	
	l := lexer.New(input)
	p := New(l) // p stands for parser not for program

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, should have %d. instead got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T instead", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not *ast.IntegerLiteral. got=%T instead", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got %d instead", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLitera not %s. got=%s instead", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct{
		input        string
		operator     string
		value        interface{}
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tc := range prefixTests {
		l := lexer.New(tc.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contains %d statements. got=%d instead", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T instead", program.Statements[0])
		}

		expression, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T instead", stmt.Expression)
		}

		if expression.Operator != tc.operator {
			t.Fatalf("expression.Operator is not '%s'. got='%s'", tc.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tc.value) { return }
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct{
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foo + bar;", "foo", "+", "bar"},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}

	for _, tc := range infixTests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expression is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, expression.Left, tc.left) { return }

		if expression.Operator != tc.operator {
			t.Fatalf("expression.Operator is not '%s'. got='%s' instead", tc.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tc.right) { return }
	}

}

func TestLiteralInfixExpressions(t *testing.T) {
	tests := []struct{
		input    string
		left		 interface{}
		operator string
		right		 interface{}
	}{
		{
			input: "5 * 5;",
			left: 5,
			operator: "*",
			right: 5,
		},
		{
			input: "2 + 4;",
			left: 2,
			operator: "+",
			right: 4,
		},
		{
			input: "8 / 16;",
			left: 8,
			operator: "/",
			right: 16,
		},
		{
			input: "32 - 16;",
			left: 32,
			operator: "-",
			right: 16,
		},
		{
			input: "foo + bar",
			left: "foo",
			operator: "+",
			right: "bar",
		},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements slice size should be %d. got=%d", 1, len(program.Statements))
			return
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tc.left, tc.operator, tc.right) {
			t.Errorf("stmt.Expression is not a valid InfixExpression")
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct{
		input    string
		expected string
	}{
		{"3 * 6 + 9 + 12 + 15 / 3 * 6 * 9 + 12", "(((((3 * 6) + 9) + 12) + (((15 / 3) * 6) * 9)) + 12)"},
		{"2 + 4 + 6 * 8 / 10 + 12 * 14", "(((2 + 4) + ((6 * 8) / 10)) + (12 * 14))"},
		{"2 + 4 * 6 + 9", "((2 + (4 * 6)) + 9)"},
		{"a + b * c", "(a + (b * c))"},
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"2 * 4 + 6", "((2 * 4) + 6)"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		current := program.String()

		if current != tc.expected {
			t.Errorf("expected=%q, got=%q", tc.expected, current)
		}
	}
}

func TestBooleanExpressions(t *testing.T) {
	tests := []struct{
		input    string
		expected string
	}{
		{input: "true;", expected: "true"},
		{input: "false;", expected: "false"},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements should have %d statements. got=%d", 1, len(program.Statements))
			return
		}

		stmt, _ := program.Statements[0].(*ast.ExpressionStatement)

		be, ok := stmt.Expression.(*ast.Boolean)

		if !ok {
			t.Errorf("stmt.Expression not a ast.Boolean. got=%T", stmt.Expression)
		}

		if be.Token.Literal != tc.expected {
			t.Errorf("*Boolean.Token.Literal should be %q. got=%q", tc.expected, be.Token.Literal)
		}

		if be.String() != tc.expected {
			t.Errorf("*Boolean.String() shoud be %q. got=%q", tc.expected, be.String())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 { return }

	t.Errorf("parser has %d errors", len(errors))
	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" { 
		t.Errorf("statement.TokenLiteral() not 'let'. got=%q instead", statement.TokenLiteral())
		return false
	}

	letStmt, ok := statement.(*ast.LetStatement)

	// case statement not a LetStatement
	if !ok {
		t.Errorf("statement is not *ast.LetStatement. got=%T instead", statement)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s intead", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name#TokenLiteral not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	ilok, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T instead", il)
		return false
	}

	if ilok.Value != value {
		t.Errorf("ilok.Value not %d. got=%d instead", value, ilok.Value)
		return false
	}

	if ilok.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("ilok.TokenLiteral not %d. got=%s", value, il.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	
	if !ok {
		t.Errorf("exp not not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	bok, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("ast.Expression not *ast.Boolean. got=%d", exp)
		return false
	}

	if bok.Value != value {
		t.Errorf("*ast.Boolean.Value not %t. got=%t", value, bok.Value)
		return false
	}

	if bok.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("*ast.Boolean.TokenLiteral() not '%t'. got='%s'", value, bok.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	ifExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not *ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, ifExp.Left, left) {
		t.Errorf("exp.Left not a Literal. got=%T", ifExp.Left)
		return false
	}

	if ifExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, ifExp.Operator)
		return false
	}

	if !testLiteralExpression(t, ifExp.Right, right) {
		t.Errorf("exp.Left not a Literal. got=%T", ifExp.Right)
		return false
	}

	return true
}
