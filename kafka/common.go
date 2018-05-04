package kafka

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/movio/kasper"
	log "github.com/sirupsen/logrus"
)

// NewKafkaConfig - preparse the client and configuration for TopicProcessor
func NewKafkaConfig(address string, topic string, partition int) *kasper.Config {
	client, _ := sarama.NewClient([]string{address}, sarama.NewConfig())
	config := &kasper.Config{
		TopicProcessorName: "example-processor",
		Client:             client,
		InputTopics:        []string{topic},
		InputPartitions:    []int{partition},
	}
	return config
}

// HandleSignals - closes topic processor on sigint/sigterm
func HandleSignals(tp *kasper.TopicProcessor) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Topic processor is running...")
	for range signals {
		signal.Stop(signals)
		tp.Close()
		return
	}
}
