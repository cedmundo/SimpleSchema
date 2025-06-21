package lexer

import (
	"fmt"
)

// Location is a token coordinate, relative to build path
type Location struct {
	File string
	Row  int
	Col  int
}

// TokenTag classifies a token
type TokenTag int

// Token is a text span with a tag
type Token struct {
	Tag   TokenTag
	Loc   Location
	Value string
}

const (
	TokenTagEOF     TokenTag = iota // TokenTagEOF end of file
	TokenTagEOL                     // TokenTagEOL end of line
	TokenTagComment                 // TokenTagComment only single-line comments at the moment
	TokenTagDecInt                  // TokenTagDecInt a decimal integer number
	TokenTagBinInt                  // TokenTagBinInt a binary integer number
	TokenTagOctInt                  // TokenTagOctInt a octal integer number
	TokenTagHexInt                  // TokenTagHexInt a hexadecimal integer number
	TokenTagFloat                   // TokenTagFloat a decimal floating point number
	TokenTagString                  // TokenTagString a string literal
	TokenTagWord                    // TokenTagWord both ids and keywords
	TokenTagPunct                   // TokenTagPunct any punctuation symbol
)

// String returns a standard file coordinate format
func (l Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.File, l.Row, l.Col)
}

// String returns debuggable info about the token
func (t Token) String() string {
	switch t.Tag {
	case TokenTagEOF:
		return "`EOF`"
	case TokenTagEOL:
		return "`EOL`"
	case TokenTagComment:
		return fmt.Sprintf("`COMMENT '%s'`", t.Value)
	case TokenTagDecInt, TokenTagBinInt,
		TokenTagOctInt, TokenTagHexInt:
		return fmt.Sprintf("`INT '%s'`", t.Value)
	case TokenTagFloat:
		return fmt.Sprintf("`FLOAT '%s'`", t.Value)
	case TokenTagString:
		return fmt.Sprintf("`STRING '%s'`", t.Value)
	case TokenTagWord:
		return fmt.Sprintf("`WORD '%s'`", t.Value)
	case TokenTagPunct:
		return fmt.Sprintf("`PUNCT '%s'`", t.Value)
	}
	panic("unreachable code: unhandled tag in Token.String()")
}

func (t Token) GetErrorf(msg string, args ...any) error {
	return fmt.Errorf("%s:%d:%d: %s", t.Loc.File, t.Loc.Row, t.Loc.Col, fmt.Sprintf(msg, args...))
}
