package main

import (
	"github.com/astaxie/beego/logs"
	"time"
	"ylogs/config"
	"ylogs/kafka"
	"ylogs/tailf"
)

func serverRun() (err error) {
	//加载配置文件
	c, err := config.LoadConfig("..\\app.conf")
	if err != nil {
		return
	}
	//初始化tail
	t, err := tailf.NewTailfManager(c.Collectpath, c.Messagechan)
	if err != nil {
		return
	}
	//初始化kafka
	k, err := kafka.NewKafkaProducer(c.KafkaAddr)

	if err != nil {
		return
	}
	logs.Debug("initialize all succ")

	for {
		msg := t.GetOneLine()
		err := k.Sendmessage(msg.Message, msg.Topic)
		if err != nil {
			logs.Error("send to kafka failed, err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}
