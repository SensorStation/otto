package device

type DeviceManager struct {
	Devices map[string]Name `json:"-"`
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
			Devices: make(map[string]Name),
		}
	}
	return devices
}

func (dm *DeviceManager) Add(d Name) {
	if dm.Devices == nil {
		dm.Devices = make(map[string]Name)
	}
	dm.Devices[d.Name()] = d
}

func (dm *DeviceManager) Get(name string) (Name, bool) {
	d, ex := dm.Devices[name]
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
