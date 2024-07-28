package log

import (
	"io"
	"log"
	"os"
)

var LogWriter = os.Stdout

var Engine = New(LogWriter, "[engine] ")

type Logger struct {
	log.Logger
}

func New(out io.Writer, name string) *Logger {
	return &Logger{
		Logger: *log.New(out, name, log.Llongfile),
	}
}
