package main

import (
	"github.com/HsimWong/device-plugin-switch/dpmanager"
	"github.com/HsimWong/device-plugin-switch/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debugf("Begins")
	dpManager := dpmanager.NewDpManager()
	go dpManager.Run()
	utils.ThreadBlock()
}
