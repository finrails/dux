package evaluator

import (
	"dux/lexer"
	"dux/object"
	"dux/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"2 - 2 - 2", -2},
		{"2 * -2", -4},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct{
		input    string
		expected bool
	}{
		{"true == true", true},
		{"true", true},
		{"false", false},
		{"1 > 2", false},
		{"1 < 2", true},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == false", false},
		{"false == false", true},
		{"false == true", false},
		{"true != false", true},
		{"true != true", false},
		{"false != true", true},
		{"false != false", false},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

func TestExclamationOperator(t *testing.T) {
	tests := []struct{
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct{
		input string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		integer, ok := tc.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct{
		input    string
		expected int64
	}{
		{"if (10 > 1) { if (5 > 1) { return 10 } return 1 }", 10},
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"9; return 10", 10},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object value should be %d. got=%d", expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object value should be %t. got=%t", expected, result.Value)
		return false
	}

	return true
}

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj != NIL {
		t.Errorf("object is not NIL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}