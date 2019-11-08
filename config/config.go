package config

import (
	"errors"
	"github.com/astaxie/beego/config"
	"ylogs/tailf"
)

type Config struct {
	Collectpath []tailf.CollectPath
	KafkaAddr   string
	Messagechan int
}

func LoadConfig(filename string) (*Config, error) {

	appconfig := new(Config)
	conf, err := config.NewConfig("ini", filename)
	if err != nil {
		return nil, err
	}
	appconfig.KafkaAddr = conf.String("kafka::addr")
	appconfig.Messagechan, err = conf.Int("collect::messagechan")
	if err != nil {
		return nil, err
	}
	tpath := conf.Strings("collect::path_name")
	ttopic := conf.Strings("collect::topic_name")

	if len(tpath) != len(ttopic) {
		return nil, errors.New("配置失败：path数量与topic数量应该相等")
	}
	if len(tpath) == 0 || len(ttopic) == 0 {
		return nil, errors.New("配置加载失败：path与topic数量不应为0\n")
	}

	tcollectpath := tailf.CollectPath{}
	for i, _ := range tpath {
		tcollectpath.Pathname = tpath[i]
		tcollectpath.Topic = ttopic[i]
		appconfig.Collectpath = append(appconfig.Collectpath, tcollectpath)
	}
	return appconfig, nil

}
