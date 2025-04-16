package logutil

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLog(name string) {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	log.SetReportCaller(true)

	log.AddHook(&FunctionHook{})
}

func GetLogger() *logrus.Logger {
	return log
}

type FunctionHook struct{}

func (hook *FunctionHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *FunctionHook) Fire(entry *logrus.Entry) error {
	if pc, _, _, ok := runtime.Caller(8); ok {
		fn := runtime.FuncForPC(pc)
		entry.Data["function"] = fn.Name()
	}
	return nil
}
