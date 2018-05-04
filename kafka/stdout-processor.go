package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/movio/kasper"
	log "github.com/sirupsen/logrus"
)

// Process - Callback
func (cs *ProcessorStdout) Process(msgs []*sarama.ConsumerMessage, sender kasper.Sender) error {
	for _, msg := range msgs {
		cs.ProcessMessage(msg)
	}
	return nil
}

// ProcessorStdout -  Kafka message processor that shows how to read messages from Kafka topic
type ProcessorStdout struct {
}

// ProcessMessage -  Processes Kafka messages from topic prints info to stdout
func (*ProcessorStdout) ProcessMessage(msg *sarama.ConsumerMessage) {
	key := string(msg.Key)
	value := string(msg.Value)
	offset := msg.Offset
	topic := msg.Topic
	partition := msg.Partition
	log.WithFields(log.Fields{"key": key, "value": value, "offset": offset, "topic": topic, "partition": partition}).Info("Message received")
}
