package main

import (
	"fmt"
	"github.com/HsimWong/device-plugin-switch/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

func Thread1() {
	log.SetLevel(log.DebugLevel)
	router := utils.GetMessageRouter()
	for {
		for i := 0; i < 10; i++ {
			recvResp := router.Call("thread2",
				"helloTest", fmt.Sprintf("HelloFromThread1+%d", i)).(string)
			log.Infof("Received response: %s", recvResp)
			time.Sleep(1 * time.Second)
		}

	}
}

func Thread2(msgSvr *utils.SyncMessenger) {
	for {
		//router := utils.GetMessageRouter()
		recv := msgSvr.Serve().(string)
		msgSvr.Respond("Receive from thread 2" + recv)
	}
}

func main() {
	msgSvr := utils.NewSyncMessenger()
	router := utils.GetMessageRouter()
	router.NewClient("thread2", "helloTest", msgSvr)
	go Thread2(msgSvr)
	go Thread1()
	utils.ThreadBlock()
}
