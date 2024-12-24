package main

type DeviceManager struct {
	devices map[string]Device
}

var (
	stationName string = "station"
	devices     DeviceManager
)

func init() {
	devices.devices = make(map[string]Device)
}

func (dm *DeviceManager) Add(d Device) {
	if dm.devices == nil {
		dm.devices = make(map[string]Device)
	}
	dm.devices[d.Name()] = d
}

func (dm *DeviceManager) Get(name string) Device {
	d, ex := dm.devices[name]
	if !ex {
		l.Error("device does not exist", "device", name)
		return nil
	}
	return d
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
