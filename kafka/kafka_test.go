package kafka

import (
	"testing"
)

func TestKafka(t *testing.T) {
	// t.Logf("hello")
	// fmt.Printf("hello")
	p := NewKafkaProducer("0.0.0.0:9092")
	err := p.Sendmessage("hello!", "test")
	if err != nil {
		return
	}

}
