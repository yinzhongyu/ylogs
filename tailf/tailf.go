package tailf

import (
	"context"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"os"
	"time"
)

//一个收集路径
type CollectPath struct {
	Pathname string
	Topic    string
}

type KafkaMessage struct {
	Topic   string
	Message string
}

//一个tailf结构
type TailfObj struct {
	collectpath CollectPath
	exit        chan bool
	tail        *tail.Tail
}

//tailf管理器
type TailfManager struct {
	tailobjs   []*TailfObj
	ctx        context.Context
	cancelfunc context.CancelFunc
	msgchan    chan *KafkaMessage
}

func NewTailfManager(conf []CollectPath, resultsize int) (*TailfManager, error) {

	tailmgr := &TailfManager{
		msgchan: make(chan *KafkaMessage, resultsize),
	}

	tailmgr.ctx, tailmgr.cancelfunc = context.WithCancel(context.Background())
	err := tailmgr.init(conf)
	if err != nil {
		return nil, err
	}
	return tailmgr, nil
}

func (tm *TailfManager) init(conf []CollectPath) error {

	if len(conf) == 0 {
		return errors.New("invalid config for log collect, conf:%v")
	}

	for _, v := range conf {
		tm.newTailTask(v)
	}

	tm.readFromTailObj()

	return nil

}

func (tm *TailfManager) newTailTask(c CollectPath) {
	obj := &TailfObj{
		collectpath: c,
		exit:        make(chan bool, 1),
	}

	tails, err := tail.TailFile(c.Pathname, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		logs.Error("collect filename[%s] failed, err:%v", c.Pathname, err)
		return
	}

	obj.tail = tails
	tm.tailobjs = append(tm.tailobjs, obj)
}

func (tm *TailfManager) GetOneLine() *KafkaMessage {
	msg := <-tm.msgchan
	return msg
}

func (tm *TailfManager) readFromTailObj() {

	for _, v := range tm.tailobjs {

		go func(v *TailfObj) {
			for {
				select {
				case line, ok := <-v.tail.Lines:
					if !ok {
						logs.Warn("tail file close reopen, filename:%s\n", v.tail.Filename)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					textMsg := &KafkaMessage{
						Message: line.Text,
						Topic:   v.collectpath.Topic,
					}

					tm.msgchan <- textMsg
				case <-v.exit:
					logs.Warn("tail obj will exited, conf:%v", v.collectpath)
					return
				case <-tm.ctx.Done():
					logs.Warn("tail obj will exited, conf:%v", v.collectpath)
					return
				}
			}
		}(v)

	}
}
