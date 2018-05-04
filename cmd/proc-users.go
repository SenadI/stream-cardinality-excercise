package main

import (
	sce "github.com/senadi/stream-cardinality-excercise"
	"github.com/senadi/stream-cardinality-excercise/alg"
	kafka "github.com/senadi/stream-cardinality-excercise/kafka"

	log "github.com/sirupsen/logrus"

	"github.com/movio/kasper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(userCounterCmd)
}

var userCounterCmd = &cobra.Command{
	Use:   "processor-user",
	Short: "Process messages and count unique users per minute",
	Long:  "Process kafka messages and count unique users per minute and outputs them to os.Stdout",
	Run: func(cmd *cobra.Command, args []string) {

		config := kafka.NewKafkaConfig(sce.KafkaAddress, sce.KafkaTopic, 0)

		counter := &alg.Counter{Hll: nil, Now: 0}
		messageProcessors := map[int]kasper.MessageProcessor{0: &kafka.UserProcessor{Counter: counter}}
		tp := kasper.NewTopicProcessor(config, messageProcessors)

		// Stop processor on SIGINT/SIGTERM
		go kafka.HandleSignals(tp)
		err := tp.RunLoop()
		log.WithField("error", err).Error("Topic processor finished with error")
	},
}
