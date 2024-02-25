package object

import (
	"bytes"
	"dux/ast"
	"fmt"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	NIL_OBJ          = "NIL"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type Nil struct { }

type ReturnValue struct {
	Value Object
}

type Error struct {
	Message string
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

type String struct {
	Value string
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

type Array struct {
	Elements []Object
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[uint64]HashPair
}

type Hashable interface {
	HashKey() uint64
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out strings.Builder

	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() uint64 {
	return uint64(i.Value)
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() uint64 {
	var value uint64

	if b.Value { value = 1 } else { value = 0 }

	return value
}

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string { return "nil" }

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR " + e.Message }

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string { return fmt.Sprintf("%q", s.Value) }
func (s *String) HashKey() uint64 {
	hk := fnv.New64a()
	hk.Write([]byte(s.Value))

	return hk.Sum64()
}

func (bi *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (bi *Builtin) Inspect() string { return "builtin function" }

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out strings.Builder

	elements := []string{}

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
