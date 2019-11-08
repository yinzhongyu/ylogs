package main

import (
	"github.com/astaxie/beego/logs"
)

func main() {
	err := serverRun()

	if err != nil {
		logs.Error("serverRUN failed, err:%v", err)
		return
	}

	logs.Info("program exited")
}
