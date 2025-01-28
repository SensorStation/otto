package devices

import "time"

type AnalogMock struct {
    BaseDevice
}

func (a *AnalogMock) Get() float64 {
    return 54.5
}

func (a *AnalogMock) ReadContinuous() <-chan float64 {
    q := make(chan float64)

    go func() {
        for {
            v := a.Get()
            q <-v
            time.Sleep(1 * time.Second)
        }
    }()

    return q
}
