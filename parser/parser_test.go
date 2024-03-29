package parser

import (
	"dux/ast"
	"dux/lexer"
	"fmt"
	"testing"
)

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not hash literal. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash obj has wrong number of pairs. want=%d, got=%d", 3, len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("no test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: "ok", false: "not ok"}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 2 {
		t.Fatalf("hash exp has wrong number os pairs. want=%d, got=%d", 2, len(hash.Pairs))
	}

	expected := map[bool]string{
		true: "ok",
		false: "not ok",
	}

	for key, value := range hash.Pairs {
		bol, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("hash exp's key is not *ast.Boolean. got=%T", key)
		}

		expectedValue := expected[bol.Value]

		testString(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: "one", 2: "two", 3: "three"}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash exp has wrong number of pairs. want=%d, got=%d", 3, len(hash.Pairs))
	}

	expected := map[int64]string{
		1: "one",
		2: "two",
		3: "three",
	}

	for key, value := range hash.Pairs {
		il, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("hash key is not ast.IntegerLiteral. got=%T", key)
		}

		expectedValue := expected[il.Value]

		testString(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash exp has wrong number of pairs. want=%d, got=%d", 0, len(hash.Pairs))
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 { 
		t.Fatalf("hash exp has wrong number of pairs. want=%d, got=%d", 3, len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatal("exp not *ast.ExpressionStatement")
	}

	ie, ok := stmt.Expression.(*ast.IndexExpresssion)
	if !ok {
		t.Fatalf("stmt.exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, ie.Left, "myArray") { return }

	if !testInfixExpression(t, ie.Index, 1, "+", 1) { return }
}

func TestParsingArrayLiterals(t *testing.T) {
	tests := []struct{
		input            string
		expectedLen      int
		expectedElements []interface{}
		left             interface{}
		right            interface{}
		op               string
	}{
		{input: "[1, 2, 3]", expectedLen: 3, expectedElements: []interface{}{1, 2, 3}},
		{input: `["foo", 1, 2, "bar"]`, expectedLen: 4, expectedElements: []interface{}{"foo", 1, 2, "bar"}},
		{input: "[1, 2]", expectedLen: 2, expectedElements: []interface{}{1, 2}},
		{input: "[]", expectedLen: 0, expectedElements: []interface{}{}},
		{input: "[1, 2, 6 + 2]", expectedLen: 3, expectedElements: []interface{}{1, 2, 8}, left: 6, right: 2, op: "+"},
		{input: `[3 * 3, "foobar"]`, expectedLen: 2, expectedElements: []interface{}{9, "foobar"}, left: 3, right: 3, op: "*"},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if len(program.Statements) != 1 {
			t.Fatalf("program has wrong number of statements. got=%d, want=%d", len(program.Statements), 1)
		}

		if !ok {
			t.Fatalf("exp not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		array, ok := stmt.Expression.(*ast.ArrayLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not an *ast.ArrayLiteral. got=%T", stmt.Expression)
		}

		if len(array.Elements) != tc.expectedLen {
			t.Fatalf("array has wrong number of elements. want=%d, got=%d", tc.expectedLen, len(array.Elements))
		}

		testArrayElements(t, array.Elements, tc.expectedElements, tc.left, tc.right, tc.op)
	}
}

func testArrayElements(t *testing.T, elems []ast.Expression, expectedElements []interface{}, left, right interface{}, op string) {
	var caseNum int
	for index, exp := range elems {
		caseNum++

		switch exp.(type) {
		case *ast.InfixExpression:
			if !testInfixExpression(t, exp, left, op, right) {
				t.Errorf("test case (%d): array[%d] infix exp has wrong element. want=%T", caseNum, index, expectedElements[index])
			}
		default:
			if !testLiteralExpression(t, exp, expectedElements[index]) {
				t.Errorf("test case (%d): array[%d] was wrong element. got=%T, want=%T", caseNum, index, exp.String(), expectedElements[index])
			}
		}
	}

}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello, world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello, world" {
		t.Errorf("literal.Value not %q. got=%q", "hello, world", literal.Value)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct{
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements length should be %d. got=%d", 1, len(program.Statements))
		}

		stmt := program.Statements[0]

		if !testLetStatement(t, stmt, tc.expectedIdentifier) { return }

		exp := stmt.(*ast.LetStatement).Value

		if !testLiteralExpression(t, exp, tc.expectedValue) { return }
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct{
		input         string
		expectedValue interface{}
	}{
		{input: "return 5;", expectedValue: 5},
		{input: "return foobar;", expectedValue: "foobar"},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements length should be %d. got=%d", 1, len(program.Statements))
		}

		rex, ok := program.Statements[0].(*ast.ReturnStatement)

		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ReturnStatement. got=%T", program.Statements[0])
		}

		if !testLiteralExpression(t, rex.ReturnValue, tc.expectedValue) { return }
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
		{"2 * (5 + 5)", "(2 * (5 + 5))"},
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
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"add(2, add(5, 5))", "add(2, add(5, 5))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
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

func TestIfExpressions(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements should have %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifok, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, ifok.Condition, "x", "<", "y") {
		return
	}

	if len(ifok.Consequence.Statements) != 1 {
		t.Errorf("Consequence.Statements should be %d. got=%d", 1, len(ifok.Consequence.Statements))
	}

	consequence, ok := ifok.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Consequence.Statements[0] not *ast.ExpressionStatement. got=%T", ifok.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") { return }

	if ifok.Alternative != nil {
		t.Errorf("IfExpression.Alternative not nil. got=%T", ifok.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements should be %d. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifok, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("program.Expression not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, ifok.Condition, "x", "<", "y") { return }

	if len(ifok.Consequence.Statements) != 1 {
		t.Fatalf("Consequence.Statements length shoul be %d. got=%d", 1, len(ifok.Consequence.Statements))
	}

	consequence, ok := ifok.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Consequence.Statements[0] not *ast.ExpressionStatement. got=%T", ifok.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") { return }

	if len(ifok.Alternative.Statements) != 1 {
		t.Fatalf("Alternative.Statements[0] should be %d. got=%d", 1, ifok.Alternative.Statements)
	}

	alternative, ok := ifok.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Alternative.Statements[0] not *ast.ExpressionStatement. got=%T", ifok.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") { return }
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements length should be %d. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fl, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(fl.Parameters) != 2 {
		t.Fatalf("function literal parameters length, should be %d. got=%d", 2, len(fl.Parameters))
	}

	testLiteralExpression(t, fl.Parameters[0], "x")
	testLiteralExpression(t, fl.Parameters[1], "y")

	if len(fl.Body.Statements) != 1 {
		t.Fatalf("function literal Body.Statements should be %d. got=%d", 1, len(fl.Body.Statements))
	}

	bodyStmt, ok := fl.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body statement not *ast.ExpressionStatement. got=%T", fl.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct{
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements length should be %d. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		fel, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.FunctionLiteral. got=%T", stmt.Expression)
		}

		if len(fel.Parameters) != len(tc.expectedParams) {
			t.Errorf("length of parameters wrong. want %d, got=%d", len(tc.expectedParams), len(fel.Parameters))
		}

		for i, id := range tc.expectedParams {
			testLiteralExpression(t, fel.Parameters[i], id)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	
	checkParserErrors(t, p)

	if len(program.Statements) > 1 {
		t.Fatalf("program.Statements length should be %d. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	cex, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, cex.Function, "add") { return }

	if len(cex.Arguments) != 3 {
		t.Fatalf("wrong lenth of arguments. want=%d, got=%d", 3, len(cex.Arguments))
	}

	testLiteralExpression(t, cex.Arguments[0], 1)
	testInfixExpression(t, cex.Arguments[1], 2, "*", 3)
	testInfixExpression(t, cex.Arguments[2], 4, "+", 5)
}

func TestCallExpressionArgumentParsing(t *testing.T) {
	tests := []struct{
		input string
		expectedIdent string
		expectedArgs []string
	}{
		{
			input: "add();",
			expectedIdent: "add",
			expectedArgs: []string{},
		},
		{
			input: "add(1);",
			expectedIdent: "add",
			expectedArgs: []string{"1"},
		},
		{
			input: "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs: []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tc := range tests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements length should be %d. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		cex, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.CallExpresion. got=%T", stmt.Expression)
		}

		if !testIdentifier(t, cex.Function, tc.expectedIdent) { return }

		if len(cex.Arguments) != len(tc.expectedArgs) {
			t.Fatalf("cex.Arguments length should be %d. got=%d", len(tc.expectedArgs), len(cex.Arguments))
		}

		for i, arg := range tc.expectedArgs {
			if cex.Arguments[i].String() != arg {
				t.Errorf("argument wrong. want=%q, got=%q", arg, cex.Arguments[i].String())
			}
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
		if str, ok := exp.(*ast.StringLiteral); ok {
			return testString(t, str, v)
		}
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testString(t *testing.T, exp ast.Expression, expected string) bool {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.StringLiteral. got=%T", exp)
		return false
	}

	if str.Value != expected {
		t.Errorf("string has invalid value. got=%q, want=%q", str.Value, expected)
		return false
	}

	return true
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
