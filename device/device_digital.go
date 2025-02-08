package device

// import "github.com/warthog618/go-gpiocdev"

// func NewDigitalDevice(name string, offset int, opts ...gpiocdev.LineReqOption) *Device {
// 	d := NewDevice(name)
// 	g := GetGPIO()
// 	p := g.Pin(name, offset, opts...)
// 	d.ReadWriteCloser = p
// 	return d
// }

// type DigitalMock struct {
// }

// func GetDigitalMock(name string) *DigitalMock {
// 	return &DigitalMock{}
// }

// func (d *DigitalMock) Get() (int, error) {
// 	return 1, nil
// }

// func (d *DigitalMock) Set(val int) error {
// 	return nil
// }

// func (d *DigitalMock) On() error {
// 	return d.Set(1)
// }

// func (d *DigitalMock) Off() error {
// 	return d.Set(0)
// }

// func (d *DigitalMock) Close() error {
// 	return nil
// }

// func (d *DigitalMock) Read([]byte) (n int, err error) {
// 	return n, err
// }

// func (d *DigitalMock) Write([]byte) (n int, err error) {
// 	return n, err
// }
