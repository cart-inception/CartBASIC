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
		{name: "function", obj: &Function{Parameters: []*ast.Identifier{{Value: "x"}}, Body: &ast.BlockStatement{}}, expectedType: FUNCTION, expectedText: "fn(x) {}"},
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
