package runtime

import (
	"fmt"
	"math"
	"time"
)

type FunctionHandler func(e *Evaluator, args []Object) (Object, error)

type Function struct {
	name    string
	arity   int
	handler FunctionHandler
}

func (f Function) Type() ObjectType                                 { return FUNC_OBJ }
func (f Function) String() string                                   { return f.name }
func (f Function) Arity() int                                       { return f.arity }
func (f Function) Call(e *Evaluator, args []Object) (Object, error) { return f.handler(e, args) }

var NativeFunctions = []Function{
	{
		name:    "sleep",
		arity:   1,
		handler: sleepHandler,
	},
	{
		name:    "time",
		arity:   0,
		handler: timeHandler,
	},
	{
		name:    "abs",
		arity:   1,
		handler: absHandler,
	},
	{
		name:    "max",
		arity:   2,
		handler: maxHandler,
	},
	{
		name:    "min",
		arity:   2,
		handler: minHandler,
	},
	{
		name:    "type",
		arity:   1,
		handler: typeHandler,
	},
}

func sleepHandler(e *Evaluator, args []Object) (Object, error) {
	arg, ok := args[0].(Float64)
	if !ok {
		return NIL, fmt.Errorf("sleep expects an argument of type float64")
	}

	time.Sleep(time.Duration(arg.Value) * time.Second)
	return NIL, nil
}

func timeHandler(e *Evaluator, args []Object) (Object, error) {
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	return NewFloat64(float64(ms)), nil
}

func absHandler(e *Evaluator, args []Object) (Object, error) {
	arg, ok := args[0].(Float64)
	if !ok {
		return NIL, fmt.Errorf("abs expects an argument of type float64")
	}
	return NewFloat64(math.Abs(arg.Value)), nil
}

func maxHandler(e *Evaluator, args []Object) (Object, error) {
	arg1, ok := args[0].(Float64)
	if !ok {
		return NIL, fmt.Errorf("first argument to max() must be a number")
	}

	arg2, ok := args[1].(Float64)
	if !ok {
		return NIL, fmt.Errorf("second argument to max() must be a number")
	}

	return NewFloat64(math.Max(arg1.Value, arg2.Value)), nil
}

func minHandler(e *Evaluator, args []Object) (Object, error) {
	arg1, ok := args[0].(Float64)
	if !ok {
		return NIL, fmt.Errorf("first argument to min() must be a number")
	}

	arg2, ok := args[1].(Float64)
	if !ok {
		return NIL, fmt.Errorf("second argument to min() must be a number")
	}

	return NewFloat64(math.Min(arg1.Value, arg2.Value)), nil
}

func typeHandler(e *Evaluator, args []Object) (Object, error) {
	arg := args[0]
	return NewType(arg.Type()), nil
}