package lex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "operators_paren",
			input: "= / + - * , ( ) { } == ! != < <= > >= && ||",
			want: []Token{
				{Type: TT_ASSIGN, Literal: "=", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 1}},
				{Type: TT_DIVIDE, Literal: "/", BeginPosition: Position{Line: 1, Column: 3}, EndPosition: Position{Line: 1, Column: 3}},
				{Type: TT_PLUS, Literal: "+", BeginPosition: Position{Line: 1, Column: 5}, EndPosition: Position{Line: 1, Column: 5}},
				{Type: TT_MINUS, Literal: "-", BeginPosition: Position{Line: 1, Column: 7}, EndPosition: Position{Line: 1, Column: 7}},
				{Type: TT_MULTIPLY, Literal: "*", BeginPosition: Position{Line: 1, Column: 9}, EndPosition: Position{Line: 1, Column: 9}},
				{Type: TT_COMMA, Literal: ",", BeginPosition: Position{Line: 1, Column: 11}, EndPosition: Position{Line: 1, Column: 11}},
				{Type: TT_LPAREN, Literal: "(", BeginPosition: Position{Line: 1, Column: 13}, EndPosition: Position{Line: 1, Column: 13}},
				{Type: TT_RPAREN, Literal: ")", BeginPosition: Position{Line: 1, Column: 15}, EndPosition: Position{Line: 1, Column: 15}},
				{Type: TT_LBRACE, Literal: "{", BeginPosition: Position{Line: 1, Column: 17}, EndPosition: Position{Line: 1, Column: 17}},
				{Type: TT_RBRACE, Literal: "}", BeginPosition: Position{Line: 1, Column: 19}, EndPosition: Position{Line: 1, Column: 19}},
				{Type: TT_EQ, Literal: "==", BeginPosition: Position{Line: 1, Column: 21}, EndPosition: Position{Line: 1, Column: 22}},
				{Type: TT_NOT, Literal: "!", BeginPosition: Position{Line: 1, Column: 24}, EndPosition: Position{Line: 1, Column: 24}},
				{Type: TT_NEQ, Literal: "!=", BeginPosition: Position{Line: 1, Column: 26}, EndPosition: Position{Line: 1, Column: 27}},
				{Type: TT_LT, Literal: "<", BeginPosition: Position{Line: 1, Column: 29}, EndPosition: Position{Line: 1, Column: 29}},
				{Type: TT_LTE, Literal: "<=", BeginPosition: Position{Line: 1, Column: 31}, EndPosition: Position{Line: 1, Column: 32}},
				{Type: TT_GT, Literal: ">", BeginPosition: Position{Line: 1, Column: 34}, EndPosition: Position{Line: 1, Column: 34}},
				{Type: TT_GTE, Literal: ">=", BeginPosition: Position{Line: 1, Column: 36}, EndPosition: Position{Line: 1, Column: 37}},
				{Type: TT_LOGICAL_AND, Literal: "&&", BeginPosition: Position{Line: 1, Column: 39}, EndPosition: Position{Line: 1, Column: 40}},
				{Type: TT_LOGICAL_OR, Literal: "||", BeginPosition: Position{Line: 1, Column: 42}, EndPosition: Position{Line: 1, Column: 43}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 42}, EndPosition: Position{Line: 1, Column: 43}},
			},
		},
		{
			name:  "integers",
			input: "123 456 7890",
			want: []Token{
				{Type: TT_NUMBER, Literal: "123", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 3}},
				{Type: TT_NUMBER, Literal: "456", BeginPosition: Position{Line: 1, Column: 5}, EndPosition: Position{Line: 1, Column: 7}},
				{Type: TT_NUMBER, Literal: "7890", BeginPosition: Position{Line: 1, Column: 9}, EndPosition: Position{Line: 1, Column: 12}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 9}, EndPosition: Position{Line: 1, Column: 12}},
			},
		},
		{
			name:  "floats",
			input: "0.123 1.23",
			want: []Token{
				{Type: TT_NUMBER, Literal: "0.123", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 5}},
				{Type: TT_NUMBER, Literal: "1.23", BeginPosition: Position{Line: 1, Column: 7}, EndPosition: Position{Line: 1, Column: 10}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 7}, EndPosition: Position{Line: 1, Column: 10}},
			},
		},
		{
			name:  "bad_floats",
			input: ".123 1.23",
			want: []Token{
				{Type: TT_ILLEGAL, Literal: ".", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 1}},
				{Type: TT_NUMBER, Literal: "123", BeginPosition: Position{Line: 1, Column: 2}, EndPosition: Position{Line: 1, Column: 4}},
				{Type: TT_NUMBER, Literal: "1.23", BeginPosition: Position{Line: 1, Column: 6}, EndPosition: Position{Line: 1, Column: 9}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 6}, EndPosition: Position{Line: 1, Column: 9}},
			},
		},
		{
			name:  "identifiers",
			input: "X Y Z aa bb cc_c d",
			want: []Token{
				{Type: TT_IDENTIFIER, Literal: "X", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 1}},
				{Type: TT_IDENTIFIER, Literal: "Y", BeginPosition: Position{Line: 1, Column: 3}, EndPosition: Position{Line: 1, Column: 3}},
				{Type: TT_IDENTIFIER, Literal: "Z", BeginPosition: Position{Line: 1, Column: 5}, EndPosition: Position{Line: 1, Column: 5}},
				{Type: TT_IDENTIFIER, Literal: "aa", BeginPosition: Position{Line: 1, Column: 7}, EndPosition: Position{Line: 1, Column: 8}},
				{Type: TT_IDENTIFIER, Literal: "bb", BeginPosition: Position{Line: 1, Column: 10}, EndPosition: Position{Line: 1, Column: 11}},
				{Type: TT_IDENTIFIER, Literal: "cc_c", BeginPosition: Position{Line: 1, Column: 13}, EndPosition: Position{Line: 1, Column: 16}},
				{Type: TT_IDENTIFIER, Literal: "d", BeginPosition: Position{Line: 1, Column: 18}, EndPosition: Position{Line: 1, Column: 18}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 18}, EndPosition: Position{Line: 1, Column: 18}},
			},
		},
		{
			name:  "mixed",
			input: " {(a = b * 5) (c = 10.5 / z)} ",
			want: []Token{
				{Type: TT_LBRACE, Literal: "{", BeginPosition: Position{Line: 1, Column: 2}, EndPosition: Position{Line: 1, Column: 2}},
				{Type: TT_LPAREN, Literal: "(", BeginPosition: Position{Line: 1, Column: 3}, EndPosition: Position{Line: 1, Column: 3}},
				{Type: TT_IDENTIFIER, Literal: "a", BeginPosition: Position{Line: 1, Column: 4}, EndPosition: Position{Line: 1, Column: 4}},
				{Type: TT_ASSIGN, Literal: "=", BeginPosition: Position{Line: 1, Column: 6}, EndPosition: Position{Line: 1, Column: 6}},
				{Type: TT_IDENTIFIER, Literal: "b", BeginPosition: Position{Line: 1, Column: 8}, EndPosition: Position{Line: 1, Column: 8}},
				{Type: TT_MULTIPLY, Literal: "*", BeginPosition: Position{Line: 1, Column: 10}, EndPosition: Position{Line: 1, Column: 10}},
				{Type: TT_NUMBER, Literal: "5", BeginPosition: Position{Line: 1, Column: 12}, EndPosition: Position{Line: 1, Column: 12}},
				{Type: TT_RPAREN, Literal: ")", BeginPosition: Position{Line: 1, Column: 13}, EndPosition: Position{Line: 1, Column: 13}},
				{Type: TT_LPAREN, Literal: "(", BeginPosition: Position{Line: 1, Column: 15}, EndPosition: Position{Line: 1, Column: 15}},
				{Type: TT_IDENTIFIER, Literal: "c", BeginPosition: Position{Line: 1, Column: 16}, EndPosition: Position{Line: 1, Column: 16}},
				{Type: TT_ASSIGN, Literal: "=", BeginPosition: Position{Line: 1, Column: 18}, EndPosition: Position{Line: 1, Column: 18}},
				{Type: TT_NUMBER, Literal: "10.5", BeginPosition: Position{Line: 1, Column: 20}, EndPosition: Position{Line: 1, Column: 23}},
				{Type: TT_DIVIDE, Literal: "/", BeginPosition: Position{Line: 1, Column: 25}, EndPosition: Position{Line: 1, Column: 25}},
				{Type: TT_IDENTIFIER, Literal: "z", BeginPosition: Position{Line: 1, Column: 27}, EndPosition: Position{Line: 1, Column: 27}},
				{Type: TT_RPAREN, Literal: ")", BeginPosition: Position{Line: 1, Column: 28}, EndPosition: Position{Line: 1, Column: 28}},
				{Type: TT_RBRACE, Literal: "}", BeginPosition: Position{Line: 1, Column: 29}, EndPosition: Position{Line: 1, Column: 29}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 29}, EndPosition: Position{Line: 1, Column: 29}},
			},
		},
		{
			name:  "keywords",
			input: "var x = 10 y = fun foo(){} if else true false return",
			want: []Token{
				{Type: TT_VAR, Literal: "var", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 3}},
				{Type: TT_IDENTIFIER, Literal: "x", BeginPosition: Position{Line: 1, Column: 5}, EndPosition: Position{Line: 1, Column: 5}},
				{Type: TT_ASSIGN, Literal: "=", BeginPosition: Position{Line: 1, Column: 7}, EndPosition: Position{Line: 1, Column: 7}},
				{Type: TT_NUMBER, Literal: "10", BeginPosition: Position{Line: 1, Column: 9}, EndPosition: Position{Line: 1, Column: 10}},
				{Type: TT_IDENTIFIER, Literal: "y", BeginPosition: Position{Line: 1, Column: 12}, EndPosition: Position{Line: 1, Column: 12}},
				{Type: TT_ASSIGN, Literal: "=", BeginPosition: Position{Line: 1, Column: 14}, EndPosition: Position{Line: 1, Column: 14}},
				{Type: TT_FUNCTION, Literal: "fun", BeginPosition: Position{Line: 1, Column: 16}, EndPosition: Position{Line: 1, Column: 18}},
				{Type: TT_IDENTIFIER, Literal: "foo", BeginPosition: Position{Line: 1, Column: 20}, EndPosition: Position{Line: 1, Column: 22}},
				{Type: TT_LPAREN, Literal: "(", BeginPosition: Position{Line: 1, Column: 23}, EndPosition: Position{Line: 1, Column: 23}},
				{Type: TT_RPAREN, Literal: ")", BeginPosition: Position{Line: 1, Column: 24}, EndPosition: Position{Line: 1, Column: 24}},
				{Type: TT_LBRACE, Literal: "{", BeginPosition: Position{Line: 1, Column: 25}, EndPosition: Position{Line: 1, Column: 25}},
				{Type: TT_RBRACE, Literal: "}", BeginPosition: Position{Line: 1, Column: 26}, EndPosition: Position{Line: 1, Column: 26}},
				{Type: TT_IF, Literal: "if", BeginPosition: Position{Line: 1, Column: 28}, EndPosition: Position{Line: 1, Column: 29}},
				{Type: TT_ELSE, Literal: "else", BeginPosition: Position{Line: 1, Column: 31}, EndPosition: Position{Line: 1, Column: 34}},
				{Type: TT_TRUE, Literal: "true", BeginPosition: Position{Line: 1, Column: 36}, EndPosition: Position{Line: 1, Column: 39}},
				{Type: TT_FALSE, Literal: "false", BeginPosition: Position{Line: 1, Column: 41}, EndPosition: Position{Line: 1, Column: 45}},
				{Type: TT_RETURN, Literal: "return", BeginPosition: Position{Line: 1, Column: 47}, EndPosition: Position{Line: 1, Column: 52}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 47}, EndPosition: Position{Line: 1, Column: 52}},
			},
		},
		{
			name:  "strings",
			input: "\"foo\" \"bar\" \"foo bar\"",
			want: []Token{
				{Type: TT_STRING, Literal: "foo", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 5}},
				{Type: TT_STRING, Literal: "bar", BeginPosition: Position{Line: 1, Column: 7}, EndPosition: Position{Line: 1, Column: 11}},
				{Type: TT_STRING, Literal: "foo bar", BeginPosition: Position{Line: 1, Column: 13}, EndPosition: Position{Line: 1, Column: 21}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 13}, EndPosition: Position{Line: 1, Column: 21}},
			},
		},
		{
			name:  "comments",
			input: "// my very very long comment",
			want: []Token{
				{Type: TT_COMMENT, Literal: "// my very very long comment", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 28}},
				{Type: TT_EOF, Literal: "0", BeginPosition: Position{Line: 1, Column: 1}, EndPosition: Position{Line: 1, Column: 28}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for _, want_tok := range tt.want {
				got_tok := l.NextToken()
				assert.Equal(t, want_tok, got_tok)
			}
		})
	}
}
