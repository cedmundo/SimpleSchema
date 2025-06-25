package parser

import "github.com/cedmundo/SimpleSchema/lexer"

// ParseDecl parses either type proc or module
func (p *Parser) ParseDecl() (Decl, error) {
	obj, err := p.expect(
		lexer.Token{Tag: lexer.TokenTagWord, Value: "module"},
		lexer.Token{Tag: lexer.TokenTagWord, Value: "type"},
		lexer.Token{Tag: lexer.TokenTagWord, Value: "proc"},
	)
	if err != nil {
		return nil, err
	}

	name, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}

	var expr Expr
	if obj.Value == "type" {
		expr, err = p.ParseExpr()
		if err != nil {
			return nil, err
		}
	} else if obj.Value == "proc" {
		expr, err = p.parseArgsWithReturnType()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.expect(lexer.Token{Tag: lexer.TokenTagEOL})
	if err != nil {
		return nil, err
	}

	if obj.Value == "module" {
		return &ModuleDecl{Name: name}, nil
	}

	if obj.Value == "proc" {
		return &ProcDecl{Name: name, Type: expr}, nil
	}

	return &TypeDecl{Name: name, Type: expr}, nil
}

// ParseAnnotatedDecl annotations followed by types
func (p *Parser) ParseAnnotatedDecl() (Decl, error) {
	annotations, err := p.parseAnnotations()
	if err != nil {
		return nil, err
	}

	decl, err := p.ParseDecl()
	if err != nil {
		return nil, err
	}

	return &AnnotatedDecl{
		Annotations: annotations,
		Decl:        decl,
	}, nil
}
