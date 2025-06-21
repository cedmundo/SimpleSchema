package lexer_test

import (
	"testing"

	"github.com/cedmundo/SimpleSchema/lexer"
	"github.com/stretchr/testify/require"
)

func TestLexer_SingleScans(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedTokens []lexer.Token
		expectedError  error
	}{
		{
			name:  "lex EOF",
			input: "",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex EOF", Row: 0, Col: 0}},
			},
		},
		{
			name:  "lex spaces",
			input: "    ",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex spaces", Row: 0, Col: 4}},
			},
		},
		{
			name:  "lex new lines",
			input: "  \n",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagEOL, Loc: lexer.Location{File: "lex new lines", Row: 0, Col: 2}},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex new lines", Row: 1, Col: 0}},
			},
		},
		{
			name:  "comments",
			input: "# a comment\n",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagComment, Loc: lexer.Location{File: "comments", Row: 0, Col: 0}, Value: "# a comment"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "comments", Row: 1, Col: 0}},
			},
		},
		{
			name:  "lex int zero",
			input: "0",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagDecInt, Loc: lexer.Location{File: "lex int zero", Row: 0, Col: 0}, Value: "0"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex int zero", Row: 0, Col: 1}},
			},
		},
		{
			name:  "lex dec int",
			input: "123",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagDecInt, Loc: lexer.Location{File: "lex dec int", Row: 0, Col: 0}, Value: "123"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex dec int", Row: 0, Col: 3}},
			},
		},
		{
			name:  "lex bin int",
			input: "0b1010",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagBinInt, Loc: lexer.Location{File: "lex bin int", Row: 0, Col: 0}, Value: "1010"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex bin int", Row: 0, Col: 6}},
			},
		},
		{
			name:  "lex oct int",
			input: "0o766",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagOctInt, Loc: lexer.Location{File: "lex oct int", Row: 0, Col: 0}, Value: "766"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex oct int", Row: 0, Col: 5}},
			},
		},
		{
			name:  "lex hex int",
			input: "0xF0F0",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagHexInt, Loc: lexer.Location{File: "lex hex int", Row: 0, Col: 0}, Value: "F0F0"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex hex int", Row: 0, Col: 6}},
			},
		},
		{
			name:  "lex float one",
			input: "1.0",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagFloat, Loc: lexer.Location{File: "lex float one", Row: 0, Col: 0}, Value: "1.0"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex float one", Row: 0, Col: 3}},
			},
		},
		{
			name:  "lex float zero",
			input: "0.0",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagFloat, Loc: lexer.Location{File: "lex float zero", Row: 0, Col: 0}, Value: "0.0"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex float zero", Row: 0, Col: 3}},
			},
		},
		{
			name:  "lex float with exp",
			input: "1.0e5",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagFloat, Loc: lexer.Location{File: "lex float with exp", Row: 0, Col: 0}, Value: "1.0e5"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex float with exp", Row: 0, Col: 5}},
			},
		},
		{
			name:  "lex float with neg exp",
			input: "1.0e-5",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagFloat, Loc: lexer.Location{File: "lex float with neg exp", Row: 0, Col: 0}, Value: "1.0e-5"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex float with neg exp", Row: 0, Col: 6}},
			},
		},
		{
			name:  "lex float with pos exp",
			input: "1.0e+5",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagFloat, Loc: lexer.Location{File: "lex float with pos exp", Row: 0, Col: 0}, Value: "1.0e+5"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex float with pos exp", Row: 0, Col: 6}},
			},
		},
		{
			name:  "lex int with exp",
			input: "1e4",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagDecInt, Loc: lexer.Location{File: "lex int with exp", Row: 0, Col: 0}, Value: "1e4"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex int with exp", Row: 0, Col: 3}},
			},
		},
		{
			name:          "lex malformed float",
			input:         "1..3",
			expectedError: lexer.ErrMalformedFloatLiteral,
		},
		{
			name:          "lex malformed exp",
			input:         "1.0e",
			expectedError: lexer.ErrMalformedFloatLiteral,
		},
		{
			name:          "lex malformed float with invalid base",
			input:         "0x1.2",
			expectedError: lexer.ErrMalformedFloatLiteral,
		},
		{
			name:  "lex empty string",
			input: "\"\"",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagString, Loc: lexer.Location{File: "lex empty string", Row: 0, Col: 0}, Value: ""},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex empty string", Row: 0, Col: 2}},
			},
		},
		{
			name:  "lex non-empty string",
			input: "\"hello\"",
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagString, Loc: lexer.Location{File: "lex non-empty string", Row: 0, Col: 0}, Value: "hello"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex non-empty string", Row: 0, Col: 7}},
			},
		},
		{
			name:  "lex escaped string",
			input: `"\tTABS"`,
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagString, Loc: lexer.Location{File: "lex escaped string", Row: 0, Col: 0}, Value: "\tTABS"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex escaped string", Row: 0, Col: 8}},
			},
		},
		{
			name:  "lex byte-escaped string",
			input: `"\xC0"`,
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagString, Loc: lexer.Location{File: "lex byte-escaped string", Row: 0, Col: 0}, Value: "\xC0"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex byte-escaped string", Row: 0, Col: 6}},
			},
		},
		{
			name:  "lex unicode-escaped string",
			input: `"\u3071"`,
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagString, Loc: lexer.Location{File: "lex unicode-escaped string", Row: 0, Col: 0}, Value: "\u3071"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex unicode-escaped string", Row: 0, Col: 8}},
			},
		},
		{
			name:  "lex large unicode-escaped string",
			input: `"\U0001F617"`,
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagString, Loc: lexer.Location{File: "lex large unicode-escaped string", Row: 0, Col: 0}, Value: "\U0001F617"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex large unicode-escaped string", Row: 0, Col: 12}},
			},
		},
		{
			name:          "lex unterminated string",
			input:         `"a`,
			expectedError: lexer.ErrMalformedEscapeSequence,
		},
		{
			name:          "lex invalid escape sequence",
			input:         `"\M"`,
			expectedError: lexer.ErrMalformedEscapeSequence,
		},
		{
			name:          "lex malformed escape sequence",
			input:         `"\xNO"`,
			expectedError: lexer.ErrMalformedEscapeSequence,
		},
		{
			name:  "lex word",
			input: `_hello_world_10`,
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagWord, Loc: lexer.Location{File: "lex word", Row: 0, Col: 0}, Value: "_hello_world_10"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex word", Row: 0, Col: 15}},
			},
		},
		{
			name:  "lex punct",
			input: `:=`,
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagPunct, Loc: lexer.Location{File: "lex punct", Row: 0, Col: 0}, Value: ":="},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex punct", Row: 0, Col: 2}},
			},
		},
		{
			name:  "lex punct with juxtaposition",
			input: `+(`,
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagPunct, Loc: lexer.Location{File: "lex punct with juxtaposition", Row: 0, Col: 0}, Value: "+"},
				{Tag: lexer.TokenTagPunct, Loc: lexer.Location{File: "lex punct with juxtaposition", Row: 0, Col: 1}, Value: "("},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex punct with juxtaposition", Row: 0, Col: 2}},
			},
		},
		{
			name:  "lex single character word",
			input: `a+`,
			expectedTokens: []lexer.Token{
				{Tag: lexer.TokenTagWord, Loc: lexer.Location{File: "lex single character word", Row: 0, Col: 0}, Value: "a"},
				{Tag: lexer.TokenTagPunct, Loc: lexer.Location{File: "lex single character word", Row: 0, Col: 1}, Value: "+"},
				{Tag: lexer.TokenTagEOF, Loc: lexer.Location{File: "lex single character word", Row: 0, Col: 2}},
			},
		},
		{
			name:          "lex unknown symbol",
			input:         `°`,
			expectedError: lexer.ErrInvalidCharacter,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.NewFromString(tt.name, tt.input)
			for _, expectedToken := range tt.expectedTokens {
				actualToken, err := lex.Read()
				if tt.expectedError != nil {
					require.ErrorIs(t, err, tt.expectedError)
					return
				}

				require.NoError(t, err)
				require.Equal(t, expectedToken, actualToken)
			}

		})
	}
}
