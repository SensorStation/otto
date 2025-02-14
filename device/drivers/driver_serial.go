package device

func NewSerialDevice(name string, port string, opts any) *Device {
	d := NewDevice(name)
	sp := GetSerial(name, port, opts)
	err := sp.Open()
	d.Error = err
	return d
}

// type SerialDevice interface {
// 	Device
// 	Open() error
// 	PortName() string
// 	SerialReader
// 	SerialWriter
// }

// type SerialWriter interface {
// 	io.Reader
// }

// type SerialReader interface {
// 	io.Writer
// }

// type Closer interface {
// 	Close() error
// }

// func NewSerialDevice(name string, port string, opts any) SerialDevice {
// 	if mock {
// 		return GetSerialMock(name)
// 	}
// 	return GetSerialPort(name, port, opts)
// }

// type SerialMock struct {
// 	*BaseDevice
// 	portName string
// }

// func GetSerialMock(name string) *SerialMock {
// 	return &SerialMock{
// 		BaseDevice: NewDevice(name),
// 	}
// }

// func (d *SerialMock) Open() error {
// 	return nil
// }

// func (d *SerialMock) Read([]byte) (n int, err error) {
// 	return n, err
// }

// func (d *SerialMock) Write([]byte) (n int, err error) {
// 	return n, err
// }

// func (d *SerialMock) Close() error {
// 	return nil
// }

// func (d *SerialMock) PortName() string {
// 	return d.portName
// }
