package dpmanager

import (
	"encoding/json"
	"github.com/HsimWong/device-plugin-switch/deviceinstance"
	"github.com/HsimWong/device-plugin-switch/deviceplugin"
	"github.com/HsimWong/device-plugin-switch/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
)

type DpManager struct {
	devicePlugins map[string]*deviceplugin.Instance
	// k: deviceCategoryID; v: devicePluginInstance
	msgRouter      *utils.MessageRouter
	deviceManagers map[string]*deviceinstance.DeviceCategory
}

func NewDpManager() *DpManager {
	return &DpManager{
		devicePlugins:  make(map[string]*deviceplugin.Instance),
		msgRouter:      utils.GetMessageRouter(),
		deviceManagers: make(map[string]*deviceinstance.DeviceCategory),
	}
}

func (d *DpManager) Register(RegisReq DeviceRegisterRequest) {
	deviceCateID := uuid.NewString()
	log.Debugf("uuid:%s", deviceCateID)
	dpMessenger := utils.NewSyncMessenger()
	d.msgRouter.NewClient("devicePlugin"+deviceCateID, dpMessenger)
	deviceMsg := utils.NewSyncMessenger()
	d.msgRouter.NewClient("device"+deviceCateID, deviceMsg)
	d.deviceManagers[deviceCateID] =
		deviceinstance.NewDevice(deviceCateID, deviceMsg)
	d.deviceManagers[deviceCateID].AddDevice(RegisReq.AccessPoint, RegisReq.DeviceBlockNum)
	d.devicePlugins[deviceCateID] =
		deviceplugin.NewDevicePluginInstance(RegisReq.DeviceCategoryType,
			deviceCateID, dpMessenger, d.deviceManagers[deviceCateID])

	log.Debugf("Device Established, trying to fire up")
	go d.deviceManagers[deviceCateID].Run()
	// Should be adding devices here `AddDevice`
	go d.devicePlugins[deviceCateID].Run()
}

func (d *DpManager) ProcessRequest(conn *net.Conn) {
	log.Debugf("Process Request StartWorking...")
	content, err := ioutil.ReadAll(*conn)

	utils.Check(err, "Reading socket failed")
	log.Debugf("%s", string(content))
	var message MessagePackage
	err = json.Unmarshal(content, &message)
	utils.Check(err, "Unmarshalling failed")

	log.Debugf(message.Type)
	switch message.Type {
	case "Register":
		go d.Register(message.Info)
		break
	default:
		log.Warning("Register type not exist")
	}
	//log.Debugf("%s", message.Info.(map[string]interface{}))

}

func (d *DpManager) Run() {
	log.Debugf("DPmanager start running")
	go utils.StartJsonSvr("tcp4",
		"0.0.0.0"+utils.DpManagerPort, d.ProcessRequest)

}
