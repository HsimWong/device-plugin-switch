package dpmanager

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/HsimWong/device-plugin-switch/deviceinstance"
	"github.com/HsimWong/device-plugin-switch/deviceplugin"
	"github.com/HsimWong/device-plugin-switch/utils"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
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

func (d *DpManager) registerNewDeviceCate(RegisReq DeviceRegisterRequest,
	deviceCateID string) string {
	ret := ""
	switch RegisReq.RegisterType {
	case "DeviceType":
		//log.Debugf("uuid:%s", deviceCateID)
		deviceMsg := utils.NewSyncMessenger()
		d.msgRouter.NewClient("device", deviceCateID, deviceMsg)
		//d.msgRouter.NewClient("device"+deviceCateID, deviceMsg)
		d.deviceManagers[deviceCateID] =
			deviceinstance.NewDevice(deviceCateID, deviceMsg)
		//log.Debugf("Router established, start adding new group")
		ret += d.deviceManagers[deviceCateID].AddGroup(RegisReq.AccessPoint,
			RegisReq.DeviceBlockNum)
		//log.Debugf("AddingGroup finished")

		dpMessenger := utils.NewSyncMessenger()
		d.msgRouter.NewClient("devicePlugin", deviceCateID, dpMessenger)
		//d.msgRouter.NewClient("devicePlugin"+deviceCateID, dpMessenger)
		d.devicePlugins[deviceCateID] =
			deviceplugin.NewDevicePluginInstance(RegisReq.DeviceCategoryType,
				deviceCateID, dpMessenger)

		//log.Debugf("Device Established, trying to fire up")
		go d.deviceManagers[deviceCateID].Run()
		go d.devicePlugins[deviceCateID].Run()
		break
	default:
		ret += fmt.Sprintf("TryRegister%sToNonExistdevice%s;",
			strings.Title(RegisReq.RegisterType),
			strings.Title(RegisReq.DeviceCategoryType))
		break
	}
	return ret
}

func (d *DpManager) regisExistDeviceCate(RegisReq DeviceRegisterRequest,
	deviceCateID string) string {
	//log.Debugf("Device Exists, won't be re-registered")
	switch RegisReq.RegisterType {
	case "Group":
		return d.deviceManagers[deviceCateID].AddGroup(RegisReq.AccessPoint,
			RegisReq.DeviceBlockNum)
	case "Block":
		return d.deviceManagers[deviceCateID].AddBlock(RegisReq.AccessPoint,
			RegisReq.DeviceBlockNum)
	case "DeviceType":
		return "DeviceTypeAlreadyExist"
	default:
		return "RegisterTypeNotExist;"
	}
}

func (d *DpManager) Register(RegisReq DeviceRegisterRequest) string {

	deviceCateID := fmt.Sprintf("%x",
		md5.Sum([]byte(RegisReq.DeviceCategoryType)))
	//log.Debugf("deviceCate: %s", deviceCateID)
	retMsg := ""
	if _, exist := d.devicePlugins[deviceCateID]; !exist {
		//log.Debugf("Registering New Device")
		retMsg = d.registerNewDeviceCate(RegisReq, deviceCateID)
	} else {
		retMsg = d.regisExistDeviceCate(RegisReq, deviceCateID)
	}
	if len(retMsg) == 0 {
		retMsg = fmt.Sprintf("Success,%s", deviceCateID)
	}
	return retMsg
}

func (d *DpManager) ProcessRequest(conn *net.Conn) {
	//log.Debugf("Process Request StartWorking...")
	//content, err := ioutil.ReadAll(*conn)
	//
	//content, err := bufio.NewReader(*conn).ReadString('`')

	decoder := json.NewDecoder(*conn)

	var message MessagePackage

	err := decoder.Decode(&message)

	utils.Check(err, "Unmarshalling failed")

	//log.Debugf(message.Type)
	switch message.Type {
	case "Register":
		go func() {
			log.Info("Start registering")
			regResult := d.Register(message.Info)
			log.Infof("RegisterResult:%s", regResult)

			_, err := (*conn).Write([]byte(regResult))
			utils.Check(err, "Returning message failed")
			log.Infof("RegisterResult:%s, has been written back", regResult)

		}()
		break
	default:
		log.Warning("Register type not exist")
	}
}

func (d *DpManager) Run() {
	//log.Debugf("DPmanager start running")
	go utils.StartJsonSvr("tcp4",
		"0.0.0.0"+utils.DpManagerPort, d.ProcessRequest)
}
