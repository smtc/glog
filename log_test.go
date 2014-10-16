package glog

import (
	//"os"
	"testing"
)

func TestConsoleLog(t *testing.T) {
	InitLogger(DEV, nil)
	Debug("this is a debug info\n")
	Info("this is a info %s", "logger init successfully.\n")
	Warn("this is a warning: base value should not be %d\n", 0)
	Error("this is a error log\n")

	SetPrefix(InfoLevel, "information:")
	Info("modify info log prefix to %s\n", Prefix(InfoLevel))
	Info("no new line")
	Info("new prefix info log\n")
}

func TestFileLog(t *testing.T) {
	InitLogger(PRO, map[string]interface{}{"typ": "file"})
	Debug("this is a debug info\n")
	Info("this is a info %s", "logger init successfully.\n")
	Warn("this is a warning: base value should not be %d\n", 0)
	Error("this is a error log\n")

	SetPrefix(InfoLevel, "information:")
	Info("modify info log prefix to %s\n", Prefix(InfoLevel))
	Info("no new line")
	Info("new prefix info log\n")
}
