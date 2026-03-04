package parser

import (
	"fmt"
	"modern-basic/ast"
	"modern-basic/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = true;
let foobar = y;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return true;
return foobar;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return'. got=%q", returnStmt.TokenLiteral())
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
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("ident.Value not %q. got=%q", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral not %q. got=%q", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Fatalf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal.TokenLiteral not %q. got=%q", "5", literal.TokenLiteral())
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testStringLiteral(t, stmt.Expression, "hello world") {
		return
	}
}

func TestParsingStringInfixExpression(t *testing.T) {
	input := `"foo" + "bar";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.InfixExpression. got=%T", stmt.Expression)
	}

	if !testStringLiteral(t, exp.Left, "foo") {
		return
	}

	if exp.Operator != "+" {
		t.Fatalf("exp.Operator is not %q. got=%q", "+", exp.Operator)
	}

	if !testStringLiteral(t, exp.Right, "bar") {
		return
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q. got=%q", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
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
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Fatalf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfElseIfElseExpression(t *testing.T) {
	input := `if (x < y) { x; } elseif (x == y) { y; } else { 0; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong statement count. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expression not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if len(ifExp.Branches) != 2 {
		t.Fatalf("expected 2 branches, got=%d", len(ifExp.Branches))
	}

	if !testInfixExpression(t, ifExp.Branches[0].Condition, "x", "<", "y") {
		return
	}

	if len(ifExp.Branches[0].Consequence.Statements) != 1 {
		t.Fatalf("expected first consequence statement count 1, got=%d", len(ifExp.Branches[0].Consequence.Statements))
	}

	if !testInfixExpression(t, ifExp.Branches[1].Condition, "x", "==", "y") {
		return
	}

	if ifExp.Alternative == nil {
		t.Fatalf("expected else alternative block")
	}

	if len(ifExp.Alternative.Statements) != 1 {
		t.Fatalf("expected alternative statement count 1, got=%d", len(ifExp.Alternative.Statements))
	}
}

func TestWhileExpression(t *testing.T) {
	input := `while (x < 10) { x = x + 1; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	whileExp, ok := stmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("expression not *ast.WhileExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, whileExp.Condition, "x", "<", 10) {
		return
	}

	if len(whileExp.Body.Statements) != 1 {
		t.Fatalf("expected body statement count 1, got=%d", len(whileExp.Body.Statements))
	}
}

func TestForExpression(t *testing.T) {
	input := `for (let i = 0; i < 3; i = i + 1) { let sum = i; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	forExp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expression not *ast.ForExpression. got=%T", stmt.Expression)
	}

	init, ok := forExp.Init.(*ast.LetStatement)
	if !ok {
		t.Fatalf("for init not *ast.LetStatement. got=%T", forExp.Init)
	}
	if init.Name.Value != "i" {
		t.Fatalf("for init variable wrong. got=%s", init.Name.Value)
	}

	if !testInfixExpression(t, forExp.Condition, "i", "<", 3) {
		return
	}

	post, ok := forExp.Post.(*ast.AssignStatement)
	if !ok {
		t.Fatalf("for post not *ast.AssignStatement. got=%T", forExp.Post)
	}

	if post.Name.Value != "i" {
		t.Fatalf("for post variable wrong. got=%s", post.Name.Value)
	}

	if len(forExp.Body.Statements) != 1 {
		t.Fatalf("expected for body statement count 1, got=%d", len(forExp.Body.Statements))
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expression not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("wrong parameter count. got=%d", len(function.Parameters))
	}

	if function.Parameters[0].Value != "x" || function.Parameters[1].Value != "y" {
		t.Fatalf("wrong parameters: got=%s,%s", function.Parameters[0].Value, function.Parameters[1].Value)
	}

	if len(function.Body.Statements) != 1 {
		t.Fatalf("wrong body statement count. got=%d", len(function.Body.Statements))
	}
}

func TestFunctionDefinitionParsing(t *testing.T) {
	input := `function add(x, y) { return x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("statement not *ast.LetStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Fatalf("function name wrong. got=%s", stmt.Name.Value)
	}

	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("let value not *ast.FunctionLiteral. got=%T", stmt.Value)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("wrong parameter count. got=%d", len(function.Parameters))
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong arg count. got=%d", len(exp.Arguments))
	}

	if !testLiteralExpression(t, exp.Arguments[0], 1) {
		return
	}
	if !testInfixExpression(t, exp.Arguments[1], 2, "*", 3) {
		return
	}
	if !testInfixExpression(t, exp.Arguments[2], 4, "+", 5) {
		return
	}
}

func TestDotExpressionParsing(t *testing.T) {
	input := `File.Read`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	dotExp, ok := stmt.Expression.(*ast.DotExpression)
	if !ok {
		t.Fatalf("expression not *ast.DotExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, dotExp.Left, "File") {
		return
	}

	if dotExp.Right == nil || dotExp.Right.Value != "Read" {
		t.Fatalf("dot right side not Read. got=%v", dotExp.Right)
	}
}

func TestNamespacedCallExpressionParsing(t *testing.T) {
	input := `File.Read(path)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	callExp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression not *ast.CallExpression. got=%T", stmt.Expression)
	}

	dotExp, ok := callExp.Function.(*ast.DotExpression)
	if !ok {
		t.Fatalf("call function not *ast.DotExpression. got=%T", callExp.Function)
	}

	if !testIdentifier(t, dotExp.Left, "File") {
		return
	}

	if dotExp.Right == nil || dotExp.Right.Value != "Read" {
		t.Fatalf("dot right side not Read. got=%v", dotExp.Right)
	}

	if len(callExp.Arguments) != 1 {
		t.Fatalf("wrong arg count. got=%d", len(callExp.Arguments))
	}

	if !testIdentifier(t, callExp.Arguments[0], "path") {
		return
	}
}

func TestArrayLiteralParsing(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) wrong. got=%d", len(array.Elements))
	}

	if !testIntegerLiteral(t, array.Elements[0], 1) {
		return
	}
	if !testInfixExpression(t, array.Elements[1], 2, "*", 2) {
		return
	}
	if !testInfixExpression(t, array.Elements[2], 3, "+", 3) {
		return
	}
}

func TestIndexExpressionParsing(t *testing.T) {
	input := `myList[1 + 1]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myList") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestHashLiteralParsingStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{"one": 1, "two": 2, "three": 3}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("key is not *ast.StringLiteral. got=%T", key)
		}

		expectedValue, ok := expected[literal.Value]
		if !ok {
			t.Fatalf("unexpected hash key %q", literal.Value)
		}

		if !testIntegerLiteral(t, value, expectedValue) {
			return
		}
	}
}

func TestHashLiteralParsingWithExpressionValues(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
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
			t.Fatalf("key is not *ast.StringLiteral. got=%T", key)
		}

		testFn, ok := tests[literal.Value]
		if !ok {
			t.Fatalf("no test function for key %q", literal.Value)
		}

		testFn(value)
	}
}

func TestIndexExpressionPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))[0]", "(add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))[0])"},
		{"arr[0][1]", "((arr[0])[1])"},
		{"File.Read(path)", "File.Read(path)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Fatalf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name == nil || letStmt.Name.Value != name {
		t.Fatalf("letStmt.Name.Value not %q. got=%v", name, letStmt.Name)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Fatalf("letStmt.Name.TokenLiteral not %q. got=%q", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Fatalf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("integ.TokenLiteral not %d. got=%q", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Fatalf("ident.Value not %q. got=%q", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Fatalf("ident.TokenLiteral not %q. got=%q", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Fatalf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Fatalf("bo.TokenLiteral not %t. got=%q", value, bo.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", exp)
		return false
	}

	if str.Value != value {
		t.Fatalf("str.Value not %q. got=%q", value, str.Value)
		return false
	}

	if str.TokenLiteral() != value {
		t.Fatalf("str.TokenLiteral not %q. got=%q", value, str.TokenLiteral())
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
		return testBooleanLiteral(t, exp, v)
	default:
		t.Fatalf("type of exp not handled. got=%T", exp)
		return false
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Fatalf("exp.Operator is not %q. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Fatalf("parser has %d errors:\n%s", len(errors), joinLines(errors))
}

func joinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	out := ""
	for i, line := range lines {
		if i > 0 {
			out += "\n"
		}
		out += line
	}

	return out
}
