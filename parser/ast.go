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

// TypeExpr represent any expresion that is compatible with types, such ids and arrays
type TypeExpr interface {
	typeExpr()
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

func (i *Ident) typeExpr() {}

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

// TypeIndex represents a selection expression in types domain (Type[index])
type TypeIndex struct {
	Base  TypeExpr
	Index Expr
}

func (ti *TypeIndex) typeExpr() {}

// StructDef represents the definition of a struct body(struct { fields ... })
type StructDef struct {
	Block Block
}

func (sd *StructDef) typeExpr() {}

// UnionDef represents the definition of a union body(union { fields ... })
type UnionDef struct {
	Block Block
}

func (ud *UnionDef) typeExpr() {}

// EnumDef represents the definition of a enum body(enum { fields ... })
type EnumDef struct {
	Block Block
}

func (sd *EnumDef) typeExpr() {}

// PrototypeDef represents the definition of a prototype (proc(int, int) -> int)
type PrototypeDef struct {
	Params     []Field
	ReturnType TypeExpr
}

func (pd *PrototypeDef) typeExpr() {}

// Block represents a sequence of declarations within a scope ({})
type Block struct {
	Decls []Decl
}

// Field represents a binding declaration (name : Type = value)
type Field struct {
	Name  Ident
	Type  TypeExpr
	Value Expr
	Block *Block
}

func (fi *Field) decl() {}

// Option represents metadata attached to the current scope (option name = "x")
type Option struct {
	Name      Ident
	Overwrite bool
	Block     *Block
}

func (op *Option) decl() {}

// TypeDecl represents a type declaration ("type Name Type" or "proc Name(arg: Type) -> Type")
type TypeDecl struct {
	Name   Ident
	Params []Field
	Type   TypeExpr
}

func (ty *TypeDecl) decl() {}

// Schema represents the data of an entire schema file
type Schema struct {
	Decls []Decl
}
