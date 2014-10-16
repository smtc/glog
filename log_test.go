package glog

import (
	"testing"
)

func TestConsoleLog(t *testing.T) {
	InitLogger(DEV, nil)
	Debug("this is a debug info\n")
	Info("this is a info %s", "logger init successfully.\n")
	Warn("this is a warning: base value should not be %d\n", 0)
	Error("this is a error log\n")
}
