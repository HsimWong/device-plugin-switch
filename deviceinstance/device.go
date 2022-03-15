package deviceinstance

import (
	"container/list"
	"github.com/HsimWong/device-plugin-switch/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"time"
)

type deviceController struct {
	accessPoint string
	Anchor      int        // Starts from 0
	BrokenParts *list.List // index starts from 0
	UsedParts   *list.List
	FreeParts   *list.List
}

type DeviceCategory struct {
	deviceCategoryID    string
	DevicePartsAmount   int
	ControllerIndexBars []int
	DeviceParts         []*plugin.Device
	DeviceControllers   []*deviceController
	deviceMsgSrv        *utils.SyncMessenger
}

func NewDevice(deviceCategoryID string,
	deviceMsgSrv *utils.SyncMessenger) *DeviceCategory {
	return &DeviceCategory{
		deviceCategoryID:    deviceCategoryID,
		DevicePartsAmount:   0,
		ControllerIndexBars: nil,
		DeviceParts:         nil,
		DeviceControllers:   nil,
		deviceMsgSrv:        deviceMsgSrv,
	}
}

func (d *DeviceCategory) Run() {
	log.Debugf("Device start running")
	go d.monitor()
	go d.deviceStatusResponder()
}

func (d *DeviceCategory) AddDevice(accessPoint string, deviceAmount int) {
	log.Debugf("Trying to add device blocks")
	d.ControllerIndexBars = append(d.ControllerIndexBars,
		d.DevicePartsAmount+deviceAmount)
	d.DeviceControllers = append(d.DeviceControllers, &deviceController{
		accessPoint: accessPoint,
		Anchor:      d.DevicePartsAmount,
		BrokenParts: list.New(),
		UsedParts:   list.New(),
		FreeParts:   list.New(),
	})

	newDeviceParts := make([]*plugin.Device, deviceAmount)
	d.DeviceParts = append(d.DeviceParts, newDeviceParts...)
	for i := 0; i < deviceAmount; i++ {
		//log.Debugf("Device Parts: %s", d.DeviceParts)
		d.DeviceParts[i+d.DevicePartsAmount] = &plugin.Device{
			ID:     uuid.NewString(),
			Health: plugin.Unhealthy,
		}
	}
	d.DevicePartsAmount += deviceAmount
}

func (d *DeviceCategory) deviceStatusResponder() {
	for {
		log.Debugf("Waiting for query")
		queryInfo := d.deviceMsgSrv.Serve().(string) // queryInfo = (string)"getAll"
		if queryInfo == "getAllDevices" {
			d.deviceMsgSrv.Respond(d.DeviceParts)
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
	for {
		log.Debugf("Current block num:%d", len(d.DeviceParts))
		time.Sleep(10 * time.Second)

	}

}
