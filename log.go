package glog

import (
	"fmt"
	"log"
)

type logLevel int

const (
	DEV logLevel = iota
	PRO
)

var (
	_       = fmt.Printf
	_logger logger
)

type logger interface {
	Debug(s string, data ...interface{})
	Info(s string, data ...interface{})
	Warn(s string, data ...interface{})
	Error(s string, data ...interface{})
	Fatal(s string, data ...interface{})

	Prefix(lv string) string
	SetPrefix(lv string, prefix string)
	Flags() int
	SetFlags(flag int)
}

func InitLogger(level logLevel, options map[string]interface{}) {
	if level == DEV {
		_logger = console{
			prefixes: map[string]string{
				"debug": "DEBUG",
				"info":  "INFO",
				"warn":  "WARN",
				"error": "ERROR",
				"fatal": "FATAL",
			},
		}
	} else {
		if options == nil {
			_logger = console{}
			return
		}
		switch options["typ"].(string) {
		case "file":
			_logger = createFileLogger(options)
		case "nsq":
			fallthrough
		default:
			_logger = console{}
		}
	}
}
func Flags() int {
	return _logger.Flags()
}

func SetFlags(flag int) {
	_logger.SetFlags(flag)
}

// Prefix returns the output prefix for the standard logger.
func Prefix(lv string) string {
	return _logger.Prefix(lv)
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(lv string, prefix string) {
	_logger.SetPrefix(lv, prefix)
}

func Debug(format string, v ...interface{}) {
	_logger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	_logger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	_logger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	_logger.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	_logger.Fatal(format, v...)
}

// 为了简单，这里修改prefix时就不加锁了
type console struct {
	prefixes map[string]string
}

func (c console) Prefix(lv string) string {
	return c.prefixes[lv]
}

func (c console) SetPrefix(lv string, prefix string) {
	c.prefixes[lv] = prefix
}

func (c console) Flags() int {
	return log.Flags()
}

func (c console) SetFlags(flag int) {
	log.SetFlags(flag)
}

func (c console) Debug(format string, v ...interface{}) {
	log.Printf("DEBUG "+format, v...)
}

func (c console) Info(format string, v ...interface{}) {
	log.Printf("INFO "+format, v...)
}

func (c console) Warn(format string, v ...interface{}) {
	log.Printf("WARN "+format, v...)
}

func (c console) Error(format string, v ...interface{}) {
	log.Printf("ERROR "+format, v...)
}

func (c console) Fatal(format string, v ...interface{}) {
	log.Fatalf("FATAL "+format, v...)
}
