package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type DisplayVars struct {
	VSync bool
}

// @Incomplete: For now, only one global Inits Variable is handled.

type Inits map[string]any

var Variables Inits = map[string]any{}

// @Incomplete: Add Map type support

func (Inits) Add(v any) {
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Ptr {
		panicString := fmt.Sprintf("struct must be a pointer, got \n\t[type]: '%s'\n\t[name]: '%s'", t.Kind(), t.Name())
		panic(panicString)
	}

	tName := t.Elem().Name()
	Variables[tName] = v
}

func (Inits) AddNamed(v any, name string) {
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Ptr {
		panicString := fmt.Sprintf("struct must be a pointer, got \n\t[type]: '%s'\n\t[name]: '%s'", t.Kind(), t.Name())
		panic(panicString)
	}

	Variables[name] = v
}

func (Inits) Find(name string) (any, error) {
	if v, ok := Variables[name]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("struct '%s' not found", name)
}

func InitVariables(file string) error {
	//
	// Add structs.
	//

	Variables.AddNamed(&GameVolume, "Volume")

	//
	// Parse .variables file and map content to structs.
	//

	var currentFolder *interface{}
	var currentFolderName string
	lexer, err := NewLexer(file)
	if err != nil {
		return err
	}

	vars := lexer.Parse()
	for _, v := range vars {
		// fmt.Printf("%s\n", v.String())
		if v.Type == Folder {
			findStruct, err := Variables.Find(v.Name)
			if err != nil {
				return err
			}

			currentFolder = &findStruct
			currentFolderName = v.Name
		}

		if v.Type == Variable && currentFolder == nil {
			return fmt.Errorf("the variable '%s' is not in any folder", v.Name)
		}

		if v.Type == Variable {
			// @Cleanup: Maybe store Var value in any, but this way gives little bit of a type safety
			field := reflect.ValueOf(*currentFolder).Elem().FieldByName(v.Name)

			if field.Kind() != v.ValueType {
				return fmt.Errorf("the variable '%s.%s' is not of type '%s'", currentFolderName, v.Name, field.Kind())
			}

			// @Cleanup: If stick to separate fields - types, when maybe create a handy function what does a map from Kind to Set Function?
			switch v.ValueType {
			case reflect.Int:
				field.SetInt(int64(v.IntValue))
			case reflect.Float32:
				field.SetFloat(float64(v.FloatValue))
			case reflect.Bool:
				field.SetBool(v.BoolValue)
			case reflect.String:
				field.SetString(v.StringValue)
			}
		}
	}

	return nil
}

type Token int

const (
	Variable Token = iota
	Int
	Float
	String
	Bool
	EndOfFile
	Folder
)

// @Incomplete: Handle folders ans subfolders in 'parser' maybe?

func (t Token) String() string {
	return [...]string{
		"Variable",
		"Int",
		"Float",
		"String",
		"Bool",
		"EndOfFile",
		"Folder",
	}[t]
}

type Var struct {
	Type        Token
	ValueType   reflect.Kind
	Name        string
	StringValue string
	IntValue    int
	FloatValue  float32
	BoolValue   bool
}

func (v *Var) String() string {
	headerString := fmt.Sprintf("[%s] %s", v.Type.String(), v.Name)
	bodyString := ""

	switch v.ValueType {
	case reflect.Int:
		bodyString = fmt.Sprintf("\t[Int]: %d", v.IntValue)
	case reflect.Float32:
		bodyString = fmt.Sprintf("\t[Float]: %f", v.FloatValue)
	case reflect.String:
		bodyString = fmt.Sprintf("\t[String]: %s", v.StringValue)
	case reflect.Bool:
		bodyString = fmt.Sprintf("\t[Bool]: %t", v.BoolValue)
	}

	return fmt.Sprintf("%s\n%s\n-----\n", headerString, bodyString)
}

type Lexer struct {
	CursorLine         int
	CursorLinePosition int

	Input  []rune
	Cursor int
}

func NewLexer(filePath string) (*Lexer, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading %s vars : %s\n", filePath, err)
	}

	return &Lexer{
		Input: []rune(string(fileData)),
	}, nil
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

		if ch == '/' {
			l.eatCharacter()
			ch, err := l.peekCharater()
			if err != nil {
				return Var{}, err
			}

			if ch == ':' {
				// @Incomplete: Handle folders :)
				l.eatCharacter()

				ch, err = l.peekCharater()
				if err != nil {
					return Var{}, err
				}

				// Eat spaces
				for unicode.IsSpace(ch) {
					l.eatCharacter()
					ch, err = l.peekCharater()
					if err != nil {
						return *v, err
					}
				}

				v, err := l.eatIdentifier()
				if err != nil {
					// @Cleanup: Return pointer instead of copy
					return Var{}, err
				}

				v.Type = Folder
				return *v, nil
			}
		}

		// @Hack: sinse we are no displaying erros if type is EndOfFile
		v.Type = Variable
		return *v, fmt.Errorf("unexpected character '%c' at [%d:%d]", ch, l.CursorLine, l.CursorLinePosition)
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
		return fmt.Errorf("expected part of variable at [%d:%d], but got '%c'", l.CursorLine, l.CursorLinePosition, ch)
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
		v.ValueType = reflect.String

		if v.StringValue == "true" {
			v.BoolValue = true
			v.ValueType = reflect.Bool
		} else if v.StringValue == "false" {
			v.BoolValue = false
			v.ValueType = reflect.Bool
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
			v.ValueType = reflect.Float32

			return nil
		} else {
			v.IntValue, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing int '%s' at [%d:%d]", value, l.CursorLine, l.CursorLinePosition)
			}

			v.ValueType = reflect.Int
			return nil
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
		v.ValueType = reflect.String

		return nil
	}

	if v.StringValue == "" && v.IntValue == 0 && v.FloatValue == 0 {
		return fmt.Errorf("expected '%s' value at [%d:%d]", v.Name, l.CursorLine, l.CursorLinePosition)
	}

	return nil
}
