package helper

import (
	"fmt"
	"time"

	"github.com/buguang01/Logger"
)

//可以挂载到别的框架中
func NewLoggerIO(loglv Logger.LogLevel) *LoggerIO {
	return &LoggerIO{
		LogLv: loglv,
	}
}

//一般第三方的框架都是可以设置一个支持io.Write接口的对象传入就可以了.test
type LoggerIO struct {
	LogLv Logger.LogLevel
}

func (this *LoggerIO) Write(p []byte) (n int, err error) {
	Logger.PAlart(this.LogLv, string(p))
	return 0, nil
}

func (this *LoggerIO) Print(values ...interface{}) {
	Logger.PrintLog(&Logger.LogMsgModel{
		Msg:   fmt.Sprint(values...),
		LogLv: this.LogLv,
		Stack: "",
		KeyID: -1,
	})
}

func GetDate(d time.Time) time.Time {
	result := time.Date(
		d.Year(),
		d.Month(),
		d.Day(),
		0,
		0,
		0,
		0,
		time.Local)
	return result
}
