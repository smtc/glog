package glog

import (
	"time"
)

type fileLogger struct {
	Logger
	dir    string
	format string // file suffix, such as "{{pid}}-{{yyyy}}-{{}}"
}

// options:
//    flag: int
//    prefix: map[int]string
//    dir: string
//    format: string
//
func createFileLogger(options map[string]interface{}) *fileLogger {
	var (
		ok       bool
		err      error
		flag     int
		dir      string
		fnSuffix string
		prefix   map[int]string
	)

	if flag, ok = options["flag"].(int); !ok {
		flag = 0
	}
	if prefix, ok = options["prefix"].(map[int]string); !ok {
		prefix = map[int]string{
			DebugLevel: "DEBUG",
			InfoLevel:  "INFO",
			WarnLevel:  "WARN",
			ErrorLevel: "ERROR",
			FatalLevel: "FATAL",
			PanicLevel: "PANIC",
		}
	}

	fl := fileLogger{
		Logger{
			flag: flag,
		},
		dir,
		fnSuffix,
	}
	fl.Logger.out, err = buildFileOut()
	if err != nil {
		panic(err)
	}

	return &fl
}

func buildFileOut() (o outputer, err error) {
	return
}

func openLogFile() {

}
