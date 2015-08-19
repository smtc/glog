package glog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/smtc/goutils"
)

type fileLogger struct {
	Logger
	dir    string
	format string // file suffix, such as "{{program}}-{{host}}-{{username}}-{{yyyy}}{{mm}}{{dd}}-{{HH}}{{MM}}{{SS}}-{{pid}}"

	rtSeconds int64
	rtItems   int64
	rtNbytes  int64
	sequence  int
	natureDay bool // 自然日模式

	rotDuration time.Duration
	rot         chan struct{}
	exist       chan struct{}
	existed     chan bool
	contact     bool
}

const (
	HourlyDuration = (time.Duration(3600) * time.Second)
	DailyDuration  = (time.Duration(86400) * time.Second)
)

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

func cleanTmpLogs(dir string, contact bool) {
	var (
		err error
		fns []string = make([]string, 0)
	)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".tmp") {
			fns = append(fns, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("clean tmp logs failed: walk dir %s failed: %v\n", dir, err)
	}

	for _, fn := range fns {
		renameTmpLogs(dir, fn, contact)
	}

	return
}

func renameTmpLogs(dir, fn string, contact bool) {
	if len(fn) <= 4 {
		return
	}

	// 不断的尝试rename fn - ".tmp" + xxx + ".log", 直到成功为止
	var (
		seq    int = 1
		err    error
		baseFn = filepath.Base(fn)
		nfn    string
	)

	if len(baseFn) < 4 {
		fmt.Printf("tmp log filename %s %s invalid: filename should contain .tmp\n", fn, baseFn)
		return
	}
	baseFn = baseFn[0 : len(baseFn)-4]

	if contact {
		contactLog(path.Join(dir, baseFn+".log"), fn)
	} else {
		seq = checkSequence(dir, baseFn) + 1
		nfn = path.Join(dir, baseFn+fmt.Sprintf(".seq%03d", seq)+".log")

		if err = os.Rename(fn, nfn); err != nil {
			fmt.Printf("Cannot rename tmp log file %s, remove it\n", fn, err)
			os.Remove(fn)
		}
	}

	return
}

func contactLog(log, tmp string) {
	file, err := os.OpenFile(log, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm|os.ModeTemporary)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tmpFile, err := os.Open(tmp)
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	buff, _ := ioutil.ReadAll(tmpFile)
	file.Write(buff)
	os.Remove(tmp)
}

// 检查是否由以tmp结尾的文件, 可能有雨程序崩溃, 存在tmp结尾的文件，把这些文件转换为对应的log后缀
func checkSequence(dir, fnDate string) int {
	var (
		err      error
		seq      int
		sequence int
		fns      []string = make([]string, 0)
	)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		baseFn := filepath.Base(path)
		if strings.Contains(baseFn, fnDate) {
			fns = append(fns, baseFn)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("check sequence: walk dir %s failed: %v\n", dir, err)
	}

	for _, fn := range fns {
		segs := strings.Split(fn, ".")
		if len(segs) < 2 {
			continue
		}
		sseq := segs[len(segs)-2]
		if strings.HasPrefix(sseq, "seq") == false {
			continue
		}
		seq, err = strconv.Atoi(sseq[3:])
		if err != nil {
			fmt.Printf("convert sequence %s failed: %v\n", sseq[3:], err)
			continue
		}
		if sequence < seq {
			sequence = seq
		}
	}

	return sequence
}

// options:
//    flag: int
//    prefix: map[int]string
//    dir: string
//    contcat:
//    format: string
//
func createFileLogger(options map[string]interface{}) *fileLogger {
	var (
		ok       bool
		err      error
		flag     int
		sequence int
		dir      string
		fnSuffix string
		contact  bool
		prefix   map[int]string
	)

	if flag, ok = options["flag"].(int); !ok {
		flag = Ldate | Ltime
	}

	// 使用不同文件来记录不同等级的log，不需要加前缀
	if prefix, ok = options["prefix"].(map[int]string); !ok {
		prefix = nil
	}

	if contact, ok = options["contact"].(bool); !ok {
		contact = false
	}

	if fnSuffix, ok = options["suffix"].(string); !ok {
		fnSuffix = "-{{yyyy}}{{mm}}{{dd}}-{{HH}}{{MM}}{{SS}}-{{pid}}"
	}

	if dir, ok = options["dir"].(string); !ok {
		dir = "./logs"
	} else if dir == "" {
		dir = "./"
	}

	// 判断目录是否存在
	_, err = os.Stat(dir)
	dirExisted := err == nil || os.IsExist(err)
	if !dirExisted {
		os.Mkdir(dir, 0666)
	}

	// 清理tmp log文件
	cleanTmpLogs(dir, contact)
	sequence = checkSequence(dir, formatSuffix(fnSuffix))

	fl := &fileLogger{
		Logger{
			flag: flag,
		},
		dir,
		fnSuffix,
		goutils.ToInt64(options["seconds"], 86400),
		goutils.ToInt64(options["items"], 0),
		goutils.ToInt64(options["nbytes"], 0),
		sequence, // sequence
		true,
		0,
		make(chan struct{}),
		make(chan struct{}),
		make(chan bool),
		contact,
	}

	//fmt.Println("createFileLogger: sequence=", fl.sequence, fl.rtSeconds, options["seconds"])
	if fl.rtSeconds < 5 {
		fl.rtSeconds = 5
	}

	fl.rotDuration = time.Duration(fl.rtSeconds) * time.Second

	err = fl.buildFileOut(prefix)
	if err != nil {
		panic(err)
	}

	fl.rotate()
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

func (fl *fileLogger) openLogFiles() (wr map[int]io.WriteCloser, err error) {
	var (
		f  *os.File
		fn string
	)

	suffix := formatSuffix(fl.format)
	wr = make(map[int]io.WriteCloser)

	for i := DebugLevel; i < LevelCount; i++ {
		fn = path.Join(fl.dir, prefixFn[i]+suffix+".tmp")
		//log.Printf("open log level %d, fn=%s\n", i, fn)
		if f, err = os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666); err != nil {
			log.Printf("open log file %s failed: %v\n", fn, err)
			continue
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
	var left int64

	tm := time.Now().Unix()
	if fl.rtSeconds%86400 == 0 && fl.natureDay {
		// 按自然日生成日志
		_, offset := time.Now().Zone()
		left = fl.rtSeconds - (tm-int64(offset))%fl.rtSeconds
	} else {
		left = fl.rtSeconds - tm%fl.rtSeconds
	}

	go func() {
		for {
			select {
			case <-fl.rot:
				//log.Println("time is up, rotate ....")
				fl.mu.Lock()
				owr := fl.out.out
				fl.closeLogFiles(owr)
				wr, err := fl.openLogFiles()
				if err != nil {
					fl.Error("rotate log files failed: %v\n", err)
				} else {
					fl.out.out = wr
				}
				fl.mu.Unlock()

				tm := time.Now().Unix()
				left := fl.rtSeconds - tm%fl.rtSeconds
				time.AfterFunc(time.Duration(left)*time.Second, func() {
					fl.rot <- struct{}{}
				})
			case <-fl.exist:
				//log.Println("file log exist ....")
				fl.mu.Lock()
				owr := fl.out.out
				fl.closeLogFiles(owr)
				fl.mu.Unlock()
				close(fl.existed)
				return
			}
		}
	}()

	// log.Println("rotate: left=", left)
	time.AfterFunc(time.Duration(left)*time.Second, func() {
		fl.rot <- struct{}{}
	})
}

func (fl *fileLogger) closeLogFiles(fs map[int]io.WriteCloser) {
	var (
		err     error
		fn, nfn string
	)

	fl.sequence++
	for _, r := range fs {
		f := r.(*os.File)
		f.Close()

		if fl.contact {
			contactLog(f.Name()[0:len(f.Name())-4]+".log", f.Name())
		} else {
			if fl.rotDuration < DailyDuration {
				nfn = f.Name()[0:len(f.Name())-4] + fmt.Sprintf(".seq%03d", fl.sequence) + ".log"
			} else {
				nfn = f.Name()[0:len(f.Name())-4] + ".log"
			}
			fn = f.Name()
			if err = os.Rename(fn, nfn); err != nil {
				log.Println("closeLogFiles failed:", err)
			} else {
				//log.Println("closeLogFiles:", fn, nfn)
			}
		}
	}
}

func (fl *fileLogger) Close() {
	fl.exist <- struct{}{}
	<-fl.existed
}
