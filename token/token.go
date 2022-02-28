package token

import (
	"fmt"
)

var keywords = map[string]TokenType{
	"var":      TT_VAR,
	"fun":      TT_FUNCTION,
	"if":       TT_IF,
	"else":     TT_ELSE,
	"true":     TT_TRUE,
	"false":    TT_FALSE,
	"return":   TT_RETURN,
	"while":    TT_WHILE,
	"nil":      TT_NIL,
	"break":    TT_BREAK,
	"continue": TT_CONTINUE,
	"defer":    TT_DEFER,
	"assert":   TT_ASSERT,
}

func LookupIdentifierType(v string) TokenType {
	if val, ok := keywords[v]; ok {
		return val
	}
	return TT_IDENTIFIER
}

type Position struct {
	Line   int
	Column int
}

func (p Position) String() string {
	return fmt.Sprintf("Pos %d:%d", p.Line, p.Column)
}

type TokenType int

type Token struct {
	Type          TokenType
	Literal       string
	BeginPosition Position
	EndPosition   Position
}

func (t Token) String() string {
	if t.Type == TT_IDENTIFIER || t.Type == TT_NUMBER {
		return fmt.Sprintf("%s:%s (%s - %s)", t.Type, t.Literal, t.BeginPosition, t.EndPosition)
	}
	return fmt.Sprintf("%s (%s - %s)", t.Literal, t.BeginPosition, t.EndPosition)
}

const (
	TT_ILLEGAL TokenType = iota
	TT_EOF

	// Identifier + Literals
	TT_IDENTIFIER
	TT_NUMBER
	TT_STRING

	// Operators
	TT_ASSIGN
	TT_PLUS
	TT_MINUS
	TT_DIVIDE
	TT_MULTIPLY
	TT_MODULO
	TT_NOT
	TT_EQ
	TT_NEQ
	TT_LT
	TT_LTE
	TT_GT
	TT_GTE
	TT_LOGICAL_AND
	TT_LOGICAL_OR

	// Delimiters
	TT_COMMA
	TT_COLON
	TT_QUESTION

	// Parens + Braces
	TT_LPAREN
	TT_RPAREN
	TT_LBRACE
	TT_RBRACE
	TT_LBRACKET
	TT_RBRACKET

	// Keywords
	TT_FUNCTION
	TT_VAR
	TT_IF
	TT_ELSE
	TT_TRUE
	TT_FALSE
	TT_RETURN
	TT_WHILE
	TT_BREAK
	TT_CONTINUE
	TT_NIL
	TT_DEFER
	TT_ASSERT

	// Misc
	TT_COMMENT
)

func (t TokenType) String() string {
	switch t {
	case TT_ILLEGAL:
		return "ILLEGAL"
	case TT_EOF:
		return "EOF"
	case TT_IDENTIFIER:
		return "IDENT"
	case TT_NUMBER:
		return "NUM"
	case TT_STRING:
		return "STR"
	case TT_ASSIGN:
		return "="
	case TT_PLUS:
		return "+"
	case TT_MINUS:
		return "-"
	case TT_DIVIDE:
		return "/"
	case TT_MULTIPLY:
		return "*"
	case TT_NOT:
		return "!"
	case TT_EQ:
		return "=="
	case TT_NEQ:
		return "!="
	case TT_LT:
		return "<"
	case TT_LTE:
		return "<="
	case TT_GT:
		return ">"
	case TT_GTE:
		return ">="
	case TT_LOGICAL_AND:
		return "&&"
	case TT_LOGICAL_OR:
		return "||"
	case TT_MODULO:
		return "%"
	case TT_COMMA:
		return ","
	case TT_QUESTION:
		return "?"
	case TT_COLON:
		return ":"
	case TT_COMMENT:
		return "//"
	case TT_LPAREN:
		return "("
	case TT_RPAREN:
		return ")"
	case TT_LBRACE:
		return "{"
	case TT_RBRACE:
		return "}"
	case TT_LBRACKET:
		return "["
	case TT_RBRACKET:
		return "]"
	case TT_FUNCTION:
		return "fun"
	case TT_VAR:
		return "var"
	case TT_IF:
		return "if"
	case TT_ELSE:
		return "else"
	case TT_TRUE:
		return "true"
	case TT_FALSE:
		return "false"
	case TT_RETURN:
		return "return"
	case TT_WHILE:
		return "while"
	case TT_BREAK:
		return "break"
	case TT_CONTINUE:
		return "continue"
	case TT_NIL:
		return "nil"
	case TT_DEFER:
		return "defer"
	case TT_ASSERT:
		return "assert"
	default:
		return "<UNKNOWN>"
	}
}
