package logic

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
)

const lineBreak = '\n'

// ProcessFiles will process the passed file and do import tidy and output the file
func ProcessFile(filename string, localPkgs []string) {
	originalCode, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("skipping file %s: failed to read with err: %v\n", filename, err)
		return
	}

	newCode, err := processFile(filename, originalCode, localPkgs)
	if err != nil {
		fmt.Printf("skipping file %s: failed to generate with err: %v\n", filename, err)
		return
	}

	updateSourceFile(filename, originalCode, newCode)
}

func updateSourceFile(filename string, originalCode, newCode []byte) {
	if !bytes.Equal(originalCode, newCode) {
		err := os.WriteFile(filename, newCode, 0)
		if err != nil {
			fmt.Printf("skipping file %s: failed to write with err: %v\n", filename, err)
			return
		}
	} else {
		fmt.Println("no need to update, there is no change")
	}
}

func processFile(filename string, source []byte, localPkgs []string) ([]byte, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// special case: no imports
	if len(file.Imports) == 0 {
		fmt.Println("no imports, skip")
		return source, nil
	}

	visitor := newVisitor(filename, fileSet, source, localPkgs)
	ast.Walk(visitor, file)

	if visitor.err != nil {
		return nil, visitor.err
	}

	visitor.updateFile(file)

	return visitor.output, nil
}

func newVisitor(filename string, fileSet *token.FileSet, source []byte, locakPkgs []string) *myVisitor {
	return &myVisitor{
		filename:  filename,
		fileSet:   fileSet,
		source:    source,
		startPos:  math.MaxInt32,
		endPos:    -1,
		locakPkgs: locakPkgs,
	}
}
