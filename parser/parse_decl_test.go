package parser_test

import (
	"testing"

	"github.com/cedmundo/SimpleSchema/lexer"
	"github.com/cedmundo/SimpleSchema/parser"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseDecl(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		expectedDecl parser.Decl
		expectedErr  error
	}{
		{
			name:  "parse module decl",
			input: "module name;",
			expectedDecl: &parser.ModuleDecl{
				Name: &parser.Ident{Token: lexer.Token{
					Tag:   lexer.TokenTagWord,
					Loc:   lexer.Location{File: "parse module decl", Row: 0, Col: 7},
					Value: "name",
				}},
			},
		},
		{
			name:  "parse type decl",
			input: "type name int;",
			expectedDecl: &parser.TypeDecl{
				Name: &parser.Ident{Token: lexer.Token{
					Tag:   lexer.TokenTagWord,
					Loc:   lexer.Location{File: "parse type decl", Row: 0, Col: 5},
					Value: "name",
				}},
				Type: &parser.Ident{Token: lexer.Token{
					Tag:   lexer.TokenTagWord,
					Loc:   lexer.Location{File: "parse type decl", Row: 0, Col: 10},
					Value: "int",
				}},
			},
		},
		{
			name:  "parse proc decl",
			input: "proc name () -> void;",
			expectedDecl: &parser.ProcDecl{
				Name: &parser.Ident{Token: lexer.Token{
					Tag:   lexer.TokenTagWord,
					Loc:   lexer.Location{File: "parse proc decl", Row: 0, Col: 5},
					Value: "name",
				}},
				Type: &parser.PrototypeDef{
					Params: []parser.Field{},
					ReturnType: &parser.Ident{Token: lexer.Token{
						Tag:   lexer.TokenTagWord,
						Loc:   lexer.Location{File: "parse proc decl", Row: 0, Col: 16},
						Value: "void",
					}},
				},
			},
		},
		{
			name:  "parse proc decl with args",
			input: "proc name (s: int) -> void;",
			expectedDecl: &parser.ProcDecl{
				Name: &parser.Ident{Token: lexer.Token{
					Tag:   lexer.TokenTagWord,
					Loc:   lexer.Location{File: "parse proc decl with args", Row: 0, Col: 5},
					Value: "name",
				}},
				Type: &parser.PrototypeDef{
					Params: []parser.Field{
						{
							Name: &parser.Ident{Token: lexer.Token{
								Tag:   lexer.TokenTagWord,
								Loc:   lexer.Location{File: "parse proc decl with args", Row: 0, Col: 11},
								Value: "s",
							}},
							Type: &parser.Ident{Token: lexer.Token{
								Tag:   lexer.TokenTagWord,
								Loc:   lexer.Location{File: "parse proc decl with args", Row: 0, Col: 14},
								Value: "int",
							}},
						},
					},
					ReturnType: &parser.Ident{Token: lexer.Token{
						Tag:   lexer.TokenTagWord,
						Loc:   lexer.Location{File: "parse proc decl with args", Row: 0, Col: 22},
						Value: "void",
					}},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewFromString(tt.name, tt.input)
			actualExpr, actualErr := p.ParseDecl()
			if tt.expectedErr != nil {
				require.ErrorIs(t, actualErr, tt.expectedErr)
				return
			}

			require.NoError(t, actualErr)
			require.Equal(t, tt.expectedDecl, actualExpr)
		})
	}
}
