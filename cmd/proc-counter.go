package main

import (
	"time"

	counter "github.com/paulbellamy/ratecounter"
	sce "github.com/senadi/stream-cardinality-excercise"
	kafka "github.com/senadi/stream-cardinality-excercise/kafka"

	log "github.com/sirupsen/logrus"

	"github.com/movio/kasper"
	"github.com/spf13/cobra"
)

var rate time.Duration
var outputRate time.Duration
var decodeJSON int

func init() {
	frameCounterCmd.Flags().DurationVarP(&rate, "rate", "r", rate, "Frames per rate")
	frameCounterCmd.Flags().DurationVarP(&outputRate, "output-rate", "o", outputRate, "Results output rate")
	frameCounterCmd.Flags().IntVarP(&decodeJSON, "json", "j", decodeJSON, "Decode json")
	rootCmd.AddCommand(frameCounterCmd)
}

var frameCounterCmd = &cobra.Command{
	Use:   "processor-counter",
	Short: "Process messages and outputs the frame-per-rate stats",
	Long:  "Process kafka messages and outputs the frame-per-rate to os.Stdout",
	Run: func(cmd *cobra.Command, args []string) {

		config := kafka.NewKafkaConfig(sce.KafkaAddress, sce.KafkaTopic, 0)
		counter := counter.NewRateCounter(rate)

		processor := &kafka.FrameProcessor{Counter: counter, DecodeJSON: decodeJSON}

		messageProcessors := map[int]kasper.MessageProcessor{0: processor}
		tp := kasper.NewTopicProcessor(config, messageProcessors)

		// Stop processor on SIGINT/SIGTERM
		go kafka.HandleSignals(tp)

		// Output the results based on outputRate
		go kafka.LogRate(processor, outputRate)

		err := tp.RunLoop()
		log.WithField("error", err).Error("Topic processor finished with error")
	},
}
