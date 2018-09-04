package main

import (
	"github.com/sirupsen/logrus"
	"test.com/mine/services/initializer"
	"test.com/mine/services/signal"
)

func main() {
	defer initializer.Initialize()()
	sig := signal.WaitExitSignal()
	logrus.Infof("Sig %s received, exiting...", sig)
}
