package logic

import (
	"go/ast"
)

// Implements sort.Interface for []*ast.ImportSpec
type byImportPath []*ast.ImportSpec

func (a byImportPath) Len() int {
	return len(a)
}

func (a byImportPath) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byImportPath) Less(i, j int) bool {
	return a[i].Path.Value < a[j].Path.Value
}
