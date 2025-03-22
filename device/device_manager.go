package device

type DeviceManager struct {
	devices map[string]Name `json:"devices"`
}

var (
	stationName string = "station"
	devices     *DeviceManager
)

func init() {
}

func GetDeviceManager() *DeviceManager {
	if devices == nil {
		devices = &DeviceManager{
			devices: make(map[string]Name),
		}
	}
	return devices
}

func (dm *DeviceManager) Add(d Name) {
	if dm.devices == nil {
		dm.devices = make(map[string]Name)
	}
	dm.devices[d.Name()] = d
}

func (dm *DeviceManager) Get(name string) (Name, bool) {
	d, ex := dm.devices[name]
	return d, ex
}

// func (dm *DeviceManager) FindPin(offset int) Device {
// 	for _, d := range dm.devices {
// 		switch d.(type) {
// 		case
// 		}
// 		if d.Offset() == offset {
// 			return d
// 		}
// 	}
// 	return nil
// }
