package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Log interface{
	New()
	WriteInfo(string)
	WriteWarn(string)
	WriteFatal(string)
}

type log struct{
	l *logrus.Logger
}

func (l *log)New(){
	l.l = logrus.New()
	l.l.SetOutput(os.Stdout)
	l.l.SetLevel(logrus.InfoLevel)
}

func (l *log)WriteInfo(t string){
	l.l.Infoln(t)
}

func (l *log)WriteWarn(t string){
	l.l.Warnln(t)
}

func (l *log)WriteFatal(t string){
	l.l.Warnln(t)
}

func New()Log{
	l := new(log)
	l.New()
	return l
}