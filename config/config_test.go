package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	appconfig, err := LoadConfig("..\\app.conf")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf(appconfig.KafkaAddr)
	t.Log(appconfig.Collectpath)

}
