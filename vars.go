package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Token int

const (
	Variable Token = iota
	Int
	Float
	String
	Bool
	EndOfFile
	// Folder
)

type Var struct {
	Type        Token
	ValueType   Token
	Name        string
	StringValue string
	IntValue    int
	FloatValue  float32
	BoolValue   bool
}

// @Cleanup: Refactor all prints to use normal logs

func InitVariables(file string) {
	lexer := NewLexer(file)
	vars := lexer.Parse()
	for _, v := range vars {
		fmt.Printf("%+v\n", v)
	}
}

type Lexer struct {
	CursorLine         int
	CursorLinePosition int

	Input  []rune
	Cursor int
}

func NewLexer(filePath string) *Lexer {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading %s vars : %s\n", filePath, err)
		return nil
	}

	return &Lexer{
		Input: []rune(string(fileData)),
	}
}

func (l *Lexer) Parse() []Var {
	vars := []Var{}

	for {
		v, err := l.composeNewVar()
		if err != nil && v.Type != EndOfFile {
			panic(err)
		}

		if v.Type == EndOfFile {
			break
		}

		vars = append(vars, v)
	}

	return vars
}

func (l *Lexer) composeNewVar() (Var, error) {
	v := &Var{}
	v.Type = EndOfFile

	for {
		ch, err := l.peekCharater()
		if err != nil {
			break
		}

		for unicode.IsSpace(ch) {
			l.eatCharacter()
			ch, err = l.peekCharater()
			if err != nil {
				return *v, err
			}
		}

		// Comments
		if ch == '#' {
			err = l.eatComment()
			if err != nil {
				return *v, err
			}

			continue
		}

		// Is start of identifier
		if unicode.IsLetter(ch) {
			v, err := l.eatIdentifier()
			if err != nil {
				return *v, err
			}

			err = l.FillVariable(v)
			if err != nil {
				return *v, err
			}

			return *v, nil
		}
	}

	return *v, nil
}

func (l *Lexer) peekCharater() (rune, error) {
	if l.Cursor >= len(l.Input) {
		return 0, fmt.Errorf("end of file at [%d:%d]", l.CursorLine, l.CursorLinePosition)
	}

	return l.Input[l.Cursor], nil
}

func (l *Lexer) eatCharacter() {
	if l.Cursor >= len(l.Input) {
		return
	}

	if l.Input[l.Cursor] == '\n' {
		l.CursorLine++
		l.CursorLinePosition = 0
	}

	l.Cursor++
	l.CursorLinePosition++
}

func (l *Lexer) eatComment() error {
	ch, err := l.peekCharater()
	if err != nil {
		return err
	}

	// @Cleanup: Double check for that?
	if ch == '#' {
		return l.eatLine()
	}

	return fmt.Errorf("u§nexpected character at [%d:%d]", l.CursorLine, l.CursorLinePosition)
}

func (l *Lexer) eatLine() error {
	for {
		ch, err := l.peekCharater()
		if err != nil {
			return err
		}

		l.eatCharacter()

		if ch == '\n' {
			break
		}
	}

	return nil
}

func (l *Lexer) eatIdentifier() (*Var, error) {
	t := &Var{}
	t.Type = Variable

	err := l.parseIdentifier(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (l *Lexer) parseIdentifier(v *Var) error {
	stringBuilder := strings.Builder{}

	ch, err := l.peekCharater()
	if err != nil {
		return err
	}

	// Is start of variable name
	if unicode.IsLetter(ch) {
		for {
			ch, err = l.peekCharater()
			if err != nil {
				// Found something
				if stringBuilder.Len() > 0 {
					break
				}

				return err
			}

			if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
				stringBuilder.WriteRune(ch)
				l.eatCharacter()
				continue
			}

			break
		}
	} else {
		return fmt.Errorf("expected part of variable at [%d:%d]", l.CursorLine, l.CursorLinePosition)
	}

	v.Name = stringBuilder.String()

	return nil
}

func (l *Lexer) FillVariable(v *Var) error {
	stringBuilder := strings.Builder{}

	ch, err := l.peekCharater()
	if err != nil {
		return fmt.Errorf("unexpected end of file while parsing variable [%s] at [%d:%d]", v.Name, l.CursorLine, l.CursorLinePosition)
	}

	// eat all spaces between variable name and value
	if !unicode.IsSpace(ch) {
		return fmt.Errorf("unexpected charater '%c' between variable and value at [%d:%d]", ch, l.CursorLine, l.CursorLinePosition)
	}

	for unicode.IsSpace(ch) {
		l.eatCharacter()
		ch, err = l.peekCharater()
		if err != nil {
			return err
		}
	}

	// Is start of string
	if unicode.IsLetter(ch) {
		for {
			ch, err = l.peekCharater()
			if err != nil {
				// Found something
				if stringBuilder.Len() > 0 {
					break
				}

				return err
			}

			if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
				stringBuilder.WriteRune(ch)
				l.eatCharacter()
				continue
			}

			break
		}

		v.StringValue = stringBuilder.String()
		v.ValueType = String

		if v.StringValue == "true" {
			v.BoolValue = true
			v.ValueType = Bool
		} else if v.StringValue == "false" {
			v.BoolValue = false
			v.ValueType = Bool
		}

		return nil
	}

	// Is start of int
	if unicode.IsDigit(ch) {
		for {
			ch, err = l.peekCharater()
			if err != nil {
				// Found something
				if stringBuilder.Len() > 0 {
					break
				}
			}

			if unicode.IsDigit(ch) || ch == '.' {
				stringBuilder.WriteRune(ch)
				l.eatCharacter()
				continue
			}

			break
		}

		value := stringBuilder.String()
		if strings.Contains(value, ".") {
			fl, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return fmt.Errorf("error parsing float '%s' at [%d:%d]", value, l.CursorLine, l.CursorLinePosition)
			}

			v.FloatValue = float32(fl)
		} else {
			v.IntValue, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing int '%s' at [%d:%d]", value, l.CursorLine, l.CursorLinePosition)
			}
		}
	}

	// Is start of string
	if ch == '"' {
		l.eatCharacter()

		for {
			ch, err = l.peekCharater()
			if err != nil {
				// Found something
				if stringBuilder.Len() > 0 {
					break
				}

				return err
			}

			if ch == '\n' {
				return fmt.Errorf("unexpected newline while parsing string value '%s'at [%d:%d]", stringBuilder.String(), l.CursorLine, l.CursorLinePosition)
			}

			if ch == '"' {
				l.eatCharacter()
				break
			}

			stringBuilder.WriteRune(ch)
			l.eatCharacter()
		}

		v.StringValue = stringBuilder.String()
	}

	if v.StringValue == "" && v.IntValue == 0 && v.FloatValue == 0 {
		return fmt.Errorf("expected '%s' value at [%d:%d]", v.Name, l.CursorLine, l.CursorLinePosition)
	}

	return nil
}