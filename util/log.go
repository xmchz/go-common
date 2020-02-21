package util

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)


func InitLog(name string, logPath string) {
	initLogPath(logPath)
	errWriter, _ := newRotateLogWriter(name+ ".error", logPath)
	infoWriter, _ := newRotateLogWriter(name+ ".info", logPath)
	log.SetReportCaller(true)
	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  infoWriter,
			log.ErrorLevel: errWriter,
		},
		&log.TextFormatter{},
	))
	log.Infof("log init success, logPath: %s", logPath)
}

func NewLogger(name string, logPath string, formatter log.Formatter) *log.Logger {
	initLogPath(logPath)
	writer, _ := newRotateLogWriter(name, logPath)
	Logger := log.New()
	if formatter == nil {
		formatter = &log.TextFormatter{}
	}
	Logger.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  writer,
			log.ErrorLevel: writer,
		},
		formatter,
	))
	return Logger
}

func initLogPath(logPath string) {
	if _, err := os.Stat(logPath); err != nil {
		if err := os.Mkdir(logPath, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func newRotateLogWriter(name string, logPath string) (*rotatelogs.RotateLogs, error) {
	return rotatelogs.New(
		logPath+string(os.PathSeparator)+name+".%Y%m%d.log",
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)
}
