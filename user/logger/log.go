package logger

import (
	"log"
	"os"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Logger(f *os.File) (*log.Logger, *log.Logger, *log.Logger, *log.Logger) {
	Info = log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	Warning = log.New(f, "WARN: ", log.Ldate|log.Ltime|log.Llongfile)
	Error = log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
	Debug = log.New(f, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile)
	Info.SetOutput(f)
	Warning.SetOutput(f)
	Error.SetOutput(f)
	Debug.SetOutput(f)
	return Debug, Info, Error, Warning
}
