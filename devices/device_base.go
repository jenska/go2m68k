package devices

type DeviceBase struct {
	name       string
	owner      *DeviceBase
	subDevices []*DeviceBase
}

func (device *DeviceBase) String() string {
	return device.name
}

type DeviceInterface interface {
}
