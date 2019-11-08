package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type Producer struct {
	client sarama.SyncProducer
}

func NewKafkaProducer(addr string) (*Producer, error) {
	// logs.Debug("准备初始化\n")
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	c, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		// logs.Error("init kafka  failed,err:%s", err.Error())
		// fmt.Printf("初始化错误:%s\n", err.Error())
		return nil, err
	}
	// logs.Debug("初始化成功\n")
	return &Producer{client: c}, nil
}

func (p *Producer) Sendmessage(msg, topic string) error {
	m := &sarama.ProducerMessage{}
	m.Value = sarama.StringEncoder(msg)
	m.Topic = topic
	_, _, err := p.client.SendMessage(m)
	if err != nil {
		logs.Error("send msg failed,err:%s", err.Error())
		return err
	}
	logs.Debug("成功发送msg:%s,topic:%s", msg, topic)
	return nil
}
