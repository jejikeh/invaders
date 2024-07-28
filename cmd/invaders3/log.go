package main

import (
	"io"
	"log"
	"os"
)

var defalutOut = os.Stdout

var logger = NewLogger(defalutOut)

type Logger struct {
	engine log.Logger
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{
		engine: *log.New(out, "[engine] ", log.Lshortfile),
	}
}
