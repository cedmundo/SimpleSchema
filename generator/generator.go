package generator

import (
	"fmt"
	"strings"
)

// Generator transforms internal data to a plain string containg source code
type Generator interface {
	Generate(depth int) string
}

// Decl represents any declaration, including directives and macros
type Decl interface {
	Generator
	decl()
}

// Attr represents any prefix attribute
type Attr interface {
	Generator
	attr()
}

// Expr represents both type and data expressions
type Expr interface {
	Generator
	expr()
}

// File contains declarations
type File struct {
	Decls []Decl
}

// Generate get the code for each declaration and appends a new line
func (f *File) Generate(depth int) string {
	contents := &strings.Builder{}
	for _, decl := range f.Decls {
		line := decl.Generate(depth)
		contents.WriteString(line)
		contents.WriteRune('\n')
	}
	return contents.String()
}

// ModuleWard represents a ifdef,define,endif macro ward
type ModuleWard struct {
	Name  string
	Decls []Decl
}

func (m *ModuleWard) decl() {}

// Generate wraps the following declarations within the ifndef,endif
func (m *ModuleWard) Generate(depth int) string {
	contents := &strings.Builder{}
	contents.WriteString("#ifndef ")
	contents.WriteString(m.Name)
	contents.WriteString("\n")

	contents.WriteString("#define ")
	contents.WriteString(m.Name)
	contents.WriteString("\n")

	for _, decl := range m.Decls {
		line := decl.Generate(depth)
		contents.WriteString(line)
		contents.WriteRune('\n')
	}

	contents.WriteString("#endif /* ")
	contents.WriteString(m.Name)
	contents.WriteString(" */\n")
	return contents.String()
}

// Include represents an include directive
type Include struct {
	File     string
	Relative bool
}

func (i *Include) decl() {}

// Generate outputs include directive with double quotes or between <> if relative or not
func (i *Include) Generate(depth int) string {
	if i.Relative {
		return fmt.Sprintf(`#include "%s"`, i.File)
	}

	return fmt.Sprintf(`#include <%s>`, i.File)
}

// AttrList is a list containing individual attributes
type AttrList []Attr

// GenerateList makes an attribute list if available, otherwise returns an empty string
func (al AttrList) GenerateList() string {
	if len(al) == 0 {
		return ""
	}

	attrs := &strings.Builder{}
	for _, attr := range al {
		attrs.WriteString(attr.Generate(0))
		attrs.WriteString(" ")
	}

	return attrs.String()
}

// Param represents a param with name and type and optionally attributes
type Param struct {
	Attrs []Attr
	Name  Expr
	Type  Expr
}

// GenerateParam outputs the code for a single parameter
func (p *Param) GenerateParam() string {
	param := &strings.Builder{}
	param.WriteString(AttrList(p.Attrs).GenerateList())

	param.WriteString(p.Type.Generate(0))
	if p.Name != nil {
		param.WriteRune(' ')
		param.WriteString(p.Name.Generate(0))
	}

	return param.String()
}

// Prototype represents a prototype data (only type-name-args declaration)
type Prototype struct {
	Attrs  []Attr
	Type   Expr
	Name   Expr
	Params []Param
}

// GeneratePrototype outputs the code for the prototype only (without function body or trailing semicolon)
func (p *Prototype) GeneratePrototype(depth int) string {
	proto := &strings.Builder{}
	proto.WriteString(makeIndent(depth))
	proto.WriteString(AttrList(p.Attrs).GenerateList())
	proto.WriteString(p.Type.Generate(0))
	proto.WriteRune(' ')
	proto.WriteString(p.Name.Generate(0))
	proto.WriteRune('(')

	for i, param := range p.Params {
		if i != 0 {
			proto.WriteString(", ")
		}
		proto.WriteString(param.GenerateParam())
	}

	proto.WriteRune(')')
	return proto.String()
}

// PrototypeDecl represents an actual prototype in the file
type PrototypeDecl struct {
	Prototype Prototype
}

func (p *PrototypeDecl) decl() {}

// Generate outputs the prototype in the block
func (p *PrototypeDecl) Generate(depth int) string {
	return p.Prototype.GeneratePrototype(depth) + ";"
}

// Field represents a field within a struct or union
type Field struct {
	Attrs []Attr
	Type  Expr
	Name  Expr
}

// Generate outputs the actual field with indentation
func (f *Field) GenerateField(depth int) string {
	field := &strings.Builder{}
	field.WriteString(makeIndent(depth))
	field.WriteString(AttrList(f.Attrs).GenerateList())
	field.WriteString(f.Type.Generate(depth))
	field.WriteRune(' ')
	field.WriteString(f.Name.Generate(depth))
	return field.String()
}

// FieldBlock is a list of fields
type FieldBlock []Field

// GenerateBlock returns the block wrapped on "{}" containing all fields
func (fb FieldBlock) GenerateBlock(depth int) string {
	block := &strings.Builder{}
	block.WriteRune('{')

	if len(fb) > 0 {
		block.WriteRune('\n')
	}

	for _, field := range fb {
		block.WriteString(field.GenerateField(depth + 1))
		block.WriteString(";\n")
	}

	block.WriteString(makeIndent(depth))
	block.WriteRune('}')
	return block.String()
}

// Struct is an expression that can be used as type
type Struct struct {
	Attrs  []Attr
	Name   Expr
	Fields []Field
}

func (s *Struct) expr() {}

// Generate returns the equivalent code for a structure with fields
func (s *Struct) Generate(depth int) string {
	strct := &strings.Builder{}
	strct.WriteString(makeIndent(depth))
	strct.WriteString(AttrList(s.Attrs).GenerateList())
	strct.WriteString("struct ")
	if s.Name != nil {
		strct.WriteString(s.Name.Generate(depth))
		strct.WriteRune(' ')
	}
	strct.WriteString(FieldBlock(s.Fields).GenerateBlock(depth))
	return strct.String()
}

// StructDecl represents a struct declaration
type StructDecl struct {
	Struct Struct
}

func (sd *StructDecl) decl() {}

// Generates the struct expr with a trailing semicolon
func (sd *StructDecl) Generate(depth int) string {
	return sd.Struct.Generate(depth) + ";"
}

func makeIndent(depth int) string {
	indent := &strings.Builder{}
	for range depth {
		indent.WriteString("  ")
	}

	return indent.String()
}
