package runtime

import (
	"fmt"
)

// IEEE 754 floating point number type
// Implements the following interfaces
// Object
// Truthifier
// Negator
// LessThanComparator
// GreaterThanComparator
// EqualToComparator
// Adder
// Subtractor
// Multiplier
// Divider
type Number struct{ Value float64 }

func NewNumber(value float64) Number           { return Number{Value: value} }
func (f Number) Type() ObjectType              { return TypeNumber }
func (f Number) String() string                { return fmt.Sprintf("%v", f.Value) }
func (f Number) Truthy() Bool                  { return NewBool(f.Value != 0) }
func (f Number) Negate() (Object, error)       { return f.Multiply(NewNumber(-1)) }
func (f Number) LessThan(other Object) Bool    { return NewBool(f.Value < other.(Number).Value) }
func (f Number) GreaterThan(other Object) Bool { return NewBool(f.Value > other.(Number).Value) }
func (f Number) EqualTo(other Object) Bool     { return NewBool(f.Value == other.(Number).Value) }

func (f Number) Add(other Object) (Object, error) {
	return NewNumber(f.Value + other.(Number).Value), nil
}

func (f Number) Subtract(other Object) (Object, error) {
	return NewNumber(f.Value - other.(Number).Value), nil
}

func (f Number) Multiply(other Object) (Object, error) {
	return NewNumber(f.Value * other.(Number).Value), nil
}

func (f Number) Divide(other Object) (Object, error) {
	if other.(Number).Value == 0 {
		return nil, fmt.Errorf("Divide by zero error")
	}
	return NewNumber(f.Value / other.(Number).Value), nil
}

// Boolean type
// Implements the following interfaces
// Object
// Truthifier
// EqualToComparator
// Notter
type Bool struct{ Value bool }

var TRUE = NewBool(true)
var FALSE = NewBool(false)

func NewBool(value bool) Bool            { return Bool{Value: value} }
func (f Bool) Type() ObjectType          { return TypeBool }
func (f Bool) String() string            { return fmt.Sprintf("%v", f.Value) }
func (f Bool) EqualTo(other Object) Bool { return NewBool(f.Value == other.(Bool).Value) }
func (f Bool) Truthy() Bool              { return NewBool(f.Value) }
func (f Bool) Not() (Object, error)      { return NewBool(!f.Value), nil }

// String type
// Implements the following interfaces
// Object
// Sequence
// Truthifier
// Adder
// LessThanComparator
// GreaterThanComparator
// EqualToComparator
type String struct{ Value string }

func NewString(value string) String            { return String{Value: value} }
func (f String) Type() ObjectType              { return TypeString }
func (f String) String() string                { return f.Value }
func (f String) Truthy() Bool                  { return NewBool(f.Size().Value > 0) }
func (f String) LessThan(other Object) Bool    { return NewBool(f.Value < other.(String).Value) }
func (f String) GreaterThan(other Object) Bool { return NewBool(f.Value > other.(String).Value) }
func (f String) EqualTo(other Object) Bool     { return NewBool(f.Value == other.(String).Value) }
func (f String) Size() Number                  { return NewNumber(float64(len(f.Value))) }

func (f String) Add(other Object) (Object, error) {
	return NewString(f.Value + other.(String).Value), nil
}

func (f String) Index(n Number) (Object, error) {
	idx := int(n.Value)
	if idx >= len(f.Value) {
		return nil, fmt.Errorf("string index out of range")
	}
	return NewString(string(f.Value[idx])), nil
}

// Type information meta-type
// Implements the following interfaces
// Object
// Truthifier
// EqualToComparator
type Type struct{ Value ObjectType }

func NewType(value ObjectType) Type      { return Type{Value: value} }
func (f Type) Type() ObjectType          { return TypeType }
func (f Type) String() string            { return string(f.Value) }
func (f Type) Truthy() Bool              { return TRUE }
func (f Type) EqualTo(other Object) Bool { return NewBool(f.Value == other.(Type).Value) }

// Heterogenous list type
// Implements the following interfaces
// Object
// Sequence
// Truthifier
// Adder
// TODO:
// LessThanComparator
// GreaterThanComparator
// EqualToComparator
type List struct{ Values []Object }

func NewList(values []Object) List { return List{Values: values} }
func (f List) Type() ObjectType    { return TypeList }
func (f List) String() string      { return fmt.Sprintf("%v", f.Values) }
func (f List) Size() Number        { return NewNumber(float64(len(f.Values))) }
func (f List) Truthy() Bool        { return NewBool(f.Size().Value > 0) }

func (f List) Add(other Object) (Object, error) {
	l, ok := other.(List)
	if !ok {
		return nil, fmt.Errorf("cannot concatenate list with %s", other.Type())
	}

	return NewList(append(f.Values, l.Values...)), nil
}

func (f List) Index(n Number) (Object, error) {
	idx := int(n.Value)
	if idx >= len(f.Values) {
		return nil, fmt.Errorf("list index out of range")
	}
	return f.Values[idx], nil
}

// Nil type
// Implements the following interfaces
// Object
// Truthifier
// EqualToComparator
type Nil struct{}

func (f Nil) Type() ObjectType          { return TypeNil }
func (f Nil) String() string            { return "nil" }
func (f Nil) Truthy() Bool              { return FALSE }
func (f Nil) EqualTo(other Object) Bool { return TRUE }
