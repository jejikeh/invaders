package varfile

import (
	"strings"
)

type lexer struct {
	filePath string
	lines    []string
	line     int
	column   int
}

func newLexer(filePath string, content string) *lexer {
	return &lexer{
		filePath: filePath,
		lines:    strings.Split(content, "\n"),
		line:     0,
		column:   0,
	}
}

type VarType int

const (
	Group VarType = iota
	Number
	String
	Bool
)

type Var struct {
	Type  VarType
	Value any
	Name  string
}

func (l *lexer) ParseFile() ([]Var, error) {
	vars := make([]Var, 0)

	for lineIdx, line := range l.lines {
		l.line = lineIdx

	fieldLex:
		for columnIdx, field := range strings.Fields(line) {
			l.column = columnIdx

			switch {
			case strings.HasPrefix(field, "[") && strings.HasSuffix(field, "]"):
				groupName := strings.TrimSuffix(strings.TrimPrefix(field, "["), "]")
				vars = append(vars, Var{
					Type:  Group,
					Name:  groupName,
					Value: groupName,
				})
			case strings.HasPrefix(field, "#"):
				break fieldLex
			default:
			}
		}
	}

	return vars, nil
}
