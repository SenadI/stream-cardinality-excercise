package main

import (
	"os"

	sce "github.com/senadi/stream-cardinality-excercise"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&sce.KafkaAddress, "address", "a", sce.KafkaAddress, "Kafka address")
	rootCmd.PersistentFlags().StringVarP(&sce.KafkaTopic, "topic", "t", sce.KafkaTopic, "Kafka topic")
}

var rootCmd = &cobra.Command{
	Use:   sce.Name,
	Short: sce.Description,
	Long:  sce.LongDescription,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		sce.ConfigureLogging()
		if sce.KafkaAddress == "" || sce.KafkaTopic == "" {
			log.Error("Kafka address and topic are required params")
			os.Exit(1)
		}
	},
}
