package logger

import "github.com/sirupsen/logrus"

func Debug(msg ...interface{}) {
	logrus.Debug(msg...)
}

func Debugf(template string, args ...interface{}) {
	logrus.Debugf(template, args...)
}

func Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func Infof(template string, args ...interface{}) {
	logrus.Infof(template, args...)
}

func Warn(msg ...interface{}) {
	logrus.Warn(msg...)
}

func Warnf(template string, args ...interface{}) {
	logrus.Warnf(template, args...)
}

func Error(msg ...interface{}) {
	logrus.Error(msg...)
}

func Errorf(template string, args ...interface{}) {
	logrus.Errorf(template, args...)
}
