package logic

import (
	"go/token"
	"strings"
)

// It checks whether the import path contains any of the prefix,
// if yes, consider it as a local package
func contains(prefixes []string, importPath string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.Trim(importPath, "\""), prefix) { // becasue the value wil start with
			return true
		}
	}
	return false
}

// Run to end of line in both directions if not at line start/end.
func getLineBoundary(source []byte, pos token.Pos) (int, int) {
	charAtPos := source[pos]
	if charAtPos == lineBreak {
		return int(pos), int(pos)
	}

	startPos, endPos := int(pos), int(pos)+1
	for startPos > 0 && source[startPos-1] != lineBreak {
		startPos--
	}

	for endPos < len(source) && source[endPos-1] != lineBreak {
		endPos++
	}

	return startPos, endPos
}
