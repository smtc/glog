package glog

type logfile struct {
	Dir string
}

func (l *logfile) Debug(s string, data ...interface{}) {
}

func (l *logfile) Info(s string, data ...interface{}) {
}

func (l *logfile) Warn(s string, data ...interface{}) {
}

func (l *logfile) Error(s string, data ...interface{}) {
}

func (l *logfile) Fatal(s string, data ...interface{}) {
}
