package main

import (
	sce "github.com/senadi/stream-cardinality-excercise"
	kafka "github.com/senadi/stream-cardinality-excercise/kafka"

	log "github.com/sirupsen/logrus"

	"github.com/movio/kasper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(consumerStdoutCmd)
}

var consumerStdoutCmd = &cobra.Command{
	Use:   "processor-stdout",
	Short: "Process messages and outputs them to os.Stdout",
	Long:  "Process kafka messages and outputs them to os.Stdout",
	Run: func(cmd *cobra.Command, args []string) {

		config := kafka.NewKafkaConfig(sce.KafkaAddress, sce.KafkaTopic, 0)
		messageProcessors := map[int]kasper.MessageProcessor{0: &kafka.ProcessorStdout{}}
		tp := kasper.NewTopicProcessor(config, messageProcessors)

		// Stop processor on SIGINT/SIGTERM
		go kafka.HandleSignals(tp)
		err := tp.RunLoop()
		log.WithField("error", err).Error("Topic processor finished with error")
	},
}
