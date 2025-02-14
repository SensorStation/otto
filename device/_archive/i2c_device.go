package devices

type I2CDevice struct {
	*BaseDevice
	Bus  string
	Addr int
}

func NewI2CDevice(name string, bus string, addr int) *I2CDevice {
	return &I2CDevice{
		BaseDevice: NewDevice(name),
		Bus:        bus,
		Addr:       addr,
	}
}
