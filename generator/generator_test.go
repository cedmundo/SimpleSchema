package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type mockString string

func (m mockString) decl() {}

func (m mockString) attr() {}

func (m mockString) expr() {}

func (m mockString) Generate(d int) string {
	indent := ""
	for range d {
		indent += "  "
	}
	return indent + string(m)
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
				Decls: []Decl{mockString("hello")},
			},
			depth:          0,
			expectedString: "hello\n",
		},
		{
			name: "statement with padding",
			file: &File{
				Decls: []Decl{
					mockString("statement"),
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
			module:         &ModuleWard{Name: "HELLO_H", Decls: []Decl{mockString("hello")}},
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
			list:           AttrList([]Attr{mockString("example")}),
			expectedString: "example ",
		},
		{
			name:           "multiple attributes",
			list:           AttrList([]Attr{mockString("hello"), mockString("world")}),
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
			param:          &Param{Type: mockString("int")},
			expectedString: "int",
		},
		{
			name:           "param with name and type",
			param:          &Param{Name: mockString("x"), Type: mockString("int")},
			expectedString: "int x",
		},
		{
			name: "param with attributes, name and type",
			param: &Param{
				Attrs: []Attr{mockString("_Alignas(16)")},
				Name:  mockString("x"),
				Type:  mockString("int"),
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
		expectedString string
	}{
		{
			name: "prototype with type and name, nil params",
			proto: &Prototype{
				Type:   mockString("int"),
				Name:   mockString("hello"),
				Params: nil,
			},
			expectedString: "int hello()",
		},
		{
			name: "prototype with type and name, empty params",
			proto: &Prototype{
				Type:   mockString("int"),
				Name:   mockString("hello"),
				Params: []Param{},
			},
			expectedString: "int hello()",
		},
		{
			name: "prototype with type and name, single param",
			proto: &Prototype{
				Type: mockString("int"),
				Name: mockString("hello"),
				Params: []Param{
					{
						Type: mockString("int"),
						Name: mockString("x"),
					},
				},
			},
			expectedString: "int hello(int x)",
		},
		{
			name: "prototype with type and name, multiple param",
			proto: &Prototype{
				Type: mockString("int"),
				Name: mockString("hello"),
				Params: []Param{
					{
						Type: mockString("int"),
						Name: mockString("x"),
					},
					{
						Type: mockString("int"),
						Name: mockString("y"),
					},
				},
			},
			expectedString: "int hello(int x, int y)",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.proto.GeneratePrototype()
			require.Equal(t, tt.expectedString, actualString)
		})
	}
}

func TestPrototypeDecl_Generate(t *testing.T) {
	decl := &PrototypeDecl{
		Prototype: Prototype{
			Attrs: []Attr{
				mockString("__attribute__(unused)"),
			},
			Name: mockString("hello"),
			Type: mockString("int"),
			Params: []Param{
				{
					Type: mockString("int"),
					Name: mockString("x"),
				},
			},
		},
	}

	actualString := decl.Generate(0)
	expectedString := "__attribute__(unused) int hello(int x);"
	require.Equal(t, expectedString, actualString)
}
