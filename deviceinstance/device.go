package deviceinstance

import (
	"crypto/md5"
	"fmt"
	"github.com/HsimWong/device-plugin-switch/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"time"
)

type deviceController struct {
	partsMap       map[int]string //map[partGroupIndex]partTotalIndex
	accessPoint    string
	BrokenPartsNum int // index starts from 0
	UsedPartsNum   int
	FreePartsNum   int
}

type DeviceCategory struct {
	deviceCategoryID  string
	DevicePartsAmount int
	DeviceParts       map[string]*plugin.Device
	DeviceControllers map[string]*deviceController //map[accessPoint-md5]*deviceController
	deviceMsgSrv      *utils.SyncMessenger
	deviceUpdateChan  chan bool
}

func NewDevice(deviceCategoryID string,
	deviceMsgSrv *utils.SyncMessenger) *DeviceCategory {
	return &DeviceCategory{
		deviceCategoryID:  deviceCategoryID,
		DevicePartsAmount: 0,
		DeviceParts:       make(map[string]*plugin.Device),
		DeviceControllers: make(map[string]*deviceController),
		deviceMsgSrv:      deviceMsgSrv,
		deviceUpdateChan:  make(chan bool, 0xa),
	}
}

func (d *DeviceCategory) Run() {
	//log.Debugf("Device start running")
	go d.deviceStatusResponder()
}

func (d *DeviceCategory) AddGroup(accessPoint string, deviceAmount int) string {
	controllerIndex := fmt.Sprintf("%x", md5.Sum([]byte(accessPoint)))
	if _, exist := d.DeviceControllers[controllerIndex]; exist {
		return fmt.Sprintf("Re-registerExistedGroup:%s;", accessPoint)
	}
	//log.Debugf("Working on AddGroup")
	d.DeviceControllers[controllerIndex] = &deviceController{
		partsMap:       make(map[int]string),
		accessPoint:    accessPoint,
		BrokenPartsNum: deviceAmount,
		UsedPartsNum:   0,
		FreePartsNum:   0,
	}
	//newDeviceParts := make([]*plugin.Device, deviceAmount)
	//d.DeviceParts = append(d.DeviceParts, newDeviceParts...)
	for i := 0; i < deviceAmount; i++ {
		deviceBlockID := uuid.NewString()
		d.DeviceParts[deviceBlockID] = &plugin.Device{
			ID:     deviceBlockID,
			Health: plugin.Unhealthy,
		}
		d.DeviceControllers[controllerIndex].partsMap[i] = deviceBlockID
	}
	d.DevicePartsAmount += deviceAmount
	//log.Debugf("Sending signal working")
	d.deviceUpdateChan <- true

	//log.Debugf("Sending signal finished")

	return ""
}

func (d *DeviceCategory) AddBlock(accessPoint string, deviceAmount int) string {
	controllerIndex := fmt.Sprintf("%x", md5.Sum([]byte(accessPoint)))
	if _, exist := d.DeviceControllers[controllerIndex]; !exist {
		return fmt.Sprintf("addingBlockToNon-ExistGroup: %s;", accessPoint)
	}
	existedPartsNum := len(d.DeviceControllers[controllerIndex].partsMap)
	for i := 0; i < deviceAmount; i++ {
		deviceBlockID := uuid.NewString()
		d.DeviceParts[deviceBlockID] = &plugin.Device{
			ID:     deviceBlockID,
			Health: plugin.Unhealthy,
		}
		d.DeviceControllers[controllerIndex].partsMap[i+existedPartsNum] = deviceBlockID
	}
	d.DevicePartsAmount += deviceAmount
	d.deviceUpdateChan <- true
	return ""
}

//func (d *DeviceCategory) AddDevice(accessPoint string, deviceAmount int) {
////	log.Debugf("Trying to add device blocks")
//	d.DeviceControllers = append(d.DeviceControllers, &deviceController{
//		accessPoint: accessPoint,
//		BrokenParts: list.New(),
//		UsedParts:   list.New(),
//		FreeParts:   list.New(),
//	})
//
//	newDeviceParts := make([]*plugin.Device, deviceAmount)
//	d.DeviceParts = append(d.DeviceParts, newDeviceParts...)
//	for i := 0; i < deviceAmount; i++ {
////		//log.Debugf("Device Parts: %s", d.DeviceParts)
//		d.DeviceParts[i+d.DevicePartsAmount] = &plugin.Device{
//			ID:     uuid.NewString(),
//			Health: plugin.Unhealthy,
//		}
//	}
//	d.DevicePartsAmount += deviceAmount
//}

func (d *DeviceCategory) listAllDevices() []*plugin.Device {
	devicePartsLength := len(d.DeviceParts)
	deltas := make([]*plugin.Device, devicePartsLength)
	index := 0
	for _, v := range d.DeviceParts {
		deltas[index] = v
		index++
	}
	return deltas
}

func (d *DeviceCategory) deviceStatusResponder() {
	for {
		//log.Debugf("Waiting for query")
		queryInfo := d.deviceMsgSrv.Serve().(string) // queryInfo = (string)"getAll"
		log.Debugf("Received: %s", queryInfo)
		if queryInfo == "getAllDevices" {
			//log.Debugf("GetAllDevicesCalled")
			d.deviceMsgSrv.Respond(d.listAllDevices())
			log.Debugf("GetAllDeviceResponded")
		} else if queryInfo == "invokeMonitor" {
			go d.monitor()
			//log.Debugf("Monitor invoked!!!")
		} else {
			d.deviceMsgSrv.Respond("Request Failed")
		}
	}
}

func (d *DeviceCategory) FreeDevice() {

}

func (d *DeviceCategory) AllocateDevice() {

}

func (d *DeviceCategory) monitor() {
	log.Debugf("Monitor Invoked")
	go func() {
		router := utils.GetMessageRouter()
		time.Sleep(10 * time.Second)
		for {
			_ = <-d.deviceUpdateChan
			//log.Debugf("Monitor received changes signal")
			deltas := d.listAllDevices()
			log.Debugf("Trying to call devicePlugin to update............")
			//router.Call("devicePlugin", d.deviceCategoryID, deltas)

			respond := router.Call("devicePlugin", d.deviceCategoryID, deltas).(string)
			log.Debugf("call devicePlugin to update finished, received respond: %s", respond)
		}
	}()
	for {
		//log.Debugf("Current block num total:%d", len(d.DeviceParts))
		time.Sleep(10 * time.Second)

		for _, group := range d.DeviceControllers {
			log.Debugf("Group: %s, groupblocks: %d", group.accessPoint, len(group.partsMap))
		}
		d.deviceUpdateChan <- true

	}

}
