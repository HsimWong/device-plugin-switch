package utils

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

type SyncMessenger struct {
	Requester chan interface{}
	Responder chan interface{}
}

func NewSyncMessenger() *SyncMessenger {
	return &SyncMessenger{
		Requester: make(chan interface{}),
		Responder: make(chan interface{}),
	}
}

func (s *SyncMessenger) Request(reqMsg interface{}) interface{} {
	s.Requester <- reqMsg
	return <-s.Responder
}

func (s *SyncMessenger) Serve() interface{} {
	return <-s.Requester
}

func (s *SyncMessenger) Respond(resMsg interface{}) {
	s.Responder <- resMsg
}

var MsgRouter *MessageRouter
var once sync.Once

type MessageRouter struct {
	messageReceivers map[string]*SyncMessenger
}

func (m *MessageRouter) run() {

}

func GetMessageRouter() *MessageRouter {
	once.Do(func() {
		messageReceivers := make(map[string]*SyncMessenger)
		MsgRouter = &MessageRouter{messageReceivers: messageReceivers}
		go MsgRouter.run()
	})
	return MsgRouter
}

func (m *MessageRouter) NewClient(pathID string, msgSvr *SyncMessenger) {
	m.messageReceivers[pathID] = msgSvr
}

func (m *MessageRouter) Call(target string, deviceID string,
	content interface{}) interface{} {
	log.Debugf("Target: %s, cateID:%s", target, deviceID)
	return m.messageReceivers[target+deviceID].Request(content)
}
