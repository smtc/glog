package glog

import (
	"fmt"
	"github.com/smtc/goutils"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type fileLogger struct {
	Logger
	dir    string
	format string // file suffix, such as "{{program}}-{{host}}-{{username}}-{{yyyy}}{{mm}}{{dd}}-{{HH}}{{MM}}{{SS}}-{{pid}}"
}

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
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
	if fnSuffix, ok = options["suffix"].(string); !ok {
		fnSuffix = "-{{yyyy}}{{mm}}{{dd}}-{{HH}}{{MM}}{{SS}}-{{pid}}.log"
	}

	fl := &fileLogger{
		Logger{
			flag: flag,
		},
		dir,
		fnSuffix,
	}
	err = fl.buildFileOut(prefix)
	if err != nil {
		panic(err)
	}

	return fl
}

func (fl *fileLogger) buildFileOut(prefix map[int]string) (err error) {
	if err = goutils.CreateDirIfNotExist(fl.dir); err != nil {
		return
	}

	return
}

func (fl *fileLogger) openLogFile() (err error) {
	return
}

func formatSuffix(format string) (res string) {
	if format == "" {
		return
	}

	tm := time.Now()
	res = strings.Replace(format, "{{program}}", program, -1)
	res = strings.Replace(res, "{{host}}", host, -1)
	res = strings.Replace(format, "{{username}}", userName, -1)
	res = strings.Replace(res, "{{yyyy}}", fmt.Sprintf("%d", tm.Year()), -1)
	res = strings.Replace(res, "{{mm}}", fmt.Sprintf("%02d", tm.Month()), -1)
	res = strings.Replace(res, "{{dd}}", fmt.Sprintf("%02d", tm.Day()), -1)
	res = strings.Replace(res, "{{HH}}", fmt.Sprintf("%02d", tm.Hour()), -1)
	res = strings.Replace(res, "{{MM}}", fmt.Sprintf("%02d", tm.Minute()), -1)
	res = strings.Replace(res, "{{SS}}", fmt.Sprintf("%02d", tm.Second()), -1)
	res = strings.Replace(res, "{{pid}}", fmt.Sprint(pid), -1)

	return
}
