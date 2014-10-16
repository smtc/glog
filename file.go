package glog

import (
	"fmt"
	"github.com/smtc/goutils"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type fileLogger struct {
	Logger
	dir    string
	format string // file suffix, such as "{{program}}-{{host}}-{{username}}-{{yyyy}}{{mm}}{{dd}}-{{HH}}{{MM}}{{SS}}-{{pid}}"

	rtSeconds int64
	rtItems   int64
	rtNbytes  int64
}

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"

	prefixFn = map[int]string{
		DebugLevel: "DEBUG",
		InfoLevel:  "INFO",
		WarnLevel:  "WARN",
		ErrorLevel: "ERROR",
		FatalLevel: "FATAL",
		PanicLevel: "PANIC",
	}
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

	// 使用不同文件来记录不同等级的log，不需要加前缀
	if prefix, ok = options["prefix"].(map[int]string); !ok {
		prefix = nil
	}

	if fnSuffix, ok = options["suffix"].(string); !ok {
		fnSuffix = "-{{yyyy}}{{mm}}{{dd}}-{{HH}}{{MM}}{{SS}}-{{pid}}.log"
	}

	if dir, ok = options["dir"].(string); !ok {
		dir = "./logs"
	} else if dir == "" {
		dir = "./"
	}

	fl := &fileLogger{
		Logger{
			flag: flag,
		},
		dir,
		fnSuffix,
		goutils.ToInt64(options["seconds"], 86400),
		goutils.ToInt64(options["items"], 0),
		goutils.ToInt64(options["nbytes"], 0),
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

	// avoid panic on nil map
	if fl.out.prefix = prefix; prefix == nil {
		fl.out.prefix = make(map[int]string)
	}

	fl.out.out, err = fl.openLogFiles()
	return
}

func (fl *fileLogger) openLogFiles() (wr map[int]io.Writer, err error) {
	var f *os.File

	suffix := formatSuffix(fl.format)
	wr = make(map[int]io.Writer)

	for i := DebugLevel; i < LevelCount; i++ {
		if f, err = os.OpenFile(path.Join(fl.dir, prefixFn[i]+suffix), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666); err != nil {
			return
		}
		wr[i] = f
	}

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

// 2014-10-17 guotie
// TODO: rotate file logs
func (fl *fileLogger) rotate() {
	tm := time.Now().Unix()
	left := fl.rtSeconds - tm%fl.rtSeconds

	time.AfterFunc(time.Duration(left)*time.Second, func() {
		wr, err := fl.openLogFiles()
		if err != nil {
			fl.Error("rotate log files failed: %v\n", err)
			return
		}
		fl.mu.Lock()
		owr := fl.out.out
		fl.out.out = wr
		fl.mu.Unlock()
		for _, r := range owr {
			f := r.(*os.File)
			f.Close()
		}
	})
}
