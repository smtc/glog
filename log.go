package glog

type log_level int

const (
	DEV log_level = iota
	PRO
)

var (
	logger logfile
)

func InitLogger(level log_level, dir string) {
	logger = logfile{Dir: dir}
}

func Debug(s string, data ...interface{}) {
	logger.Debug(s, data)
}

func Info(s string, data ...interface{}) {
	logger.Info(s, data)
}

func Warn(s string, data ...interface{}) {
	logger.Warn(s, data)
}

func Error(s string, data ...interface{}) {
	logger.Error(s, data)
}

func Fatal(s string, data ...interface{}) {
	logger.Fatal(s, data)
}
