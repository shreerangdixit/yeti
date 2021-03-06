package eval

import (
	"fmt"
	"strconv"

	"github.com/shreerangdixit/yeti/ast"
	"github.com/shreerangdixit/yeti/lex"
)

type Evaluator struct {
	Importer *Importer

	env      *Environment
	deferred []ast.CallNode
}

func NewEvaluator() *Evaluator {
	e := Evaluator{
		env:      NewEnvironment(),
		deferred: make([]ast.CallNode, 0, 20),
	}
	e.Importer = NewImporter(&e)
	return &e
}

func (e *Evaluator) Evaluate(root ast.Node) (Object, error) {
	return e.eval(root)
}

func (e *Evaluator) eval(node ast.Node) (Object, error) {
	switch node := node.(type) {
	case ast.ProgramNode:
		obj, err := e.evalProgramNode(node)
		return e.wrapResult(node, obj, err)
	case ast.BlockNode:
		obj, err := e.evalBlockNode(node)
		return e.wrapResult(node, obj, err)
	case ast.VarStmtNode:
		obj, err := e.evalVarStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.ExpStmtNode:
		obj, err := e.evalExpStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.IfStmtNode:
		obj, err := e.evalIfStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.WhileStmtNode:
		obj, err := e.evalWhileStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.BreakStmtNode:
		obj, err := e.evalBreakStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.ContinueStmtNode:
		obj, err := e.evalContinueStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.ReturnStmtNode:
		obj, err := e.evalReturnStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.AssignmentNode:
		obj, err := e.evalAssignmentNode(node)
		return e.wrapResult(node, obj, err)
	case ast.LogicalAndNode:
		obj, err := e.evalLogicalAndNode(node)
		return e.wrapResult(node, obj, err)
	case ast.LogicalOrNode:
		obj, err := e.evalLogicalOrNode(node)
		return e.wrapResult(node, obj, err)
	case ast.ExpNode:
		obj, err := e.eval(node.Exp)
		return e.wrapResult(node, obj, err)
	case ast.TernaryOpNode:
		obj, err := e.evalTernaryOpNode(node)
		return e.wrapResult(node, obj, err)
	case ast.BinaryOpNode:
		obj, err := e.evalBinaryOpNode(node)
		return e.wrapResult(node, obj, err)
	case ast.UnaryOpNode:
		obj, err := e.evalUnaryOpNode(node)
		return e.wrapResult(node, obj, err)
	case ast.IdentifierNode:
		obj, err := e.evalIdentifierNode(node)
		return e.wrapResult(node, obj, err)
	case ast.NumberNode:
		obj, err := e.evalNumberNode(node)
		return e.wrapResult(node, obj, err)
	case ast.BooleanNode:
		obj, err := e.evalBooleanNode(node)
		return e.wrapResult(node, obj, err)
	case ast.StringNode:
		obj, err := e.evalStringNode(node)
		return e.wrapResult(node, obj, err)
	case ast.ListNode:
		obj, err := e.evalListNode(node)
		return e.wrapResult(node, obj, err)
	case ast.MapNode:
		obj, err := e.evalMapNode(node)
		return e.wrapResult(node, obj, err)
	case ast.NilNode:
		obj, err := e.evalNilNode(node)
		return e.wrapResult(node, obj, err)
	case ast.CallNode:
		obj, err := e.evalCallNode(node)
		return e.wrapResult(node, obj, err)
	case ast.IndexOfNode:
		obj, err := e.evalIndexOfNode(node)
		return e.wrapResult(node, obj, err)
	case ast.FunctionNode:
		obj, err := e.evalFunctionNode(node)
		return e.wrapResult(node, obj, err)
	case ast.DeferStmtNode:
		obj, err := e.evalDeferStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.AssertStmtNode:
		obj, err := e.evalAssertStmtNode(node)
		return e.wrapResult(node, obj, err)
	case ast.CommentNode:
		obj, err := e.evalCommentNode(node)
		return e.wrapResult(node, obj, err)
	case ast.ImportStmtNode:
		obj, err := e.evalImportNode(node)
		return e.wrapResult(node, obj, err)
	}
	return NIL, fmt.Errorf("invalid node: %T", node)
}

func (e *Evaluator) wrapResult(node ast.Node, obj Object, err error) (Object, error) {
	if err != nil {
		switch err := err.(type) {
		case BreakError:
		case ContinueError:
		case ReturnError:
			return obj, err
		case EvaluateError:
			return obj, NewEvaluateError(node, err, WithInnerError(err))
		default:
			return obj, NewEvaluateError(node, err)
		}
	}
	return obj, err
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
	return e.evalBlockNodeWithEnv(node, NewEnvironment().WithEnclosing(e.env))
}

func (e *Evaluator) evalBlockNodeWithEnv(node ast.BlockNode, env *Environment) (Object, error) {
	prev := e.env
	// Reset environment at the end of block scope
	defer func() {
		e.env = prev
	}()

	// New environment at the beginning of block scope
	e.env = env
	for _, node := range node.Declarations {
		_, err := e.eval(node)
		if err != nil {
			return NIL, err
		}
	}

	return e.runDeferred()
}

func (e *Evaluator) runDeferred() (Object, error) {
	deferred := e.deferred
	e.deferred = make([]ast.CallNode, 0, 20)
	for _, call := range deferred {
		o, err := e.eval(call)
		if err != nil {
			return o, err
		}
	}

	return NIL, nil
}

func (e *Evaluator) evalVarStmtNode(node ast.VarStmtNode) (Object, error) {
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
	} else if node.FalseStmt != nil {
		return e.eval(node.FalseStmt)
	} else {
		return NIL, nil
	}
}

func (e *Evaluator) evalWhileStmtNode(node ast.WhileStmtNode) (Object, error) {
	for {
		condition, err := e.eval(node.Condition)
		if err != nil {
			return NIL, err
		}

		if !IsTruthy(condition) {
			break
		}

		_, err = e.eval(node.Body)
		if err != nil {
			switch err := err.(type) {
			case *BreakError:
				return NIL, nil
			case *ContinueError:
				continue
			default:
				return NIL, err
			}
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalBreakStmtNode(node ast.BreakStmtNode) (Object, error) {
	return NIL, NewBreakError()
}

func (e *Evaluator) evalContinueStmtNode(node ast.ContinueStmtNode) (Object, error) {
	return NIL, NewContinueError()
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
		return TRUE, nil
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
	case lex.TT_PLUS:
		return Add(left, right)
	case lex.TT_MINUS:
		return Subtract(left, right)
	case lex.TT_DIVIDE:
		return Divide(left, right)
	case lex.TT_MULTIPLY:
		return Multiply(left, right)
	case lex.TT_MODULO:
		return Modulo(left, right)
	case lex.TT_EQ:
		return EqualTo(left, right), nil
	case lex.TT_NEQ:
		return NotEqualTo(left, right), nil
	case lex.TT_LT:
		return LessThan(left, right), nil
	case lex.TT_LTE:
		return LessThanEq(left, right), nil
	case lex.TT_GT:
		return GreaterThan(left, right), nil
	case lex.TT_GTE:
		return GreaterThanEq(left, right), nil
	}
	return NIL, fmt.Errorf("invalid binary op: %s", node.Op.Type)
}

func (e *Evaluator) evalUnaryOpNode(node ast.UnaryOpNode) (Object, error) {
	val, err := e.eval(node.Operand)
	if err != nil {
		return NIL, err
	}

	if node.Op.Type == lex.TT_MINUS {
		return Negate(val)
	} else if node.Op.Type == lex.TT_NOT {
		return Not(val)
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

	return NewNumber(val), nil
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

func (e *Evaluator) evalListNode(node ast.ListNode) (Object, error) {
	elements, err := e.evalNodes(node.Elements)
	if err != nil {
		return nil, err
	}

	return NewList(elements), nil
}

func (e *Evaluator) evalMapNode(node ast.MapNode) (Object, error) {
	m := NewMap()

	for _, kvp := range node.Elements {
		key, err := e.eval(kvp.Key)
		if err != nil {
			return NIL, err
		}

		value, err := e.eval(kvp.Value)
		if err != nil {
			return NIL, err
		}

		m, err = m.Add(key, value)
		if err != nil {
			return NIL, err
		}
	}
	return m, nil
}

func (e *Evaluator) evalNilNode(node ast.NilNode) (Object, error) {
	return NIL, nil
}

func (e *Evaluator) evalCallNode(node ast.CallNode) (Object, error) {
	calleeNode, err := e.eval(node.Callee)
	if err != nil {
		return NIL, err
	}

	callable, ok := calleeNode.(Callable)
	if !ok { // If the callee node itself isn't callable, check if it's value is callable
		calleeValue, err := e.env.Get(calleeNode.String())
		if err != nil {
			return NIL, fmt.Errorf("%s is not callable", calleeNode.Type())
		}

		callable, ok = calleeValue.(Callable)
		if !ok {
			return NIL, fmt.Errorf("%s is not callable", calleeValue.Type())
		}
	}

	if !callable.Variadic() && callable.Arity() != len(node.Arguments) {
		return NIL, fmt.Errorf("incorrect number of arguments to %s - %d expected %d provided", callable, callable.Arity(), len(node.Arguments))
	}

	argValues, err := e.evalNodes(node.Arguments)
	if err != nil {
		return NIL, err
	}

	return callable.Call(e, argValues)
}

func (e *Evaluator) evalIndexOfNode(node ast.IndexOfNode) (Object, error) {
	seq, err := e.eval(node.Sequence)
	if err != nil {
		return nil, err
	}

	idx, err := e.eval(node.Index)
	if err != nil {
		return nil, err
	}

	return ItemAtIndex(seq, idx)
}

func (e *Evaluator) evalFunctionNode(node ast.FunctionNode) (Object, error) {
	fun := NewUserFunction(node, e.env)
	return fun, e.env.Declare(fun.Name(), fun)
}

func (e *Evaluator) evalReturnStmtNode(node ast.ReturnStmtNode) (Object, error) {
	val, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	return NIL, NewReturnError(val)
}

func (e *Evaluator) evalDeferStmtNode(node ast.DeferStmtNode) (Object, error) {
	e.deferred = append(e.deferred, node.Call)
	return NIL, nil
}

func (e *Evaluator) evalAssertStmtNode(node ast.AssertStmtNode) (Object, error) {
	exp, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	if !IsTruthy(exp) {
		return NIL, NewAssertError(node.Exp)
	}
	return NIL, nil
}

func (e *Evaluator) evalImportNode(node ast.ImportStmtNode) (Object, error) {
	m := NewFileModule(node.Name.Token.Literal)
	return NIL, e.Importer.Import(m)
}

func (e *Evaluator) evalCommentNode(node ast.CommentNode) (Object, error) {
	return NIL, nil
}

func (e *Evaluator) evalNodes(argNodes []ast.Node) ([]Object, error) {
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
