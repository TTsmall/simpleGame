package log

import "github.com/astaxie/beego/logs"

var (
	Default *logs.BeeLogger
)

func Debug(format string, v ...interface{}) {
	Default.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	Default.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	Default.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	Default.Error(format, v...)
}
