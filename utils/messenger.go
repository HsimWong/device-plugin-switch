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
		Requester: make(chan interface{}, 0x1),
		Responder: make(chan interface{}, 0x1),
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

//
//func (m *MessageRouter) run() {
//
//}

func GetMessageRouter() *MessageRouter {
	once.Do(func() {
		messageReceivers := make(map[string]*SyncMessenger)
		MsgRouter = &MessageRouter{messageReceivers: messageReceivers}
		//go MsgRouter.run()
	})
	return MsgRouter
}

func (m *MessageRouter) NewClient(target string, deviceID string,
	msgSvr *SyncMessenger) {
	m.messageReceivers[target+deviceID] = msgSvr
}

func (m *MessageRouter) Call(target string, deviceID string,
	content interface{}) interface{} {
	//log.Debugf("calling Target: %s, cateID:%s", target, deviceID)
	if messenger, exist := m.messageReceivers[target+deviceID]; exist {
		//log.Debugf("Messenger: Sending request to %s, ", target+deviceID, content)
		ret := messenger.Request(content)
		log.Debugf("Received Ret:%s", ret)
		return ret
	} else {
		log.Warnf("Target:%s does not exist", target+deviceID)
	}
	return nil
	//return m.messageReceivers[target+deviceID].Request(content)
}
