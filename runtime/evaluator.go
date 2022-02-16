package runtime

import (
	"fmt"
	"github.com/shreerangdixit/lox/ast"
	"github.com/shreerangdixit/lox/token"
	"strconv"
)

type Evaluator struct {
	env *Env
}

func NewEvaluator() *Evaluator {
	return &Evaluator{env: NewEnv()}
}

func (e *Evaluator) Evaluate(root ast.Node) (Object, error) {
	return e.eval(root)
}

func (e *Evaluator) eval(node ast.Node) (Object, error) {
	switch node := node.(type) {
	case ast.ProgramNode:
		return e.evalProgramNode(node)
	case ast.BlockNode:
		return e.evalBlockNode(node)
	case ast.LetStmtNode:
		return e.evalLetStmtNode(node)
	case ast.ExpStmtNode:
		return e.evalExpStmtNode(node)
	case ast.IfStmtNode:
		return e.evalIfStmtNode(node)
	case ast.PrintStmtNode:
		return e.evalPrintStmtNode(node)
	case ast.WhileStmtNode:
		return e.evalWhileStmtNode(node)
	case ast.AssignmentNode:
		return e.evalAssignmentNode(node)
	case ast.LogicalAndNode:
		return e.evalLogicalAndNode(node)
	case ast.LogicalOrNode:
		return e.evalLogicalOrNode(node)
	case ast.ExpNode:
		return e.eval(node.Exp)
	case ast.TernaryOpNode:
		return e.evalTernaryOpNode(node)
	case ast.BinaryOpNode:
		return e.evalBinaryOpNode(node)
	case ast.UnaryOpNode:
		return e.evalUnaryOpNode(node)
	case ast.IdentifierNode:
		return e.evalIdentifierNode(node)
	case ast.NumberNode:
		return e.evalNumberNode(node)
	case ast.BooleanNode:
		return e.evalBooleanNode(node)
	case ast.StringNode:
		return e.evalStringNode(node)
	case ast.NilNode:
		return e.evalNilNode(node)
	case ast.CallNode:
		return e.evalCallNode(node)
	}
	return NIL, fmt.Errorf("invalid node: %T", node)
}

func (e *Evaluator) evalProgramNode(node ast.ProgramNode) (Object, error) {
	for _, node := range node.Declarations {
		_, err := e.eval(node)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalBlockNode(node ast.BlockNode) (Object, error) {
	// Reset environment at the end of block scope
	prev := e.env
	defer func() {
		e.env = prev
	}()

	// New environment at the beginning of block scope
	e.env = NewEnvWithEnclosing(e.env)
	for _, node := range node.Declarations {
		_, err := e.eval(node)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalLetStmtNode(node ast.LetStmtNode) (Object, error) {
	value, err := e.eval(node.Value)
	if err != nil {
		return NIL, err
	}

	if err := e.env.Declare(node.Identifier.Token.Literal, value); err != nil {
		return NIL, err
	}
	return NIL, nil
}

func (e *Evaluator) evalExpStmtNode(node ast.ExpStmtNode) (Object, error) {
	return e.eval(node.Exp)
}

func (e *Evaluator) evalIfStmtNode(node ast.IfStmtNode) (Object, error) {
	value, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(value) {
		return e.eval(node.TrueStmt)
	} else {
		return e.eval(node.FalseStmt)
	}
}

func (e *Evaluator) evalPrintStmtNode(node ast.PrintStmtNode) (Object, error) {
	result, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	fmt.Printf("%s\n", result)

	return NIL, nil
}

func (e *Evaluator) evalWhileStmtNode(node ast.WhileStmtNode) (Object, error) {
	for {
		result, err := e.eval(node.Condition)
		if err != nil {
			return NIL, err
		}

		if !IsTruthy(result) {
			break
		}

		_, err = e.eval(node.Body)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalAssignmentNode(node ast.AssignmentNode) (Object, error) {
	value, err := e.eval(node.Value)
	if err != nil {
		return NIL, err
	}
	return NIL, e.env.Assign(node.Identifier.Token.Literal, value)
}

func (e *Evaluator) evalLogicalAndNode(node ast.LogicalAndNode) (Object, error) {
	left, err := e.eval(node.LHS)
	if err != nil {
		return NIL, err
	}

	right, err := e.eval(node.RHS)
	if err != nil {
		return NIL, err
	}

	return NewBool(IsTruthy(left) && IsTruthy(right)), nil
}

func (e *Evaluator) evalLogicalOrNode(node ast.LogicalOrNode) (Object, error) {
	left, err := e.eval(node.LHS)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(left) {
		return NewBool(true), nil
	}

	right, err := e.eval(node.RHS)
	if err != nil {
		return NIL, err
	}

	return NewBool(IsTruthy(right)), nil
}

func (e *Evaluator) evalTernaryOpNode(node ast.TernaryOpNode) (Object, error) {
	value, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(value) {
		return e.eval(node.TrueExp)
	} else {
		return e.eval(node.FalseExp)
	}
}

func (e *Evaluator) evalBinaryOpNode(node ast.BinaryOpNode) (Object, error) {
	left, err := e.eval(node.LeftExp)
	if err != nil {
		return NIL, err
	}

	right, err := e.eval(node.RightExp)
	if err != nil {
		return NIL, err
	}

	switch node.Op.Type {
	case token.TT_PLUS:
		return Add(left, right)
	case token.TT_MINUS:
		return Subtract(left, right)
	case token.TT_DIVIDE:
		return Divide(left, right)
	case token.TT_MULTIPLY:
		return Multiply(left, right)
	case token.TT_EQ:
		return EqualTo(left, right), nil
	case token.TT_NEQ:
		return NotEqualTo(left, right), nil
	case token.TT_LT:
		return LessThan(left, right), nil
	case token.TT_LTE:
		return LessThanEq(left, right), nil
	case token.TT_GT:
		return GreaterThan(left, right), nil
	case token.TT_GTE:
		return GreaterThanEq(left, right), nil
	}
	return NIL, fmt.Errorf("invalid binary op: %s", node.Op.Type)
}

func (e *Evaluator) evalUnaryOpNode(node ast.UnaryOpNode) (Object, error) {
	val, err := e.eval(node.Operand)
	if err != nil {
		return NIL, err
	}

	if node.Op.Type == token.TT_MINUS || node.Op.Type == token.TT_NOT {
		return Negate(val)
	}

	return NIL, fmt.Errorf("invalid unary op: %s", node.Op.Type)
}

func (e *Evaluator) evalIdentifierNode(node ast.IdentifierNode) (Object, error) {
	return e.env.Get(node.Token.Literal)
}

func (e *Evaluator) evalNumberNode(node ast.NumberNode) (Object, error) {
	val, err := strconv.ParseFloat(node.Token.Literal, 10)
	if err != nil {
		return NIL, err
	}

	return NewFloat64(val), nil
}

func (e *Evaluator) evalBooleanNode(node ast.BooleanNode) (Object, error) {
	val, err := strconv.ParseBool(node.Token.Literal)
	if err != nil {
		return NIL, err
	}

	return NewBool(val), nil
}

func (e *Evaluator) evalStringNode(node ast.StringNode) (Object, error) {
	return NewString(node.Token.Literal), nil
}

func (e *Evaluator) evalNilNode(node ast.NilNode) (Object, error) {
	return NIL, nil
}

func (e *Evaluator) evalCallNode(node ast.CallNode) (Object, error) {
	callee, err := e.eval(node.Callee)
	if err != nil {
		return NIL, fmt.Errorf("%s is not declared", node.Callee)
	}

	calleeValue, err := e.env.Get(callee.String())
	if err != nil {
		return NIL, fmt.Errorf("%s is not callable", callee.Type())
	}

	callable, ok := calleeValue.(Callable)
	if !ok {
		return NIL, fmt.Errorf("%s is not declared", calleeValue.Type())
	}

	if callable.Arity() != len(node.Arguments) {
		return NIL, fmt.Errorf(
			"incorrect number of arguments to %s - %d expected %d provided",
			callable,
			callable.Arity(),
			len(node.Arguments),
		)
	}

	argValues, err := e.makeCallArguments(node.Arguments)
	if err != nil {
		return NIL, err
	}

	return callable.Call(e, argValues)
}

func (e *Evaluator) makeCallArguments(argNodes []ast.Node) ([]Object, error) {
	argValues := make([]Object, 0, 255)
	for _, arg := range argNodes {
		argval, err := e.eval(arg)
		if err != nil {
			return []Object{}, err
		}

		argValues = append(argValues, argval)
	}
	return argValues, nil
}