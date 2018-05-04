package kafka

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/movio/kasper"
	counter "github.com/paulbellamy/ratecounter"
	log "github.com/sirupsen/logrus"
)

// Process - Callback
func (f *FrameProcessor) Process(msgs []*sarama.ConsumerMessage, sender kasper.Sender) error {
	for _, msg := range msgs {
		f.ProcessMessage(msg)
	}
	return nil
}

// FrameProcessor  - Kafka message processor that messures "frames-per-(duration)"
type FrameProcessor struct {
	Counter    *counter.RateCounter
	DecodeJSON int
}

// ProcessMessage -  Processes Kafka messages from topic and increments the counter
func (f *FrameProcessor) ProcessMessage(msg *sarama.ConsumerMessage) {
	if f.DecodeJSON != 0 {
		var decoded map[string]*json.RawMessage
		err := json.Unmarshal(msg.Value, &decoded)
		if err != nil {
			log.Error("Error occured during decoding", err)
		}
	}
	// Lets record the event
	f.Counter.Incr(1)
}

// logs the rate
func logRate(f *FrameProcessor) {
	log.WithField("frames", f.Counter.Rate()).Info("Number of frames per duration")
}

// LogRate - Goroutine that writes the counter content every so often.
func LogRate(f *FrameProcessor, frequency time.Duration) {
	for {
		logRate(f)
		time.Sleep(frequency)
	}
}
