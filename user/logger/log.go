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

const (
	//DEBUG   string = "Debug"
	INFO string = "Info"
	//WARNING        = "Warning"
	ERROR = "Error"
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

func PrintLog(logType string, err error) {
	logs := Error
	if logType == "Debug" {
		logs = Debug
	} else if logType == "Info" {
		logs = Info
	} else if logType == "Warning" {
		logs = Warning
	}
	logs.Println(err.Error())
}
