package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"modern-basic/ast"
	"sort"
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
	CAUGHT_ERROR ObjectType = "CAUGHT_ERROR"
	FUNCTION     ObjectType = "FUNCTION"
	BUILTIN      ObjectType = "BUILTIN"
	ARRAY        ObjectType = "ARRAY"
	HASH         ObjectType = "HASH"
)

// HashKey identifies a hashable key value.
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Hashable marks runtime objects that can be used as hash keys.
type Hashable interface {
	HashKey() HashKey
}

// Object is the runtime value interface used by the evaluator.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// BuiltinFunction defines a native function callable by the interpreter.
type BuiltinFunction func(args ...Object) Object

// Integer wraps an int64 runtime value.
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// HashKey returns a stable hash key for integers.
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Boolean wraps a boolean runtime value.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// HashKey returns a stable hash key for booleans.
func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return HashKey{Type: b.Type(), Value: 1}
	}

	return HashKey{Type: b.Type(), Value: 0}
}

// String wraps a string runtime value.
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING }

func (s *String) Inspect() string { return s.Value }

// HashKey returns a stable hash key for strings.
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

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

// CaughtError wraps an error intercepted by try/catch so it can be used as a regular value.
type CaughtError struct {
	Err *Error
}

func (ce *CaughtError) Type() ObjectType { return CAUGHT_ERROR }

func (ce *CaughtError) Inspect() string {
	if ce.Err == nil {
		return "ERROR:"
	}

	return ce.Err.Inspect()
}

// Function stores a user-defined function and its lexical environment.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Builtin wraps a native Go function.
type Builtin struct {
	Fn BuiltinFunction
}

// Array represents an ordered list of runtime objects.
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY }

func (a *Array) Inspect() string {
	parts := make([]string, 0, len(a.Elements))
	for _, el := range a.Elements {
		parts = append(parts, el.Inspect())
	}

	return "[" + strings.Join(parts, ", ") + "]"
}

// HashPair stores the original key object and mapped value.
type HashPair struct {
	Key   Object
	Value Object
}

// Hash represents key-value mappings for hashable keys.
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH }

func (h *Hash) Inspect() string {
	parts := make([]string, 0, len(h.Pairs))
	for _, pair := range h.Pairs {
		parts = append(parts, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}
	sort.Strings(parts)

	return "{" + strings.Join(parts, ", ") + "}"
}

func (f *Function) Type() ObjectType { return FUNCTION }

func (b *Builtin) Type() ObjectType { return BUILTIN }

func (b *Builtin) Inspect() string { return "<builtin>" }

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
