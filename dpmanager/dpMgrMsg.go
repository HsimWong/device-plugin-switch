package dpmanager

import "sync"

type MessagePackage struct {
	Type string                `json:"Type"`
	Info DeviceRegisterRequest `json:"Info"`
}

type DeviceStatus struct {
	deviceID string
	status   string
}

type DeviceRegisterRequest struct {
	DeviceCategoryType string
	DeviceBlockNum     int
	AccessPoint        string
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
