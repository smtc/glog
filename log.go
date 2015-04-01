package glog

import (
	"fmt"
	"log"
)

type logType int

const (
	DEV logType = iota
	PRO
)
const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	LevelCount
)

var (
	_       = fmt.Printf
	_logger logger
)

type logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Panic(format string, v ...interface{})

	Prefix(lv int) string
	SetPrefix(lv int, prefix string)
	Flags() int
	SetFlags(flag int)
	Close()
}

func InitLogger(level logType, options map[string]interface{}) {
	if level == DEV {
		_logger = console{
			prefixes: map[int]string{
				DebugLevel: "DEBUG",
				InfoLevel:  "INFO",
				WarnLevel:  "WARN",
				ErrorLevel: "ERROR",
				FatalLevel: "FATAL",
				PanicLevel: "PANIC",
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
		//case "nsq":
		//	_logger = createNsqLogger(options)
		default:
			_logger = console{}
		}
	}
}

func Close() {
	_logger.Close()
}

func Flags() int {
	return _logger.Flags()
}

func SetFlags(flag int) {
	_logger.SetFlags(flag)
}

// Prefix returns the output prefix for the standard logger.
func Prefix(lv int) string {
	return _logger.Prefix(lv)
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(lv int, prefix string) {
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

func Panic(format string, v ...interface{}) {
	_logger.Panic(format, v...)
}

// 为了简单，这里修改prefix时就不加锁了
type console struct {
	prefixes map[int]string
}

func (c console) Prefix(lv int) string {
	return c.prefixes[lv]
}

func (c console) SetPrefix(lv int, prefix string) {
	c.prefixes[lv] = prefix
}

func (c console) Flags() int {
	return log.Flags()
}

func (c console) SetFlags(flag int) {
	log.SetFlags(flag)
}

func (c console) Debug(format string, v ...interface{}) {
	log.Printf(c.prefixes[DebugLevel]+" "+format, v...)
}

func (c console) Info(format string, v ...interface{}) {
	log.Printf(c.prefixes[InfoLevel]+" "+format, v...)
}

func (c console) Warn(format string, v ...interface{}) {
	log.Printf(c.prefixes[WarnLevel]+" "+format, v...)
}

func (c console) Error(format string, v ...interface{}) {
	log.Printf(c.prefixes[ErrorLevel]+" "+format, v...)
}

func (c console) Fatal(format string, v ...interface{}) {
	log.Fatalf(c.prefixes[FatalLevel]+" "+format, v...)
}

func (c console) Panic(format string, v ...interface{}) {
	log.Panicf(c.prefixes[PanicLevel]+" "+format, v...)
}

func (c console) Close() {
}
