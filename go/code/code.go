package code

import (
	"aakimov/marslang/ast"
	"aakimov/marslang/object"
)

type Code struct {
	id  string
	ast *ast.StatementsBlock
	env *object.Environment
}
