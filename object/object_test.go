package object

import (
	"modern-basic/ast"
	"testing"
)

func TestObjectInspectAndType(t *testing.T) {
	tests := []struct {
		name         string
		obj          Object
		expectedType ObjectType
		expectedText string
	}{
		{name: "integer", obj: &Integer{Value: 42}, expectedType: INTEGER, expectedText: "42"},
		{name: "boolean true", obj: &Boolean{Value: true}, expectedType: BOOLEAN, expectedText: "true"},
		{name: "string", obj: &String{Value: "hello"}, expectedType: STRING, expectedText: "hello"},
		{name: "null", obj: &Null{}, expectedType: NULL, expectedText: "null"},
		{name: "return value", obj: &ReturnValue{Value: &Integer{Value: 7}}, expectedType: RETURN_VALUE, expectedText: "7"},
		{name: "error", obj: &Error{Message: "type mismatch: INTEGER + BOOLEAN"}, expectedType: ERROR, expectedText: "ERROR: type mismatch: INTEGER + BOOLEAN"},
		{name: "caught error", obj: &CaughtError{Err: &Error{Message: "division by zero"}}, expectedType: CAUGHT_ERROR, expectedText: "ERROR: division by zero"},
		{name: "function", obj: &Function{Parameters: []*ast.Identifier{{Value: "x"}}, Body: &ast.BlockStatement{}}, expectedType: FUNCTION, expectedText: "fn(x) {}"},
		{name: "builtin", obj: &Builtin{Fn: func(args ...Object) Object { return &Null{} }}, expectedType: BUILTIN, expectedText: "<builtin>"},
		{name: "array", obj: &Array{Elements: []Object{&Integer{Value: 1}, &String{Value: "two"}}}, expectedType: ARRAY, expectedText: "[1, two]"},
	}

	for _, tt := range tests {
		if tt.obj.Type() != tt.expectedType {
			t.Fatalf("%s: expected type %q, got %q", tt.name, tt.expectedType, tt.obj.Type())
		}

		if tt.obj.Inspect() != tt.expectedText {
			t.Fatalf("%s: expected inspect %q, got %q", tt.name, tt.expectedText, tt.obj.Inspect())
		}
	}
}

func TestHashKeySupport(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Fatalf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Fatalf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Fatalf("strings with different content have same hash keys")
	}

	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}
	two := &Integer{Value: 2}

	if one1.HashKey() != one2.HashKey() {
		t.Fatalf("integers with same content have different hash keys")
	}

	if one1.HashKey() == two.HashKey() {
		t.Fatalf("integers with different content have same hash keys")
	}

	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Fatalf("booleans with same content have different hash keys")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Fatalf("booleans with different content have same hash keys")
	}
}

func TestHashObjectInspectAndType(t *testing.T) {
	key := &String{Value: "name"}
	hash := &Hash{
		Pairs: map[HashKey]HashPair{
			key.HashKey(): {Key: key, Value: &String{Value: "cart"}},
		},
	}

	if hash.Type() != HASH {
		t.Fatalf("expected type %q, got %q", HASH, hash.Type())
	}

	if got := hash.Inspect(); got != "{name: cart}" {
		t.Fatalf("expected inspect %q, got %q", "{name: cart}", got)
	}
}

func TestEnvironmentSetAndGet(t *testing.T) {
	env := NewEnvironment()
	value := &Integer{Value: 99}
	env.Set("x", value)

	obj, ok := env.Get("x")
	if !ok {
		t.Fatalf("expected binding for x")
	}

	integer, ok := obj.(*Integer)
	if !ok {
		t.Fatalf("expected *Integer, got %T", obj)
	}

	if integer.Value != 99 {
		t.Fatalf("expected value 99, got %d", integer.Value)
	}
}

func TestEnvironmentOuterScopeLookup(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("name", &String{Value: "cart"})
	inner := NewEnclosedEnvironment(outer)

	obj, ok := inner.Get("name")
	if !ok {
		t.Fatalf("expected to resolve name from outer scope")
	}

	str, ok := obj.(*String)
	if !ok {
		t.Fatalf("expected *String, got %T", obj)
	}

	if str.Value != "cart" {
		t.Fatalf("expected value cart, got %q", str.Value)
	}
}
