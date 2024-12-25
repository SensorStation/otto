package devices

type I2CDevice struct {
	*Dev
	Bus  string
	Addr int
}

func NewI2CDevice(name string, bus string, addr int) I2CDevice {
	return I2CDevice{
		Dev: &Dev{
			name: name,
		},
		Bus:  bus,
		Addr: addr,
	}
}
