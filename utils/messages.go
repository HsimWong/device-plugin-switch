package utils

type DeviceStatus struct {
	deviceID string
	status   string
}

type DeviceRegisterRequest struct {
	deviceCategoryType string
	deviceBlockNum     int
	accessPoint        string
}

type DeviceRegisterResponse struct {
	deviceRegisterID      string
	deviceRegisterResults string
}

type DeviceStatusDelta struct {
}
