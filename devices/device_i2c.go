package devices

type I2CDevice struct {
	*Device
	Bus  string
	Addr int
}

func NewI2CDevice(name string, bus string, addr int) I2CDevice {
	return I2CDevice{
		Device: &Device{
			Name: name,
		},
		Bus:  bus,
		Addr: addr,
	}
}
