package parser_test

import (
	"testing"

	"github.com/cedmundo/SimpleSchema/lexer"
	"github.com/cedmundo/SimpleSchema/parser"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseAtom(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		expectedExpr parser.Expr
		expectedErr  error
	}{
		{
			name:  "parse literal atom",
			input: "10",
			expectedExpr: &parser.Literal{
				Token: lexer.Token{
					Tag:   lexer.TokenTagDecInt,
					Value: "10",
					Loc: lexer.Location{
						File: "parse literal atom",
						Col:  0,
						Row:  0,
					},
				},
			},
		},
		{
			name:  "parse identifier atom",
			input: "hello",
			expectedExpr: &parser.Ident{
				Token: lexer.Token{
					Tag:   lexer.TokenTagWord,
					Value: "hello",
					Loc: lexer.Location{
						File: "parse identifier atom",
						Col:  0,
						Row:  0,
					},
				},
			},
		},
		{
			name:  "parse group atom",
			input: "(hello)",
			expectedExpr: &parser.Ident{
				Token: lexer.Token{
					Tag:   lexer.TokenTagWord,
					Value: "hello",
					Loc: lexer.Location{
						File: "parse group atom",
						Col:  1,
						Row:  0,
					},
				},
			},
		},
		{
			name:        "fails to parse an atom because empty input",
			input:       "",
			expectedErr: parser.ErrUnexpectedToken,
		},
		{
			name:        "fails to parse an atom",
			input:       "+",
			expectedErr: parser.ErrUnexpectedToken,
		},
		{
			name:        "fails to parse a non-closed group atom",
			input:       "(a",
			expectedErr: parser.ErrUnexpectedToken,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewFromString(tt.name, tt.input)
			actualExpr, actualErr := p.ParseAtom()
			if tt.expectedErr != nil {
				require.ErrorIs(t, actualErr, tt.expectedErr)
				return
			}

			require.NoError(t, actualErr)
			require.Equal(t, tt.expectedExpr, actualExpr)
		})
	}
}

func TestParser_ParseLookup(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		expectedExpr parser.Expr
		expectedErr  error
	}{
		{
			name:  "parse descend",
			input: "id",
			expectedExpr: &parser.Ident{
				Token: lexer.Token{
					Tag:   lexer.TokenTagWord,
					Value: "id",
					Loc: lexer.Location{
						File: "parse descend",
						Col:  0,
						Row:  0,
					},
				},
			},
		},
		{
			name:  "parse single lookup",
			input: "id1.id2",
			expectedExpr: &parser.BinaryOp{
				Operator: lexer.Token{
					Tag: lexer.TokenTagPunct, Value: ".",
					Loc: lexer.Location{
						File: "parse single lookup",
						Col:  3,
						Row:  0,
					},
				},
				Left: &parser.Ident{
					Token: lexer.Token{
						Tag:   lexer.TokenTagWord,
						Value: "id1",
						Loc: lexer.Location{
							File: "parse single lookup",
							Col:  0,
							Row:  0,
						},
					},
				},
				Right: &parser.Ident{
					Token: lexer.Token{
						Tag:   lexer.TokenTagWord,
						Value: "id2",
						Loc: lexer.Location{
							File: "parse single lookup",
							Col:  4,
							Row:  0,
						},
					},
				},
			},
		},
		{
			name:  "parse multiple lookup",
			input: "id1.id2.id3",
			expectedExpr: &parser.BinaryOp{
				Operator: lexer.Token{
					Tag:   lexer.TokenTagPunct,
					Value: ".",
					Loc: lexer.Location{
						File: "parse multiple lookup",
						Col:  7,
						Row:  0,
					},
				},
				Left: &parser.BinaryOp{
					Operator: lexer.Token{
						Tag:   lexer.TokenTagPunct,
						Value: ".",
						Loc: lexer.Location{
							File: "parse multiple lookup",
							Col:  3,
							Row:  0,
						},
					},
					Left: &parser.Ident{
						Token: lexer.Token{
							Tag:   lexer.TokenTagWord,
							Value: "id1",
							Loc: lexer.Location{
								File: "parse multiple lookup",
								Col:  0,
								Row:  0,
							},
						},
					},
					Right: &parser.Ident{
						Token: lexer.Token{
							Tag:   lexer.TokenTagWord,
							Value: "id2",
							Loc: lexer.Location{
								File: "parse multiple lookup",
								Col:  4,
								Row:  0,
							},
						},
					},
				},
				Right: &parser.Ident{
					Token: lexer.Token{
						Tag:   lexer.TokenTagWord,
						Value: "id3",
						Loc: lexer.Location{
							File: "parse multiple lookup",
							Col:  8,
							Row:  0,
						},
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewFromString(tt.name, tt.input)
			actualExpr, actualErr := p.ParseLookup()
			if tt.expectedErr != nil {
				require.ErrorIs(t, actualErr, tt.expectedErr)
				return
			}

			require.NoError(t, actualErr)
			require.Equal(t, tt.expectedExpr, actualExpr)
		})
	}
}

func TestParser_ParseSubscript(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		expectedExpr parser.Expr
		expectedErr  error
	}{
		{
			name:  "parse call without args",
			input: "a()",
			expectedExpr: &parser.Call{
				Callee: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse call without args",
							Col:  0,
							Row:  0,
						},
						Value: "a",
					},
				},
				Args: []parser.Expr{},
			},
		},
		{
			name:  "parse call with single arg",
			input: "a(1)",
			expectedExpr: &parser.Call{
				Callee: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse call with single arg",
							Col:  0,
							Row:  0,
						},
						Value: "a",
					},
				},
				Args: []parser.Expr{
					&parser.Literal{
						Token: lexer.Token{
							Tag: lexer.TokenTagDecInt,
							Loc: lexer.Location{
								File: "parse call with single arg",
								Col:  2,
								Row:  0,
							},
							Value: "1",
						},
					},
				},
			},
		},
		{
			name:  "parse call with multiple args",
			input: "a(1, 2)",
			expectedExpr: &parser.Call{
				Callee: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse call with multiple args",
							Col:  0,
							Row:  0,
						},
						Value: "a",
					},
				},
				Args: []parser.Expr{
					&parser.Literal{
						Token: lexer.Token{
							Tag: lexer.TokenTagDecInt,
							Loc: lexer.Location{
								File: "parse call with multiple args",
								Col:  2,
								Row:  0,
							},
							Value: "1",
						},
					},
					&parser.Literal{
						Token: lexer.Token{
							Tag: lexer.TokenTagDecInt,
							Loc: lexer.Location{
								File: "parse call with multiple args",
								Col:  5,
								Row:  0,
							},
							Value: "2",
						},
					},
				},
			},
		},
		{
			name:  "parse indexing",
			input: "a[1]",
			expectedExpr: &parser.Index{
				Base: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse indexing",
							Col:  0,
							Row:  0,
						},
						Value: "a",
					},
				},
				Index: &parser.Literal{
					Token: lexer.Token{
						Tag: lexer.TokenTagDecInt,
						Loc: lexer.Location{
							File: "parse indexing",
							Col:  2,
							Row:  0,
						},
						Value: "1",
					},
				},
			},
		},
		{
			name:        "parse unclosed index",
			input:       "a[1",
			expectedErr: parser.ErrUnclosedSubscription,
		},
		{
			name:        "parse invalid index",
			input:       "a[]",
			expectedErr: parser.ErrUnexpectedToken,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewFromString(tt.name, tt.input)
			actualExpr, actualErr := p.ParseSubscript()
			if tt.expectedErr != nil {
				require.ErrorIs(t, actualErr, tt.expectedErr)
				return
			}

			require.NoError(t, actualErr)
			require.Equal(t, tt.expectedExpr, actualExpr)
		})
	}
}

func TestParser_ParseUnary(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		expectedExpr parser.Expr
		expectedErr  error
	}{
		{
			name:  "parse descend",
			input: "a",
			expectedExpr: &parser.Ident{
				Token: lexer.Token{
					Tag: lexer.TokenTagWord,
					Loc: lexer.Location{
						File: "parse descend",
						Row:  0,
						Col:  0,
					},
					Value: "a",
				},
			},
		},
		{
			name:  "parse single unary operator",
			input: "-a",
			expectedExpr: &parser.UnaryOp{
				Operator: lexer.Token{
					Tag: lexer.TokenTagPunct,
					Loc: lexer.Location{
						File: "parse single unary operator",
						Row:  0,
						Col:  0,
					},
					Value: "-",
				},
				Operand: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse single unary operator",
							Row:  0,
							Col:  1,
						},
						Value: "a",
					},
				},
			},
		},
		{
			name:  "parse nested unary operator",
			input: "+-a",
			expectedExpr: &parser.UnaryOp{
				Operator: lexer.Token{
					Tag: lexer.TokenTagPunct,
					Loc: lexer.Location{
						File: "parse nested unary operator",
						Row:  0,
						Col:  0,
					},
					Value: "+",
				},
				Operand: &parser.UnaryOp{
					Operator: lexer.Token{
						Tag: lexer.TokenTagPunct,
						Loc: lexer.Location{
							File: "parse nested unary operator",
							Row:  0,
							Col:  1,
						},
						Value: "-",
					},
					Operand: &parser.Ident{
						Token: lexer.Token{
							Tag: lexer.TokenTagWord,
							Loc: lexer.Location{
								File: "parse nested unary operator",
								Row:  0,
								Col:  2,
							},
							Value: "a",
						},
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewFromString(tt.name, tt.input)
			actualExpr, actualErr := p.ParseUnary()
			if tt.expectedErr != nil {
				require.ErrorIs(t, actualErr, tt.expectedErr)
				return
			}

			require.NoError(t, actualErr)
			require.Equal(t, tt.expectedExpr, actualExpr)
		})
	}
}

func TestParser_ParseBinary(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		expectedExpr parser.Expr
		expectedErr  error
	}{
		{
			name:  "parse descend",
			input: "a",
			expectedExpr: &parser.Ident{
				Token: lexer.Token{
					Tag: lexer.TokenTagWord,
					Loc: lexer.Location{
						File: "parse descend",
						Row:  0,
						Col:  0,
					},
					Value: "a",
				},
			},
		},
		{
			name:  "parse single binary operator",
			input: "a + b",
			expectedExpr: &parser.BinaryOp{
				Operator: lexer.Token{
					Tag: lexer.TokenTagPunct,
					Loc: lexer.Location{
						File: "parse single binary operator",
						Row:  0,
						Col:  2,
					},
					Value: "+",
				},
				Left: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse single binary operator",
							Row:  0,
							Col:  0,
						},
						Value: "a",
					},
				},
				Right: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse single binary operator",
							Row:  0,
							Col:  4,
						},
						Value: "b",
					},
				},
			},
		},
		{
			name:  "parse multiple binary operators",
			input: "a + b + c",
			expectedExpr: &parser.BinaryOp{
				Operator: lexer.Token{
					Tag: lexer.TokenTagPunct,
					Loc: lexer.Location{
						File: "parse multiple binary operators",
						Row:  0,
						Col:  6,
					},
					Value: "+",
				},
				Left: &parser.BinaryOp{
					Operator: lexer.Token{
						Tag: lexer.TokenTagPunct,
						Loc: lexer.Location{
							File: "parse multiple binary operators",
							Row:  0,
							Col:  2,
						},
						Value: "+",
					},
					Left: &parser.Ident{
						Token: lexer.Token{
							Tag: lexer.TokenTagWord,
							Loc: lexer.Location{
								File: "parse multiple binary operators",
								Row:  0,
								Col:  0,
							},
							Value: "a",
						},
					},
					Right: &parser.Ident{
						Token: lexer.Token{
							Tag: lexer.TokenTagWord,
							Loc: lexer.Location{
								File: "parse multiple binary operators",
								Row:  0,
								Col:  4,
							},
							Value: "b",
						},
					},
				},
				Right: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse multiple binary operators",
							Row:  0,
							Col:  8,
						},
						Value: "c",
					},
				},
			},
		},
		{
			name:  "parse grouped binary operators",
			input: "a + (b + c)",
			expectedExpr: &parser.BinaryOp{
				Operator: lexer.Token{
					Tag: lexer.TokenTagPunct,
					Loc: lexer.Location{
						File: "parse grouped binary operators",
						Row:  0,
						Col:  2,
					},
					Value: "+",
				},
				Left: &parser.Ident{
					Token: lexer.Token{
						Tag: lexer.TokenTagWord,
						Loc: lexer.Location{
							File: "parse grouped binary operators",
							Row:  0,
							Col:  0,
						},
						Value: "a",
					},
				},
				Right: &parser.BinaryOp{
					Operator: lexer.Token{
						Tag: lexer.TokenTagPunct,
						Loc: lexer.Location{
							File: "parse grouped binary operators",
							Row:  0,
							Col:  7,
						},
						Value: "+",
					},
					Left: &parser.Ident{
						Token: lexer.Token{
							Tag: lexer.TokenTagWord,
							Loc: lexer.Location{
								File: "parse grouped binary operators",
								Row:  0,
								Col:  5,
							},
							Value: "b",
						},
					},
					Right: &parser.Ident{
						Token: lexer.Token{
							Tag: lexer.TokenTagWord,
							Loc: lexer.Location{
								File: "parse grouped binary operators",
								Row:  0,
								Col:  9,
							},
							Value: "c",
						},
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewFromString(tt.name, tt.input)
			actualExpr, actualErr := p.ParseBinary()
			if tt.expectedErr != nil {
				require.ErrorIs(t, actualErr, tt.expectedErr)
				return
			}

			require.NoError(t, actualErr)
			require.Equal(t, tt.expectedExpr, actualExpr)
		})
	}
}
