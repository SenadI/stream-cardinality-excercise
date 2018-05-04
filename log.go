package sce

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// ConfigureLogging - Configures basic logging for this excercise
func ConfigureLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
