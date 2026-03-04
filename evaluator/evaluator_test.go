package evaluator

import (
	"modern-basic/lexer"
	"modern-basic/object"
	"modern-basic/parser"
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
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestIfElseIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", int64(10)},
		{"if (false) { 10 }", nil},
		{"if (false) { 10 } elseif (true) { 20 } else { 30 }", int64(20)},
		{"if (false) { 10 } elseif (false) { 20 } else { 30 }", int64(30)},
		{"if (1) { 10 }", int64(10)},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestWhileLoopEvaluation(t *testing.T) {
	input := `let i = 0; while (i < 5) { i = i + 1; }; i;`
	testIntegerObject(t, testEval(t, input), 5)
}

func TestForLoopEvaluation(t *testing.T) {
	input := `let sum = 0; for (let i = 0; i < 5; i = i + 1) { sum = sum + i; }; sum;`
	testIntegerObject(t, testEval(t, input), 10)
}

func TestFunctionObject(t *testing.T) {
	evaluated := testEval(t, `fn(x) { x + 2; }`)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not *object.Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not x. got=%q", fn.Parameters[0])
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"function add(x, y) { return x + y; } add(3, 7);", 10},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestFunctionLocalScopeIsolation(t *testing.T) {
	input := `
let x = 10;
let f = fn() { let x = 99; return x; };
f();
x;
`

	testIntegerObject(t, testEval(t, input), 10)
}

func TestReturnPropagationInFunction(t *testing.T) {
	input := `
let f = fn() {
  if (true) {
    return 10;
  }
  return 1;
};
f();
`

	testIntegerObject(t, testEval(t, input), 10)
}

func TestNestedReturnFromLoopInFunction(t *testing.T) {
	input := `
let f = fn() {
  for (let i = 0; i < 5; i = i + 1) {
    if (i == 3) {
      return i;
    }
  }
  return 99;
};
f();
`

	testIntegerObject(t, testEval(t, input), 3)
}

func TestFunctionWrongArgumentCount(t *testing.T) {
	evaluated := testEval(t, `let add = fn(x, y) { x + y; }; add(1);`)
	testErrorObject(t, evaluated, "wrong number of arguments: got=1, want=2")
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{"\"hello\" - \"world\"", "unknown operator: STRING - STRING"},
		{"if (5 + true) { 10 }", "type mismatch: INTEGER + BOOLEAN"},
		{"while (5 + true) { 10 }", "type mismatch: INTEGER + BOOLEAN"},
		{"let add = fn(x) { x; }; add(1, 2)", "wrong number of arguments: got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testErrorObject(t, evaluated, tt.expectedMessage)
	}
}

func TestStringLiteralEvaluation(t *testing.T) {
	evaluated := testEval(t, `"hello world"`)
	strObj, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not *object.String. got=%T (%+v)", evaluated, evaluated)
	}

	if strObj.Value != "hello world" {
		t.Fatalf("String value wrong. expected=%q, got=%q", "hello world", strObj.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	evaluated := testEval(t, `"hello" + " " + "world"`)
	strObj, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not *object.String. got=%T (%+v)", evaluated, evaluated)
	}

	if strObj.Value != "hello world" {
		t.Fatalf("String value wrong. expected=%q, got=%q", "hello world", strObj.Value)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(t, input)

	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not *object.Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", int64(1)},
		{"[1, 2, 3][1]", int64(2)},
		{"[1, 2, 3][2]", int64(3)},
		{"let i = 0; [1][i];", int64(1)},
		{"[1, 2, 3][1 + 1];", int64(3)},
		{"let myArray = [1, 2, 3]; myArray[2];", int64(3)},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", int64(6)},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", int64(2)},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int64)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
{
    "one": 10 - 9,
    two: 1 + 1,
    "thr" + "ee": 6 / 2,
    4: 4,
    true: 5,
    false: 6
}`

	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(&object.Boolean{Value: true}).HashKey():   5,
		(&object.Boolean{Value: false}).HashKey():  6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("hash has wrong num pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Fatalf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, int64(5)},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, int64(5)},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, int64(5)},
		{`{true: 5}[true]`, int64(5)},
		{`{false: 5}[false]`, int64(5)},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int64)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestIndexOperatorErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{`{"name": "Cart"}[fn(x) { x }]`, "unusable as hash key: FUNCTION"},
		{`[1, 2]["x"]`, "index operator not supported: ARRAY[STRING]"},
		{`1[0]`, "index operator not supported: INTEGER[INTEGER]"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testErrorObject(t, evaluated, tt.expectedMessage)
	}
}

func testEval(t *testing.T, input string) object.Object {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser had errors: %v", p.Errors())
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	t.Helper()

	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("object is not *object.Integer. got=%T (%+v)", obj, obj)
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value. expected=%d, got=%d", expected, result.Value)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	t.Helper()

	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("object is not *object.Boolean. got=%T (%+v)", obj, obj)
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value. expected=%t, got=%t", expected, result.Value)
	}
}

func testNullObject(t *testing.T, obj object.Object) {
	t.Helper()

	if obj == nil {
		t.Fatalf("object is nil")
	}

	if obj.Type() != object.NULL {
		t.Fatalf("object is not NULL. got=%T (%+v)", obj, obj)
	}
}

func testErrorObject(t *testing.T, obj object.Object, expectedMessage string) {
	t.Helper()

	errObj, ok := obj.(*object.Error)
	if !ok {
		t.Fatalf("object is not *object.Error. got=%T (%+v)", obj, obj)
	}

	if errObj.Message != expectedMessage {
		t.Fatalf("wrong error message. expected=%q, got=%q", expectedMessage, errObj.Message)
	}
}
