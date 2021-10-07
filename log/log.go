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
var logFormat = "[%s]:%s"

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

func (l *lvlWrap) Printf(lvl int, format string, v ...interface{}) {
	if !l.canPrint(lvl) {
		return
	}
	l.Logger.Println(l.lvlFormat(lvl, format, v...))
}

func (l *lvlWrap) Print(lvl int, v ...interface{}) {
	if !l.canPrint(lvl) {
		return
	}
	l.Logger.Print(l.lvlFormat(lvl, "%s", v...))
}

func (l *lvlWrap) Println(lvl int, v ...interface{}) {
	if !l.canPrint(lvl) {
		return
	}
	l.Logger.Println(l.lvlFormat(lvl, "%s", v...))
}

func MakeLogger(l string) {
	l = strings.ToLower(l)
	for k, v := range str2lvl {
		lvl2str[v] = strings.ToUpper(k)
	}

	lo := new(lvlWrap)
	lo.Logger = log.New(os.Stdin, "", log.LstdFlags)
	lo.l = str2lvl[l]
	Logger = lo
}
