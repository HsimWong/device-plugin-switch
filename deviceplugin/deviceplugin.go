package deviceplugin

import (
	"context"
	"github.com/HsimWong/device-plugin-switch/deviceinstance"
	"github.com/HsimWong/device-plugin-switch/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"net"
	"os"
	"path"
	"time"
)

type Instance struct {
	deviceType       string
	deviceCategoryID string
	lis              net.Listener
	srv              *grpc.Server
	ctx              context.Context
	cancel           context.CancelFunc
	msgRcv           *utils.SyncMessenger
	devices          *deviceinstance.DeviceCategory
	devSocket        string
}

func NewDevicePluginInstance(deviceType string, deviceCategoryID string,
	msgRcv *utils.SyncMessenger, deviceCategory *deviceinstance.DeviceCategory) *Instance {
	return &Instance{
		deviceType:       deviceType,
		deviceCategoryID: deviceCategoryID,
		lis:              nil,
		srv:              nil,
		ctx:              nil,
		cancel:           nil,
		msgRcv:           msgRcv,
		devices:          deviceCategory,
		devSocket:        path.Join(utils.DevicePluginDir, deviceType+".sock"),
	}
}

func (d *Instance) Run() {
	_, err := os.Stat(d.devSocket)
	if err == nil {
		err := os.Remove(d.devSocket)
		utils.Check(err, "Removing Failed")
	}
	// Setting up grpc server for device plugin client
	lis, err := net.Listen("unix",
		path.Join(utils.DevicePluginDir, d.deviceType+".sock"))
	utils.Check(err, "Setting Unix listener failed")
	d.lis = lis
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.srv = grpc.NewServer()

	plugin.RegisterDevicePluginServer(d.srv, d)
	go func() {
		err = d.srv.Serve(lis)
		utils.Check(err, "Error occurs when serving device plugin")
		log.Debugf("Error device plugin: %s", d.deviceType)
	}()
	//time.Sleep(5 * time.Second)
	//conn, err := utils.Dial(dp.devSocketPath, 5 * time.Second)

	// Registering to kubelet
	conn, err := utils.Dial(path.Join(utils.DevicePluginDir,
		utils.KubeletSocket), 5*time.Second)
	//conn, err := grpc.Dial(),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.Check(err, "Error when register to kubelet")

	dpClient := plugin.NewRegistrationClient(conn)
	log.Debugf("Endpoint: %s, ResourceName:%s, deviceType:%s", d.devSocket, d.deviceType, d.deviceType)

	req := &plugin.RegisterRequest{
		Version:      plugin.Version,
		Endpoint:     path.Base(d.devSocket),
		ResourceName: "csu.ac.cn/" + d.deviceType,
	}
	_, err = dpClient.Register(context.Background(), req)
	utils.Check(err, "Registering device plugin failed")
	log.Info("Register Finished")
}

func (d *Instance) GetDevicePluginOptions(ctx context.Context, empty *plugin.Empty) (*plugin.DevicePluginOptions, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Instance) ListAndWatch(empty *plugin.Empty,
	server plugin.DevicePlugin_ListAndWatchServer) error {
	// This function runs in separate thread
	log.Info("ListAndWatch called by kubelet")
	// Initial Reporting to kubelet
	router := utils.GetMessageRouter()
	allDevices := router.Call("device",
		d.deviceCategoryID, "getAllDevices").([]*plugin.Device)
	err := server.Send(&plugin.ListAndWatchResponse{
		Devices: allDevices,
	})
	utils.Check(err, "Reporting List&Watch failed")

	// drop in loop, reporting to kubelet continuously
	for {
		deviceDeltas := d.msgRcv.Serve().([]*plugin.Device)
		err := server.Send(&plugin.ListAndWatchResponse{
			Devices: deviceDeltas,
		})
		utils.Check(err, "Reporting List&Watch device deltas failed")
	}

	return nil
}

func (d *Instance) GetPreferredAllocation(ctx context.Context, request *plugin.PreferredAllocationRequest) (*plugin.PreferredAllocationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Instance) Allocate(ctx context.Context, request *plugin.AllocateRequest) (*plugin.AllocateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Instance) PreStartContainer(ctx context.Context, request *plugin.PreStartContainerRequest) (*plugin.PreStartContainerResponse, error) {
	log.Println("PreStartContainer called")
	return &plugin.PreStartContainerResponse{}, nil
	//panic("implement me")
}
