package object

import (
	"bytes"
	"fmt"
	"modern-basic/ast"
	"strings"
)

// ObjectType identifies a runtime value category.
type ObjectType string

const (
	INTEGER      ObjectType = "INTEGER"
	BOOLEAN      ObjectType = "BOOLEAN"
	STRING       ObjectType = "STRING"
	NULL         ObjectType = "NULL"
	RETURN_VALUE ObjectType = "RETURN_VALUE"
	ERROR        ObjectType = "ERROR"
	FUNCTION     ObjectType = "FUNCTION"
)

// Object is the runtime value interface used by the evaluator.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer wraps an int64 runtime value.
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean wraps a boolean runtime value.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// String wraps a string runtime value.
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING }

func (s *String) Inspect() string { return s.Value }

// Null represents the absence of a value.
type Null struct{}

func (n *Null) Type() ObjectType { return NULL }

func (n *Null) Inspect() string { return "null" }

// ReturnValue wraps a value produced by return statements.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE }

func (rv *ReturnValue) Inspect() string {
	if rv.Value == nil {
		return ""
	}
	return rv.Value.Inspect()
}

// Error describes evaluation failures that should halt execution.
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }

func (e *Error) Inspect() string { return "ERROR: " + strings.TrimSpace(e.Message) }

// Function stores a user-defined function and its lexical environment.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION }

func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("fn(")
	for i, p := range f.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(") {")
	if f.Body != nil {
		out.WriteString(f.Body.String())
	}
	out.WriteString("}")

	return out.String()
}

// Environment stores identifier bindings.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a fresh environment with an empty store.
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

// NewEnclosedEnvironment creates a new environment with an outer scope.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get returns a bound object by name, walking outer scopes if present.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set binds a name to an object and returns the stored value.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
