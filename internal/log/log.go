package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func Init() {
	// 로거 포맷 설정
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 출력 설정
	Logger.SetOutput(os.Stdout)

	// 로그 레벨 설정 (Debug, Info, Warn, Error, Fatal, Panic)
	Logger.SetLevel(logrus.InfoLevel)

	// 호출자 정보 추가
	Logger.SetReportCaller(true)
}

// 편의성을 위한 래퍼 함수들
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return Logger.WithError(err)
}
