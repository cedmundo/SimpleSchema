// Package parser makes AST nodes representing the data of a schema file
package parser

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cedmundo/SimpleSchema/lexer"
)

var (
	ErrUnexpectedToken      = errors.New("unexpected token")
	ErrUnclosedParenthesis  = errors.New("unclosed parenthesis")
	ErrUnclosedSubscription = errors.New("unclosed subscription")
)

// Parser handle a single file parsing
type Parser struct {
	lex *lexer.Lexer
}

// New returns a new parser using only a filename and a rune reader
func New(filename string, r io.RuneReader) *Parser {
	return &Parser{lex: lexer.New(filename, r)}
}

// NewFromString returns new parser using a string as content
func NewFromString(filename, content string) *Parser {
	return New(filename, strings.NewReader(content))
}

func (p *Parser) expect(anyOf ...lexer.Token) (lexer.Token, error) {
	token, err := p.lex.Read()
	if err != nil {
		return token, err
	}

	for _, matching := range anyOf {
		matchesTag := token.Tag == matching.Tag
		matchesValue := matching.Value == "" || (matching.Value != "" && matching.Value == token.Value)
		if matchesTag && matchesValue {
			return token, nil
		}
	}

	err = p.lex.Unread(token)
	if err != nil {
		return token, err
	}

	return token, fmt.Errorf("%w `%s`", ErrUnexpectedToken, token.Value)
}

// Parse reads the entire file and descends on each rule to make an AST
func (p *Parser) Parse() (*Schema, error) {
	// Skip starting end of lines
	_, _ = p.expect(lexer.Token{Tag: lexer.TokenTagEOL})

	decls := make([]Decl, 0)
	for {
		annotatedDecl, err := p.ParseAnnotatedDecl()
		if err == nil {
			decls = append(decls, annotatedDecl)
			continue
		}

		decl, err := p.ParseDecl()
		if err == nil {
			decls = append(decls, decl)
			continue
		}

		break
	}

	// Skip trailing end of lines and EOF
	_, _ = p.expect(lexer.Token{Tag: lexer.TokenTagEOL})
	_, err := p.expect(lexer.Token{Tag: lexer.TokenTagEOF})

	return &Schema{
		Decls: decls,
	}, err
}
