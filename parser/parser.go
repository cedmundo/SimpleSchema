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
func (p *Parser) Parse() error {
	return nil
}
