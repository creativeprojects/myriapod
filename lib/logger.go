package lib

import (
	"fmt"
	stdlog "log"
	"os"
)

type Logger interface {
	Print(v ...interface{})
}

var (
	log Logger
)

type DefaultLogger struct {
	log *stdlog.Logger
}

func (d DefaultLogger) Print(v ...interface{}) {
	d.log.Output(3, fmt.Sprint(v...))
}

func init() {
	log = DefaultLogger{
		log: stdlog.New(os.Stdout, "", stdlog.LstdFlags|stdlog.Lshortfile),
	}
}

func SetLogger(logger Logger) {
	log = logger
}
