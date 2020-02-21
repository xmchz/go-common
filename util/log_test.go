package util

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	InitLog("test-log-file", "../log")
	go func() {
		for {
			log.Info("this is level info")
			log.Error("this is level error")
			time.Sleep(5 * time.Second)
		}
	}()
	<-time.After(15 * time.Second)
}
