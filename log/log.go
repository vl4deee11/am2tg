package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	Trace = iota
	Debug
	Info
	Warn
	Error
)

var Logger *lvlWrap
var lvl2str = map[int]string{}
var str2lvl = map[string]int{
	"debug": Debug,
	"info":  Info,
	"warn":  Warn,
	"error": Error,
	"trace": Trace,
}

type lvlWrap struct {
	*log.Logger
	l int
}

func (l *lvlWrap) canPrint(lvl int) bool {
	return lvl >= l.l
}

func (l *lvlWrap) lvlFormat(lvl int, format string, v ...interface{}) string {
	format = fmt.Sprintf("[%s] %s", lvl2str[lvl], format)
	return fmt.Sprintf(format, v...)
}

func (l *lvlWrap) Tracef(format string, v ...interface{}) {
	l.printf(Trace, format, v...)
}

func (l *lvlWrap) Trace(v ...interface{}) {
	l.println(Trace, v...)
}

func (l *lvlWrap) Debugf(format string, v ...interface{}) {
	l.printf(Debug, format, v...)
}

func (l *lvlWrap) Debug(v ...interface{}) {
	l.println(Debug, v...)
}

func (l *lvlWrap) Infof(format string, v ...interface{}) {
	l.printf(Info, format, v...)
}

func (l *lvlWrap) Info(v ...interface{}) {
	l.println(Info, v...)
}

func (l *lvlWrap) Warnf(format string, v ...interface{}) {
	l.printf(Warn, format, v...)
}

func (l *lvlWrap) Warn(v ...interface{}) {
	l.println(Info, v...)
}

func (l *lvlWrap) Errorf(format string, v ...interface{}) {
	l.printf(Error, format, v...)
}

func (l *lvlWrap) Error(v ...interface{}) {
	l.println(Error, v...)
}

func (l *lvlWrap) printf(lvl int, format string, v ...interface{}) {
	if !l.canPrint(lvl) {
		return
	}
	l.Logger.Println(l.lvlFormat(lvl, format, v...))
}

func (l *lvlWrap) println(lvl int, v ...interface{}) {
	if !l.canPrint(lvl) {
		return
	}
	l.Logger.Println(l.lvlFormat(lvl, "%s", v...))
}

func MakeLogger(lvl string) {
	lvl = strings.ToLower(lvl)
	for k, v := range str2lvl {
		lvl2str[v] = strings.ToUpper(k)
	}

	lo := new(lvlWrap)
	lo.Logger = log.New(os.Stdin, "", log.LstdFlags)
	if l, ok := str2lvl[lvl]; ok {
		lo.l = l
	} else {
		lo.l = 2
	}

	Logger = lo
}
