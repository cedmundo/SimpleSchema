package parser

import (
	"errors"
	"fmt"

	"github.com/cedmundo/SimpleSchema/lexer"
)

var (
	// punctuation by precedence
	punctPrec = map[int][]string{
		9: {"||"},
		8: {"&&"},
		7: {"|"},
		6: {"^"},
		5: {"&"},
		4: {"==", "!="},
		3: {"<", ">", "<=", ">="},
		2: {"+", "-"},
		1: {"*", "/", "%"},
	}
	maxPrec = 9
)

// ParseIdent tries to parse an identifier, returns error if token is not an id
func (p *Parser) ParseIdent() (Expr, error) {
	token, err := p.expect(lexer.Token{Tag: lexer.TokenTagWord})
	if err != nil {
		return nil, err
	}

	return &Ident{token}, nil
}

// ParseLiteral tries to parse a literal, returns error if token is not an literal
func (p *Parser) ParseLiteral() (Expr, error) {
	token, err := p.expect(
		lexer.Token{Tag: lexer.TokenTagBinInt},
		lexer.Token{Tag: lexer.TokenTagDecInt},
		lexer.Token{Tag: lexer.TokenTagOctInt},
		lexer.Token{Tag: lexer.TokenTagHexInt},
		lexer.Token{Tag: lexer.TokenTagString},
		lexer.Token{Tag: lexer.TokenTagFloat},
	)
	if err != nil {
		return nil, err
	}

	return &Literal{Token: token}, nil
}

func (p *Parser) parseField() (Decl, error) {
	field := &Field{}
	err := error(nil)

	// name (: type)? (= value)? (option block)?
	field.Name, err = p.ParseLookup()
	if err != nil {
		return nil, err
	}

	// type
	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: ":"})
	if err == nil {
		field.Type, err = p.ParseExpr()
		if err != nil {
			return nil, err
		}
	}

	// value
	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "="})
	if err == nil {
		field.Value, err = p.ParseExpr()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagWord, Value: "options"})
	if err == nil {
		// TODO(cedmundo): Use option block instead, leaving a block
		//	would lead to validation issues.
		block, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		field.Options = &OptionBlock{Block: block}
	}

	// end of line
	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagEOL})
	return field, err
}

func (p *Parser) parseBlock() (Block, error) {
	_, err := p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "{"})
	if err != nil {
		return Block{}, err
	}

	// Skip the end of line after "{" if needed
	_, _ = p.expect(lexer.Token{Tag: lexer.TokenTagEOL})

	decls := make([]Decl, 0)
	for {
		field, err := p.parseField()
		if err == nil {
			decls = append(decls, field)
			continue
		}

		_, err = p.expect(lexer.Token{Tag: lexer.TokenTagWord, Value: "options"})
		if err == nil {
			// TODO(cedmundo): Use option block instead, leaving a block
			//	would lead to validation issues.
			block, err := p.parseBlock()
			if err != nil {
				return Block{}, err
			}

			decls = append(decls, &OptionBlock{Block: block})
			continue
		}

		break
	}

	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "}"})
	return Block{Decls: decls}, err
}

// ParseStructDef tries to parse next expression as an struct definition
func (p *Parser) ParseStructDef() (Expr, error) {
	_, err := p.expect(lexer.Token{Tag: lexer.TokenTagWord, Value: "struct"})
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &StructDef{Block: block}, nil
}

// ParseUnionDef tries to parse next expression as an union definition
func (p *Parser) ParseUnionDef() (Expr, error) {
	_, err := p.expect(lexer.Token{Tag: lexer.TokenTagWord, Value: "union"})
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &UnionDef{Block: block}, nil
}

// ParseEnumDef tries to parse next expression as an enum definition
func (p *Parser) ParseEnumDef() (Expr, error) {
	_, err := p.expect(lexer.Token{Tag: lexer.TokenTagWord, Value: "enum"})
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &EnumDef{Block: block}, nil
}

// ParsePrototypeDef tries to parse next expression as proc prototype
func (p *Parser) ParsePrototypeDef() (Expr, error) {
	_, err := p.expect(lexer.Token{Tag: lexer.TokenTagWord, Value: "proc"})
	if err != nil {
		return nil, err
	}

	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "("})
	if err != nil {
		return nil, err
	}

	params := make([]Field, 0)
	for {
		param, err := p.ParseIdent()
		if err != nil {
			break
		}

		params = append(params, Field{Type: param})

		_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: ","})
		if err != nil {
			break
		}
	}

	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: ")"})
	if err != nil {
		return nil, err
	}

	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "->"})
	if err != nil {
		return nil, err
	}

	returnType, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	return &PrototypeDef{
		Params:     params,
		ReturnType: returnType,
	}, err
}

// ParseGroup tries to parse a grouping parenthesis
func (p *Parser) ParseGroup() (Expr, error) {
	_, err := p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "("})
	if err != nil {
		return nil, err
	}

	expr, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: ")"})
	return expr, err
}

// ParseAtom reads either an group, identifier or a literal
func (p *Parser) ParseAtom() (Expr, error) {
	atomParsers := []func() (Expr, error){
		p.ParseGroup,
		p.ParseStructDef,
		p.ParseUnionDef,
		p.ParseEnumDef,
		p.ParsePrototypeDef,
		p.ParseLiteral,
		p.ParseIdent,
	}
	for _, atomParser := range atomParsers {
		atom, err := atomParser()
		if err == nil {
			return atom, nil
		}
	}

	return nil, fmt.Errorf("%w was expecting atom", ErrUnexpectedToken)
}

// ParseLookup tries to parse a namespace lookup primitive (a.b)
func (p *Parser) ParseLookup() (Expr, error) {
	expr, err := p.ParseAtom()
	if err != nil {
		return nil, err
	}

	for {
		token, err := p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "."})
		if err != nil {
			break
		}

		right, err := p.ParseAtom()
		if err != nil {
			return nil, err
		}

		expr = &BinaryOp{
			Operator: token,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseArgs() ([]Expr, error) {
	args := make([]Expr, 0)
	_, err := p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "("})
	if err != nil {
		return args, err
	}

	for {
		expr, err := p.ParseExpr()
		if err != nil {
			break
		}

		args = append(args, expr)
		_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: ","})
		if err != nil {
			break
		}
	}

	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: ")"})
	if err != nil {
		return args, fmt.Errorf("%w: %w", err, ErrUnclosedParenthesis)
	}
	return args, nil
}

// ParseSubscript tries to parse calls and indexes
func (p *Parser) ParseSubscript() (Expr, error) {
	expr, err := p.ParseLookup()
	if err != nil {
		return nil, err
	}

	for {
		args, err := p.parseArgs()
		if err == nil {
			expr = &Call{
				Callee: expr,
				Args:   args,
			}
			continue
		} else if errors.Is(err, ErrUnclosedParenthesis) {
			return nil, err
		}

		_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "["})
		if err == nil {
			index, err := p.ParseExpr()
			if err != nil {
				return nil, err
			}

			expr = &Index{
				Base:  expr,
				Index: index,
			}

			_, err = p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: "]"})
			if err != nil {
				return nil, fmt.Errorf("%w: %w", err, ErrUnclosedSubscription)
			}
			continue
		}

		break
	}

	return expr, nil
}

// ParseUnary tries to parse unary expressions
func (p *Parser) ParseUnary() (Expr, error) {
	operator, err := p.expect(
		lexer.Token{Tag: lexer.TokenTagPunct, Value: "+"},
		lexer.Token{Tag: lexer.TokenTagPunct, Value: "-"},
		lexer.Token{Tag: lexer.TokenTagPunct, Value: "!"},
		lexer.Token{Tag: lexer.TokenTagPunct, Value: "~"},
		lexer.Token{Tag: lexer.TokenTagPunct, Value: "*"},
		lexer.Token{Tag: lexer.TokenTagPunct, Value: "&"},
	)
	if err == nil {
		expr, err := p.ParseUnary()
		if err != nil {
			return nil, err
		}

		return &UnaryOp{Operator: operator, Operand: expr}, nil
	}

	return p.ParseSubscript()
}

func (p *Parser) parseBinaryPrec(prec int) (Expr, error) {
	if prec == 0 {
		return p.ParseUnary()
	}

	expr, err := p.parseBinaryPrec(prec - 1)
	if err != nil {
		return nil, err
	}

	for {
		cont := false
		for _, punct := range punctPrec[prec] {
			token, err := p.expect(lexer.Token{Tag: lexer.TokenTagPunct, Value: punct})
			if err != nil && !errors.Is(err, ErrUnexpectedToken) {
				return nil, err
			} else if err == nil {
				right, err := p.parseBinaryPrec(prec - 1)
				if err != nil {
					return nil, err
				}

				expr = &BinaryOp{
					Operator: token,
					Left:     expr,
					Right:    right,
				}

				cont = true
				break // continue with next operator
			}
		}

		if !cont {
			break // cannot find more operators at this precedence
		}
	}

	return expr, nil
}

// ParseBinary parses common binary operators
func (p *Parser) ParseBinary() (Expr, error) {
	return p.parseBinaryPrec(maxPrec)
}

// ParseExpr parse next expression
func (p *Parser) ParseExpr() (Expr, error) {
	return p.ParseBinary()
}
