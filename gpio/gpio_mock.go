package gpio

import (
	"fmt"
	"time"

	"github.com/sensorstation/otto"
	"github.com/warthog618/go-gpiocdev"
)

// MockGPIO fakes the Line interface on computers that don't
// actually have GPIO pins mostly for mocking tests
type MockLine struct {
	offset                int `json:"offset"`
	Val                   int `json:"val"`
	gpiocdev.EventHandler `json:"event-handler"`
	start                 time.Time
}

func GetMockLine(offset int, opts ...gpiocdev.LineReqOption) *MockLine {
	l := otto.GetLogger()
	m := &MockLine{
		offset: offset,
		start:  time.Now(),
	}
	for _, opt := range opts {
		fmt.Printf("OPT Type: %T\n", opt)
		switch v := opt.(type) {
		case gpiocdev.OutputOption:
			m.Val = v[0]
			fmt.Printf("OO: %+v\n", v[0])

		case gpiocdev.EventHandler:
			m.EventHandler = opt.(gpiocdev.EventHandler)

		default:
			l.Info("MockLine does not record", "optType", v)
		}
	}

	return m
}

func (m MockLine) Close() error {
	return nil
}

func (m MockLine) Offset() int {
	return m.offset
}

func (m MockLine) SetValue(val int) error {
	m.Val = val
	return nil
}

func (m MockLine) Reconfigure(...gpiocdev.LineConfigOption) error {
	println("We have entered reconfigure")
	return nil
}

func (m MockLine) Value() (int, error) {
	return m.Val, nil
}

var seqno uint32

func getSeqno() uint32 {
	seqno += 1
	return seqno
}

func (m MockLine) MockHWInput(v int) {
	m.Val = v

	t := gpiocdev.LineEventRisingEdge
	if v == 0 {
		t = gpiocdev.LineEventFallingEdge
	}

	seq := getSeqno()
	if m.EventHandler != nil {
		evt := gpiocdev.LineEvent{
			Offset:    m.Offset(),
			Timestamp: time.Since(m.start),
			Type:      t,
			Seqno:     seq,
			LineSeqno: seq,
		}

		m.EventHandler(evt)
	}

	mqtt, err := otto.GetMQTT()
	if err != nil {
		otto.GetLogger().Error("Failed to connect mqtt", "error", err)
	}
	topic := fmt.Sprintf("ss/c/+/%d", m.offset)
	mqtt.Subscribe(topic, m)
}

func (m MockLine) SubCallback(msg *otto.Msg) {
	l := otto.GetLogger()

	// Change this to a map[string]string or map[string]interface{}
	fmt.Printf("MSG: %+v\n", msg)
	str := msg.String()
	switch str {
	case "on":
		m.MockHWInput(1)

	case "off":
		m.MockHWInput(0)

	default:
		l.Warn("bad hw mock value", "value", str)

	}
}
