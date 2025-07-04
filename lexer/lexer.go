// Package lexer handles token identification and manipulation at scanner level.
package lexer

import (
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

var (
	// ErrCannotTokenize indicates that the lexer cannot process the input stream into tokens due to an error.
	ErrCannotTokenize = errors.New("cannot tokenize stream")

	// ErrInvalidCharacter represents an error that occurs when an unexpected or invalid character is encountered in the input.
	ErrInvalidCharacter = errors.New("invalid character")

	// ErrMalformedFloatLiteral represents an error that occurs when a floating-point literal is improperly formatted.
	ErrMalformedFloatLiteral = errors.New("malformed floating literal")

	// ErrUnterminatedStringLiteral represents an error that occurs when a string literal is not properly closed before the end of the line.
	ErrUnterminatedStringLiteral = errors.New("unterminated string literal")

	// ErrMalformedEscapeSequence indicates that an escape sequence in a string or character literal is not recognized or properly formatted.
	ErrMalformedEscapeSequence = errors.New("malformed escape sequence")

	// ErrAlreadyUnread indicates an attempt to mark a token as unread when there is already an existing unread token.
	ErrAlreadyUnread = errors.New("token is already unread")

	// ErrUnbalancedGroup indicates that the grouping is not valid (there are more closes than opens)
	ErrUnbalancedGroup = errors.New("unbalanced group")

	punctuations = []string{
		"(", ")", "[", "]", "{", "}", ",", ".", ":", "=", "+", "-", "*", "/", "%",
		">", "<", "^", "~", "!", "|", "&", ":=", "==", "!=", ">=", "<=",
		">>", "<<", "&&", "||", "=>", "->", "[[", "]]",
	}
)

// Lexer is responsible for converting a sequence of characters into a sequence of tokens for parser consumption.
type Lexer struct {
	startLoc Location
	endLoc   Location
	current  rune
	consumed bool
	reader   io.RuneReader
	unread   *Token
	group    int
}

type tryReadFn func() (Token, error)

// New returns a new lexer using a rune reader
func New(file string, reader io.RuneReader) *Lexer {
	loc := Location{File: file}
	return &Lexer{
		reader:   reader,
		startLoc: loc,
		endLoc:   loc,
	}
}

// NewFromString returns a lexer using a string content
func NewFromString(file, content string) *Lexer {
	return New(file, strings.NewReader(content))
}

func (l *Lexer) advanceRune() (err error) {
	l.current, _, err = l.reader.ReadRune()
	if errors.Is(err, io.EOF) {
		l.consumed = true
		return nil
	}

	l.endLoc.Col += 1
	if l.current == '\n' {
		l.endLoc.Col = 0
		l.endLoc.Row += 1
	}
	return err
}

func (l *Lexer) skipSpaces() error {
	for l.current == ' ' || l.current == '\t' || (l.group != 0 && unicode.IsSpace(l.current)) {
		l.startLoc.Col += 1
		if l.current == '\n' {
			l.startLoc.Col = 0
			l.startLoc.Row += 1
		}

		err := l.advanceRune()
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Lexer) tryReadEOF() (Token, error) {
	if l.consumed {
		return Token{Tag: TokenTagEOF, Loc: l.startLoc}, nil
	}

	return Token{}, ErrInvalidCharacter
}

func (l *Lexer) tryReadEOL() (Token, error) {
	if l.group != 0 || (l.current != '\n' && l.current != ';') {
		return Token{}, ErrInvalidCharacter
	}

	start := l.startLoc
	for l.current == '\n' || l.current == ';' {
		err := l.advanceRune()
		if err != nil {
			return Token{}, err
		}
	}

	return Token{
		Tag: TokenTagEOL,
		Loc: start,
	}, nil
}

func (l *Lexer) tryReadComment() (Token, error) {
	if l.current != '#' {
		return Token{}, ErrInvalidCharacter
	}

	start := l.startLoc
	value := strings.Builder{}

	for l.current != '\n' && l.current != 0 {
		value.WriteRune(l.current)
		err := l.advanceRune()
		if err != nil {
			return Token{}, err
		}
	}

	// Read new line (so it skips the token further down)
	if l.current == '\n' {
		err := l.advanceRune()
		if err != nil {
			return Token{}, err
		}
	}

	return Token{
		Tag:   TokenTagComment,
		Loc:   start,
		Value: value.String(),
	}, nil
}

func (l *Lexer) tryReadNumber() (Token, error) {
	if !isDigitOfBase(l.current, TokenTagDecInt) {
		return Token{}, ErrInvalidCharacter
	}

	tag := TokenTagDecInt
	start := l.startLoc
	haveExp := false
	value := strings.Builder{}

	if l.current == '0' {
		skip := true
		err := l.advanceRune()
		if err != nil {
			return Token{}, err
		}

		switch l.current {
		case 'b':
			tag = TokenTagBinInt
		case 'o':
			tag = TokenTagOctInt
		case 'x':
			tag = TokenTagHexInt
		case '.':
			tag = TokenTagFloat
			value.WriteString("0.")
		default:
			skip = false
			value.WriteRune('0')
		}

		if skip {
			err = l.advanceRune()
			if err != nil {
				return Token{}, err
			}
		}
	}

	for {
		for isDigitOfBase(l.current, tag) {
			value.WriteRune(l.current)
			err := l.advanceRune()
			if err != nil {
				return Token{}, err
			}
		}

		if l.current == '.' && tag == TokenTagDecInt {
			value.WriteRune(l.current)
			err := l.advanceRune()
			if err != nil {
				return Token{}, err
			}

			tag = TokenTagFloat
			continue
		} else if l.current == '.' && tag != TokenTagDecInt {
			return Token{}, ErrMalformedFloatLiteral
		}

		if l.current == 'e' &&
			!haveExp &&
			(tag == TokenTagDecInt || tag == TokenTagFloat) {
			haveExp = true
			value.WriteRune(l.current)
			err := l.advanceRune()
			if err != nil {
				return Token{}, err
			}

			// exponent sign
			if l.current == '-' || l.current == '+' {
				value.WriteRune(l.current)
				err := l.advanceRune()
				if err != nil {
					return Token{}, err
				}
			}

			if !isDigitOfBase(l.current, TokenTagFloat) {
				return Token{}, ErrMalformedFloatLiteral
			}

			continue
		}

		break
	}

	return Token{
		Tag:   tag,
		Loc:   start,
		Value: value.String(),
	}, nil

}

func (l *Lexer) tryReadString() (Token, error) {
	if l.current != '"' {
		return Token{}, ErrInvalidCharacter
	}

	start := l.startLoc
	value := strings.Builder{}

	for l.current != '\n' && l.current != 0 {
		err := l.advanceRune()
		if err != nil {
			return Token{}, err
		}

		if l.current == '\\' {
			err = l.decodeEscapeSequence(&value)
			if err != nil {
				return Token{}, err
			}
		}

		if l.current == '"' {
			break
		}

		value.WriteRune(l.current)
	}

	if l.current != '"' {
		return Token{}, ErrUnterminatedStringLiteral
	}

	err := l.advanceRune()
	if err != nil {
		return Token{}, err
	}

	return Token{
		Tag:   TokenTagString,
		Loc:   start,
		Value: value.String(),
	}, nil
}

func (l *Lexer) decodeEscapeSequence(value *strings.Builder) error {
	// must already read first '\'
	err := l.advanceRune()
	if err != nil {
		return err
	}

	switch l.current {
	case 'a':
		value.WriteRune('\a')
	case 'b':
		value.WriteRune('\b')
	case 'f':
		value.WriteRune('\f')
	case 'n':
		value.WriteRune('\n')
	case 'r':
		value.WriteRune('\r')
	case 't':
		value.WriteRune('\t')
	case 'v':
		value.WriteRune('\v')
	case '\\':
		value.WriteRune('\\')
	case '\'':
		value.WriteRune('\'')
	case '"':
		value.WriteRune('"')
	default:
		// x, u, U
		takeNext := 0
		if l.current == 'x' {
			takeNext = 2
		} else if l.current == 'u' {
			takeNext = 4
		} else if l.current == 'U' {
			takeNext = 8
		} else {
			return ErrMalformedEscapeSequence
		}

		charDigits := strings.Builder{}
		for range takeNext {
			err = l.advanceRune()
			if err != nil {
				return err
			}

			if !isDigitOfBase(l.current, TokenTagHexInt) {
				return ErrMalformedEscapeSequence
			}

			charDigits.WriteRune(l.current)
		}

		charValue, _ := strconv.ParseInt(charDigits.String(), 16, 64)
		if takeNext == 2 {
			value.WriteByte(byte(charValue))
		} else {
			value.WriteRune(rune(charValue))
		}
	}

	// leave cursor just after escape
	err = l.advanceRune()
	if err != nil {
		return err
	}

	return nil
}

func (l *Lexer) tryReadWord() (Token, error) {
	if !unicode.IsLetter(l.current) && l.current != '_' {
		return Token{}, ErrInvalidCharacter
	}

	start := l.startLoc
	value := strings.Builder{}

	for unicode.IsLetter(l.current) || unicode.IsDigit(l.current) || l.current == '_' {
		value.WriteRune(l.current)
		err := l.advanceRune()
		if err != nil {
			return Token{}, err
		}
	}

	l.endLoc.Col = start.Col + value.Len()
	return Token{
		Tag:   TokenTagWord,
		Loc:   start,
		Value: value.String(),
	}, nil
}

func (l *Lexer) tryReadPunct() (Token, error) {
	value := strings.Builder{}
	start := l.startLoc
	for {
		if !slices.Contains(punctuations, value.String()+string(l.current)) {
			break
		}

		value.WriteRune(l.current)
		err := l.advanceRune()
		if err != nil {
			return Token{}, err
		}
	}

	if value.Len() == 0 {
		return Token{}, ErrInvalidCharacter
	}

	l.endLoc.Col = start.Col + value.Len()
	return Token{
		Tag:   TokenTagPunct,
		Loc:   start,
		Value: value.String(),
	}, nil
}

// Read scans the input and returns the next token or an error if an invalid character is encountered.
// It prioritizes unread tokens, attempts to classify current input, and skips spaces as necessary.
func (l *Lexer) Read() (Token, error) {
	if l.unread != nil {
		token := *l.unread
		l.unread = nil
		return token, nil
	}

	var token Token
	var err error
	if l.current == 0 && !l.consumed {
		err = l.advanceRune()
		if err != nil {
			return token, errors.Join(err, token.GetErrorf("cannot read first character"))
		}
	}

	err = l.skipSpaces()
	if err != nil {
		return token, errors.Join(err, token.GetErrorf("cannot skip spaces"))
	}

	defer func() {
		l.startLoc = l.endLoc
	}()

	// order is important
	classifiers := []tryReadFn{
		l.tryReadEOF,
		l.tryReadEOL,
		l.tryReadComment,
		l.tryReadNumber,
		l.tryReadString,
		l.tryReadWord,
		l.tryReadPunct,
	}
	for _, classifier := range classifiers {
		token, err = classifier()
		if err != nil && !errors.Is(err, ErrInvalidCharacter) {
			return token, err
		} else if err == nil {
			return token, nil
		}
	}

	token = Token{}
	return token, errors.Join(ErrCannotTokenize, ErrInvalidCharacter, token.GetErrorf("invalid character: %q", l.current))
}

// Unread attempts to set the given token as the unread token in the lexer. Returns an error if there is already an unread token.
func (l *Lexer) Unread(token Token) error {
	if l.unread != nil {
		return ErrAlreadyUnread
	}

	l.unread = &token
	return nil
}

// PushGroup pushes a group so lexer will ignore new lines
func (l *Lexer) PushGroup() {
	l.group += 1
}

// PushGroup pops a group so lexer will stop to ignore new lines
func (l *Lexer) PopGroup() error {
	if l.group <= 0 {
		return ErrUnbalancedGroup
	}

	l.group -= 1
	return nil
}

func isDigitOfBase(r rune, tag TokenTag) bool {
	switch tag {
	case TokenTagBinInt:
		return r == '0' || r == '1'
	case TokenTagOctInt:
		return r >= '0' && r <= '7'
	case TokenTagDecInt, TokenTagFloat:
		return r >= '0' && r <= '9'
	case TokenTagHexInt:
		return (r >= '0' && r <= '9') ||
			(r >= 'a' && r <= 'f') ||
			(r >= 'A' && r <= 'F')
	default:
		panic("unreachable code: invalid numeric base")
	}
}
