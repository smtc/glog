package glog

import (
//	"time"
)

type fileLogger struct {
	Dir string
}

func createFileLogger(options map[string]interface{}) *fileLogger {
	fl := fileLogger{}
	return &fl
}

func (l *fileLogger) Debug(s string, data ...interface{}) {
}

func (l *fileLogger) Info(s string, data ...interface{}) {
}

func (l *fileLogger) Warn(s string, data ...interface{}) {
}

func (l *fileLogger) Error(s string, data ...interface{}) {
}

func (l *fileLogger) Fatal(s string, data ...interface{}) {
}
