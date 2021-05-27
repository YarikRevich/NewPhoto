package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger = New()
)

type CLogB interface {
	OpenErrorLogFile()
	OpenAccessLogFile()

	UsingErrorLogFile() *Log
	UsingAccessLogFile() *Log
}

type CLogC interface {
	GetFieldError(query string) *logrus.Entry
	GetFieldAccess(command string) *logrus.Entry

	CFatalln(query string, err interface{})
	CInfoln(command string)
}

type CLog interface {

	//Basic CLog method to run it
	CLogB

	//Custom commands implemented in log
	CLogC
}

type Log struct {
	*logrus.Logger
	accessLogFile *os.File
	errorLogFile  *os.File
}

func (l *Log) OpenErrorLogFile() {
	if _, err := os.Stat("log/error.log"); os.IsNotExist(err) {
		if _, err := os.Create("log/error.log"); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	f, err := os.OpenFile("log/error.log", os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	l.errorLogFile = f
}

func (l *Log) OpenAccessLogFile() {
	if _, err := os.Stat("log/access.log"); os.IsNotExist(err) {
		if _, err := os.Create("log/access.log"); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	f, err := os.OpenFile("log/access.log", os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	l.accessLogFile = f
}

func (l *Log) UsingAccessLogFile() *Log {
	l.SetOutput(l.accessLogFile)
	return l
}

func (l *Log) UsingErrorLogFile() *Log {
	l.SetOutput(l.errorLogFile)
	return l
}

func (l *Log) GetFieldError(query string) *logrus.Entry {
	return l.WithFields(logrus.Fields{"query": query})
}

func (l *Log) GetFieldAccess(command string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"command": command,
	})
}

//Custom fatalln method which is used with fields
func (l *Log) CFatalln(query string, err interface{}) {
	l.GetFieldError(query).Fatalln(err)
}

//Custom CInfoln method which is used with fields
func (l *Log) CInfoln(command string) {
	l.GetFieldAccess(command).Infoln()
}

func New() CLog {
	l := &Log{}
	l.Logger = logrus.New()
	l.OpenAccessLogFile()
	l.OpenErrorLogFile()
	return l
}
