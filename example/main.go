package main

import "github.com/smtc/glog"

func main() {
	logInit(false)

	glog.Info("init all module successfully.\n")
	glog.Info("program serve at port 8080....\n")
	glog.Warn("this is a warning\n")
	glog.Error("you should pass a param as type %s but %s\n", "struct", "function")

	testfn()

	glog.Close()
}
