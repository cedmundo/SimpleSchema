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
func (p *Prototype) GeneratePrototype() string {
	proto := &strings.Builder{}
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
	return makeIndent(depth) + p.Prototype.GeneratePrototype() + ";"
}

func makeIndent(depth int) string {
	indent := &strings.Builder{}
	for range depth {
		indent.WriteString("  ")
	}

	return indent.String()
}
