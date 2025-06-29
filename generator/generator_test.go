package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type mockDecl string

func (m mockDecl) decl() {}

func (m mockDecl) Generate(d int) string {
	return makeIndent(d) + string(m)
}

type mockAttr string

func (m mockAttr) attr() {}

func (m mockAttr) Generate(d int) string {
	return string(m)
}

type mockExpr string

func (m mockExpr) expr() {}

func (m mockExpr) Generate(d int) string {
	return string(m)
}

func TestFile_Generate(t *testing.T) {
	cases := []struct {
		name           string
		file           *File
		depth          int
		expectedString string
	}{
		{
			name:           "empty file",
			file:           &File{},
			depth:          0,
			expectedString: "",
		},
		{
			name: "single statement",
			file: &File{
				Decls: []Decl{mockDecl("hello")},
			},
			depth:          0,
			expectedString: "hello\n",
		},
		{
			name: "statement with padding",
			file: &File{
				Decls: []Decl{
					mockDecl("statement"),
				},
			},
			depth:          1,
			expectedString: "  statement\n",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			generated := tt.file.Generate(tt.depth)
			require.Equal(t, tt.expectedString, generated)
		})
	}
}

func TestModuleWard_Generate(t *testing.T) {
	cases := []struct {
		name           string
		module         *ModuleWard
		depth          int
		expectedString string
	}{
		{
			name:           "empty ward",
			module:         &ModuleWard{Name: "HELLO_H"},
			depth:          0,
			expectedString: "#ifndef HELLO_H\n#define HELLO_H\n#endif /* HELLO_H */\n",
		},
		{
			name:           "single statement ward",
			module:         &ModuleWard{Name: "HELLO_H", Decls: []Decl{mockDecl("hello")}},
			depth:          0,
			expectedString: "#ifndef HELLO_H\n#define HELLO_H\nhello\n#endif /* HELLO_H */\n",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.module.Generate(tt.depth)
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestInclude_Generate(t *testing.T) {
	cases := []struct {
		name           string
		include        *Include
		depth          int
		expectedString string
	}{
		{
			name:           "empty include",
			include:        &Include{},
			depth:          0,
			expectedString: "#include <>",
		},
		{
			name: "non-relative include",
			include: &Include{
				File:     "hello.h",
				Relative: false,
			},
			depth:          0,
			expectedString: "#include <hello.h>",
		},
		{
			name: "relative include",
			include: &Include{
				File:     "hello.h",
				Relative: true,
			},
			depth:          0,
			expectedString: `#include "hello.h"`,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.include.Generate(tt.depth)
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestAttrList_GenerateList(t *testing.T) {
	cases := []struct {
		name           string
		list           AttrList
		expectedString string
	}{
		{
			name:           "nil list",
			list:           nil,
			expectedString: "",
		},
		{
			name:           "empty list",
			list:           AttrList([]Attr{}),
			expectedString: "",
		},
		{
			name:           "single attribute",
			list:           AttrList([]Attr{mockAttr("example")}),
			expectedString: "example ",
		},
		{
			name:           "multiple attributes",
			list:           AttrList([]Attr{mockAttr("hello"), mockAttr("world")}),
			expectedString: "hello world ",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.list.GenerateList()
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestParam_GenerateParam(t *testing.T) {
	cases := []struct {
		name           string
		param          *Param
		expectedString string
	}{
		{
			name:           "param with type",
			param:          &Param{Type: mockExpr("int")},
			expectedString: "int",
		},
		{
			name:           "param with name and type",
			param:          &Param{Name: mockExpr("x"), Type: mockExpr("int")},
			expectedString: "int x",
		},
		{
			name: "param with attributes, name and type",
			param: &Param{
				Attrs: []Attr{mockAttr("_Alignas(16)")},
				Name:  mockExpr("x"),
				Type:  mockExpr("int"),
			},
			expectedString: "_Alignas(16) int x",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.param.GenerateParam()
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestPrototype_GeneratePrototype(t *testing.T) {
	cases := []struct {
		name           string
		proto          *Prototype
		depth          int
		expectedString string
	}{
		{
			name: "prototype with type and name, nil params",
			proto: &Prototype{
				Type:   mockExpr("int"),
				Name:   mockExpr("hello"),
				Params: nil,
			},
			depth:          0,
			expectedString: "int hello()",
		},
		{
			name: "prototype with type and name, empty params",
			proto: &Prototype{
				Type:   mockExpr("int"),
				Name:   mockExpr("hello"),
				Params: []Param{},
			},
			depth:          0,
			expectedString: "int hello()",
		},
		{
			name: "prototype with type and name, single param",
			proto: &Prototype{
				Type: mockExpr("int"),
				Name: mockExpr("hello"),
				Params: []Param{{
					Type: mockExpr("int"),
					Name: mockExpr("x"),
				}},
			},
			depth:          0,
			expectedString: "int hello(int x)",
		},
		{
			name: "prototype with type and name, multiple param",
			proto: &Prototype{
				Type: mockExpr("int"),
				Name: mockExpr("hello"),
				Params: []Param{
					{
						Type: mockExpr("int"),
						Name: mockExpr("x"),
					},
					{
						Type: mockExpr("int"),
						Name: mockExpr("y"),
					},
				},
			},
			depth:          0,
			expectedString: "int hello(int x, int y)",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.proto.GeneratePrototype(tt.depth)
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestPrototypeDecl_Generate(t *testing.T) {
	decl := &PrototypeDecl{
		Prototype: Prototype{
			Attrs: []Attr{
				mockAttr("__attribute__(unused)"),
			},
			Name: mockExpr("hello"),
			Type: mockExpr("int"),
			Params: []Param{
				{
					Type: mockExpr("int"),
					Name: mockExpr("x"),
				},
			},
		},
	}

	actualString := decl.Generate(0)
	expectedString := "__attribute__(unused) int hello(int x);"
	require.Equal(t, expectedString, actualString)
}

func TestField_Generate(t *testing.T) {
	cases := []struct {
		name           string
		field          *Field
		depth          int
		expectedString string
	}{
		{
			name: "simple field",
			field: &Field{
				Type: mockExpr("int"),
				Name: mockExpr("x"),
			},
			depth:          0,
			expectedString: "int x",
		},
		{
			name: "field with attributes",
			field: &Field{
				Attrs: []Attr{mockAttr("__attr__")},
				Type:  mockExpr("int"),
				Name:  mockExpr("x"),
			},
			depth:          0,
			expectedString: "__attr__ int x",
		},
		{
			name: "field with attributes and depth",
			field: &Field{
				Attrs: []Attr{mockAttr("__attr__")},
				Type:  mockExpr("int"),
				Name:  mockExpr("x"),
			},
			depth:          1,
			expectedString: "  __attr__ int x",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.field.GenerateField(tt.depth)
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestFieldBlock_GenerateBlock(t *testing.T) {
	cases := []struct {
		name           string
		block          FieldBlock
		depth          int
		expectedString string
	}{
		{
			name:           "nil block",
			block:          nil,
			depth:          0,
			expectedString: "{}",
		},
		{
			name:           "empty block",
			block:          FieldBlock([]Field{}),
			depth:          0,
			expectedString: "{}",
		},
		{
			name: "block with one field",
			block: FieldBlock([]Field{
				{
					Type: mockExpr("int"),
					Name: mockExpr("x"),
				},
			}),
			depth:          0,
			expectedString: "{\n  int x;\n}",
		},
		{
			name: "block with multiple fields",
			block: FieldBlock([]Field{
				{
					Type: mockExpr("int"),
					Name: mockExpr("x"),
				},
				{
					Type: mockExpr("int"),
					Name: mockExpr("y"),
				},
			}),
			depth:          0,
			expectedString: "{\n  int x;\n  int y;\n}",
		},
		{
			name: "block with multiple fields and ident",
			block: FieldBlock([]Field{
				{
					Type: mockExpr("int"),
					Name: mockExpr("x"),
				},
				{
					Type: mockExpr("int"),
					Name: mockExpr("y"),
				},
			}),
			depth:          1,
			expectedString: "{\n    int x;\n    int y;\n  }",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.block.GenerateBlock(tt.depth)
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestStruct_Generate(t *testing.T) {
	cases := []struct {
		name           string
		decl           *Struct
		depth          int
		expectedString string
	}{
		{
			name:           "empty struct",
			decl:           &Struct{},
			depth:          0,
			expectedString: "struct {}",
		},
		{
			name: "struct with name but no fields",
			decl: &Struct{
				Name: mockExpr("s"),
			},
			depth:          0,
			expectedString: "struct s {}",
		},
		{
			name: "struct with name with single field",
			decl: &Struct{
				Name: mockExpr("s"),
				Fields: []Field{
					{
						Type: mockExpr("int"),
						Name: mockExpr("x"),
					},
				},
			},
			depth:          0,
			expectedString: "struct s {\n  int x;\n}",
		},
		{
			name: "struct with name with multiple fields",
			decl: &Struct{
				Name: mockExpr("s"),
				Fields: []Field{
					{
						Type: mockExpr("int"),
						Name: mockExpr("x"),
					},
					{
						Type: mockExpr("int"),
						Name: mockExpr("y"),
					},
				},
			},
			depth:          0,
			expectedString: "struct s {\n  int x;\n  int y;\n}",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.decl.Generate(tt.depth)
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestStructDecl_Generate(t *testing.T) {
	cases := []struct {
		name           string
		decl           *StructDecl
		depth          int
		expectedString string
	}{
		{
			name:           "empty struct",
			decl:           &StructDecl{},
			depth:          0,
			expectedString: "struct {};",
		},
		{
			name: "struct with name but no fields",
			decl: &StructDecl{Struct{
				Name: mockExpr("s"),
			}},
			depth:          0,
			expectedString: "struct s {};",
		},
		{
			name: "struct with name with single field",
			decl: &StructDecl{Struct{
				Name: mockExpr("s"),
				Fields: []Field{
					{
						Type: mockExpr("int"),
						Name: mockExpr("x"),
					},
				},
			}},
			depth:          0,
			expectedString: "struct s {\n  int x;\n};",
		},
		{
			name: "struct with name with multiple fields",
			decl: &StructDecl{Struct{
				Name: mockExpr("s"),
				Fields: []Field{
					{
						Type: mockExpr("int"),
						Name: mockExpr("x"),
					},
					{
						Type: mockExpr("int"),
						Name: mockExpr("y"),
					},
				},
			}},
			depth:          0,
			expectedString: "struct s {\n  int x;\n  int y;\n};",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.decl.Generate(tt.depth)
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}
