package kafka

import (
	"os"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/movio/kasper"
	alg "github.com/senadi/stream-cardinality-excercise/alg"
	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

var done int

// getKey - Fast extract value from json
func getKey(json string, key string) string {
	return gjson.Get(json, key).String()
}

// Process - Callback
func (up *UserProcessor) Process(msgs []*sarama.ConsumerMessage, sender kasper.Sender) error {
	for _, msg := range msgs {
		// For benchmarking purposes
		if done <= 5e6 {
			up.ProcessMessage(msg)
			done++
		} else {
			os.Exit(0)
		}
	}
	return nil
}

// UserProcessor -  Kafka message processor that parses the json data
// and forwards the time and user identitifer info to a counter structure.
type UserProcessor struct {
	Counter *alg.Counter
}

// ProcessMessage -  Processes Kafka messages from topic prints info to stdout
func (up *UserProcessor) ProcessMessage(msg *sarama.ConsumerMessage) {

	msgValue := msg.Value
	uid := getKey(string(msgValue), "uid")
	ts := getKey(string(msgValue), "ts")

	epoh, err := strconv.Atoi(ts)
	if err != nil {
		log.WithField("ts", epoh).Error("Failed to parse timestamp")
	}
	// When one minute frame is done just output the values.
	ok, count := up.Counter.Count(uid, epoh)

	if ok {
		log.WithFields(log.Fields{"count": count}).Info("Counting...")
	}
}
