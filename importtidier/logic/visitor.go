package logic

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"regexp"
	"sort"
	"strings"
)

// Implements ast.Visitor
type myVisitor struct {
	filename  string
	fileSet   *token.FileSet
	source    []byte
	output    []byte
	err       error
	startPos  int
	endPos    int
	locakPkgs []string
}

// Visit implements ast.Visitor
func (v *myVisitor) Visit(node ast.Node) ast.Visitor {
	switch statement := node.(type) {
	case *ast.GenDecl:
		v.detectImportDecl(statement)
	default:
		// intentionally do nothing
	}

	return v
}

func (v *myVisitor) updateFile(file *ast.File) {
	updatedImports := ""

	if len(file.Imports) > 0 {
		v.orderImports(file)
		updatedImports = v.generateImportsFragment(file)
	}
	v.output = v.replaceImports(updatedImports)

	err := v.validate(v.output)
	if err != nil {
		v.err = fmt.Errorf("generated code was invalid, err: %s", err)
		return
	}
}

func (v *myVisitor) orderImports(file *ast.File) {
	sort.Sort(byImportPath(file.Imports))
}

func (v *myVisitor) generateImportsFragment(file *ast.File) string {
	stdLibFragment := ""
	thirdPartyLibFragment := ""
	customFragment := ""

	stdLibRegex := regexp.MustCompile(`(")[a-zA-Z0-9/]+(")`)

	for _, thisImport := range file.Imports {
		// first check local packages
		if contains(v.locakPkgs, thisImport.Path.Value) {
			customFragment += "\t" + v.buildImportLine(thisImport)
		} else if stdLibRegex.MatchString(thisImport.Path.Value) {
			stdLibFragment += "\t" + v.buildImportLine(thisImport)
		} else {
			thirdPartyLibFragment += "\t" + v.buildImportLine(thisImport)
		}
	}

	stPadding, tcPadding := "", ""
	if len(stdLibFragment) > 0 && len(thirdPartyLibFragment) > 0 {
		stPadding = string(lineBreak)
	}
	if len(thirdPartyLibFragment) > 0 && len(customFragment) > 0 {
		tcPadding = string(lineBreak)
	}
	if len(stdLibFragment) > 0 && len(thirdPartyLibFragment) == 0 && len(customFragment) > 0 {
		stPadding = string(lineBreak)
	}

	output := "import (" + string(lineBreak)
	output += stdLibFragment + stPadding + thirdPartyLibFragment + tcPadding + customFragment
	output += ")" + string(lineBreak)

	return output
}

func (v *myVisitor) buildImportLine(thisImport *ast.ImportSpec) string {
	output := ""

	topComment := strings.TrimSpace(thisImport.Doc.Text())
	if len(topComment) > 0 {
		output += "// " + topComment + string(lineBreak) + "\t"
	}

	if thisImport.Name != nil {
		name := strings.TrimSpace(thisImport.Name.Name)

		// remove redundant names
		if !v.nameIsRedundant(name, thisImport.Path.Value) {
			if len(name) > 0 {
				output += name + " "
			}
		}
	}

	output += thisImport.Path.Value

	commentAfter := strings.TrimSpace(thisImport.Comment.Text())
	if len(commentAfter) > 0 {
		output += " // " + commentAfter
	}

	return output + string(lineBreak)
}

func (v *myVisitor) nameIsRedundant(name string, path string) bool {
	// case: `import io "io"`
	// compare name and path after trimming the quotes
	if name == path[1:len(path)-1] {
		return true
	}

	// case: `import proto "github.com/golang/protobuf/proto"`
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash > -1 {
		// trim the slash and the quotes
		pkgDir := path[lastSlash+1 : len(path)-1]
		if name == pkgDir {
			return true
		}
	}

	return false
}

func (v *myVisitor) replaceImports(newImports string) []byte {
	var output []byte

	// replace the imports section
	output = append(output, v.source[:v.startPos]...)
	output = append(output, newImports...)
	output = append(output, v.source[v.endPos:]...)

	return output
}

// Validate the result by running it through GoFmt
func (v *myVisitor) validate(newCode []byte) error {
	_, err := format.Source(newCode)
	return err
}

func (v *myVisitor) detectImportDecl(decl *ast.GenDecl) {
	if decl.Tok != token.IMPORT {
		return
	}

	thisStartPos, thisEndPos := getLineBoundary(v.source, decl.Pos())
	if thisStartPos < v.startPos {
		v.startPos = thisStartPos
	}

	if decl.Rparen.IsValid() {
		// override with `)` if exists
		// NOTE: add 1 for the line break
		thisEndPos = int(decl.Rparen) + 1
	}

	if thisEndPos > v.endPos {
		v.endPos = thisEndPos
	}
}
