package main

import (
	"fmt"
	"log"
	"time"
)

func configureDefaultLogger() {
	logger := log.Default()
	logger.SetFlags(0)
	logger.SetOutput(new(logWriter))
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print("[IPCHECK] " + time.Now().UTC().Format("2006/01/02 15:04:05") + " | " + string(bytes))
}
