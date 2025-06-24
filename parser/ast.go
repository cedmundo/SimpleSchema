package parser

import "github.com/cedmundo/SimpleSchema/lexer"

// Decl represents any declaration such types, fields and options
type Decl interface {
	decl()
}

// Expr represents any expressions, including literal, binary and unary operators
type Expr interface {
	expr()
}

// Literal represents any plain data in text representation
type Literal struct {
	Token lexer.Token
}

func (l *Literal) expr() {}

// Ident represents an identifier
type Ident struct {
	Token lexer.Token
}

func (i *Ident) expr() {}

// Call represents a call expression (callee(args))
type Call struct {
	Callee Expr
	Args   []Expr
}

func (ca *Call) expr() {}

// Index represents a selection expression (base[index])
type Index struct {
	Base  Expr
	Index Expr
}

func (in *Index) expr() {}

// UnaryOp represents any prefix and suffix operation
type UnaryOp struct {
	Operator lexer.Token
	Operand  Expr
}

func (uo *UnaryOp) expr() {}

// BinaryOp represents any infix operation
type BinaryOp struct {
	Operator lexer.Token
	Left     Expr
	Right    Expr
}

func (bo *BinaryOp) expr() {}

// StructDef represents the definition of a struct body(struct { fields ... })
type StructDef struct {
	Block Block
}

func (sd *StructDef) expr() {}

// UnionDef represents the definition of a union body(union { fields ... })
type UnionDef struct {
	Block Block
}

func (ud *UnionDef) expr() {}

// EnumDef represents the definition of a enum body(enum { fields ... })
type EnumDef struct {
	Block Block
}

func (sd *EnumDef) expr() {}

// PrototypeDef represents the definition of a prototype (proc(int, int) -> int)
type PrototypeDef struct {
	Params     []Field
	ReturnType Expr
}

func (pd *PrototypeDef) expr() {}

// Option represents a single metadata assignation
type Option struct {
	Name  Expr
	Value Expr
}

// OpionBlock represents group of metadata assignations
type OptionBlock struct {
	Options []Option
}

func (ob *OptionBlock) decl() {}

// Block represents a sequence of declarations within a scope ({})
type Block struct {
	Decls []Decl
}

// Field represents a binding declaration (name : Type = value)
type Field struct {
	Name    Expr
	Type    Expr
	Value   Expr
	Options *OptionBlock
}

func (fi *Field) decl() {}

// TypeDecl represents a type declaration ("type Name Type" or "proc Name(arg: Type) -> Type")
type TypeDecl struct {
	Name          Ident
	GenericParams []Field
	Type          Expr
}

func (ty *TypeDecl) decl() {}

// Schema represents the data of an entire schema file
type Schema struct {
	Decls []Decl
}
