package dpmanager

import "sync"

type MessageContent interface{}

type MessagePackage struct {
	Type string         `json:"Type"`
	Info MessageContent `json:"Info"`
}

type DeviceStatus struct {
	deviceID string
	status   string
}

type DeviceRegisterRequest struct {
	DeviceCategoryType string
	DeviceBlockNum     int
	AccessPoint        string
	RegisterType       string // DeviceType, Group, Block
}

type DeviceRegisterResponse struct {
	deviceRegisterID      string
	deviceRegisterResults string
}

type DeviceStatusDelta struct {
}

var messageTypeMapper *map[string]interface{}
var once sync.Once

func getMsgTypeMapper() *map[string]interface{} {
	once.Do(func() {
		*messageTypeMapper = make(map[string]interface{})
	})
	return messageTypeMapper
}
