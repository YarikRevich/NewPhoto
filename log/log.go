package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger = New()
)

func OpenLogFile() *os.File {
	if _, err := os.Stat("error.log"); os.IsNotExist(err) {
		if _, err := os.Create("error.log"); err != nil{
			fmt.Fprintln(os.Stderr, err)
		}
	}
	f, err := os.OpenFile("error.log", os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return f
}

func New() *logrus.Logger {
	l := logrus.New()
	l.Out = OpenLogFile()
	return l
}
