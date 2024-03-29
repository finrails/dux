package evaluator

import (
	"dux/lexer"
	"dux/object"
	"dux/parser"
	"testing"
)

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct{
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
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

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
		{
			"one": 10 - 9,
			two: 1 + 1,
			"thr" + "ee": 6/2,
			4: 4,
			true: 5,
			false: 6
		}
	`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return object.Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[uint64]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for the given key")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestReduceImplementation(t *testing.T) {
	input := `
		let reduce = fn(arr, initial, f) {
			let iter = fn(arr, result) {
				if (len(arr) == 0) {
					result
				} else {
					iter(tail(arr), f(result, first(arr)));
				}
			};

			iter(arr, initial)
		}

		let sum = fn(arr) {
			reduce(arr, 0, fn(initial, current) { initial + current });
		}

		sum([2, 4, 6])
	`

	result := testEval(input)

	num, ok := result.(*object.Integer)
	if !ok {
		t.Fatalf("object has wrong type. should be Integer, got=%s", result.Type())
	}

	if num.Value != 12 {
		t.Errorf("num has wrong value. should be 12, got=%d", num.Value)
	}
}

func TestMapImplementation(t *testing.T) {
	input := `
		let map = fn(arr, f) {
			let iter = fn(arr, accumulated) {
				if (len(arr) == 0) {
					accumulated
				} else {
					iter(tail(arr), push(accumulated, f(first(arr))))
				}
			};

			iter(arr, []);
		}

		map([1, 2, 3], fn(x) { x * 2 })
	`

	result := testEval(input)

	array, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("object should be Array. got=%T", result)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. want=3, got=%d", len(array.Elements))
	}

	if array.Inspect() != "[2, 4, 6]" {
		t.Errorf("array has wrong elements. got=%s, want=%s", array.Inspect(), "[2, 4, 6]")
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct{
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("Hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`first([1, 2, 3])`, 1},
		{`first(["poo", true])`, "poo"},
		{`let slice = [6, 2]; first(slice)`, 6},
		{`last([1, 2, 3])`, 3},
		{`let slice = [12, 8]; last(slice)`, 8},
		{`first("foobar")`, "f"},
		{`last("gap")`, "p"},
		{`tail([1, 2, 3])`, "[2, 3]"},
		{`tail([2, 3])`, "[3]"},
		{`tail([3])`, "[]"},
		{`tail([])`, nil},
		{`head([1, 2, 3])`, "[1, 2]"},
		{`head([1, 2])`, "[1]"},
		{`head([1])`, "[]"},
		{`head([])`, "nil"},
		{`push([], 5)`, "[5]"},
		{`push([1], true)`, "[1, true]"},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)

		switch expected := tc.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			switch obj := evaluated.(type) {
			case *object.String:
				if obj.Value != expected { t.Fatalf("String has invalid char sequence. want=%q, got=%q", tc.expected, obj.Value) }
				continue
			case *object.Error:
				if obj.Message != expected { t.Fatalf("wrong error message. got=%q", obj.Message) }
				continue
			case *object.Array:
				if obj.Inspect() != expected { t.Fatalf("array has wrong elements. want=%s, got=%s", tc.expected, obj.Inspect()) }
			}
		}
	}
}

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

func TestErrorHandling(t *testing.T) {
	tests := []struct{
		input    string
		expected string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 2;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 1", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (10 > 1) { return true + false } return 1 }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"foo" - "bar"`, "unknown operator: STRING - STRING"},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)

		errorObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("should return an object.Error. got %T(%v+)", evaluated, evaluated)
			continue
		}

		if errorObj.Message != tc.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tc.expected, errorObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct{
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5, let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object should be *object.Function. got=%T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T", evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("String has wrong value. want=%q, got=%q", "Hello, World!", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object not *object.String. got=%T", evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. want=%q, got=%q", "Hello, World!", str.Value)
	}
}

func TestStringMultiplication(t *testing.T) {
	input := `"foo" * 5`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object not a String. got=%T", evaluated)
	}

	if str.Value != "foofoofoofoofoo" {
		t.Errorf("String has wrong value. expected=%q, got=%q", "foofoofoofoofoo", str.Value)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct{
		input    string
		expected int64
	}{
		{"let identify = fn(x) { x; }; identify(5);", 5},
		{"let identify = fn(x) { return x; }; identify(5);", 5},
		{"let double  = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		let newAdder = fn(x) { fn(y) { x + y }; };
		let addTwo = newAdder(2);
		addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("object is not an Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d, want=%d", len(result.Elements), 3)
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct{
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i]", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6,
		},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
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
