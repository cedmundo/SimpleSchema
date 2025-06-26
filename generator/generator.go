package generator

import (
	"fmt"
	"strings"
)

type CGeneratorContext struct {
	Indentation uint
}

type CGenerator interface {
	Generate(gc CGeneratorContext) string
}

type CWord struct {
	Value string
}

func (ci *CWord) Generate(gc CGeneratorContext) string {
	return ci.Value
}

type CUnary struct {
	Prefix  string
	Suffix  string
	Operand CGenerator
}

func (cu *CUnary) Generate(gc CGeneratorContext) string {
	return strings.Join([]string{cu.Prefix, cu.Operand.Generate(gc), cu.Suffix}, " ")
}

type CBinary struct {
	Infix string
	Left  CGenerator
	Right CGenerator
}

func (cb *CBinary) Generate(gc CGeneratorContext) string {
	left := cb.Left.Generate(gc)
	right := cb.Right.Generate(gc)
	return strings.Join([]string{left, cb.Infix, right}, " ")
}

type CGroup struct {
	Wrap CGenerator
}

func (cg *CGroup) Generate(gc CGeneratorContext) string {
	return fmt.Sprintf("(%s)", cg.Generate(gc))
}

type CIndex struct {
	Base  CGenerator
	Index CGenerator
}

func (ci *CIndex) Generate(gc CGeneratorContext) string {
	return fmt.Sprintf("%s[%s]", ci.Base.Generate(gc), ci.Index.Generate(gc))
}
