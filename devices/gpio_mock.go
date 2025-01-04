package devices

import (
	"strconv"
	"time"

	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/message"
	"github.com/sensorstation/otto/messanger"
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

func GetMockGPIO() *GPIO {
	gpio := GetGPIO()
	gpio.Mock = true
	return gpio
}

func GetMockLine(offset int, opts ...gpiocdev.LineReqOption) *MockLine {
	l := logger.GetLogger()
	m := &MockLine{
		offset: offset,
		start:  time.Now(),
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case gpiocdev.OutputOption:
			m.Val = v[0]

		case gpiocdev.EventHandler:
			m.EventHandler = opt.(gpiocdev.EventHandler)
			mqtt := messanger.GetMQTT()
			mqtt.Subscribe(messanger.TopicControl("mock/"+strconv.Itoa(m.Offset())), m.Callback)

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

func (m *MockLine) SetValue(val int) error {
	m.Val = val
	return nil
}

func (m MockLine) Reconfigure(...gpiocdev.LineConfigOption) error {
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

func (m *MockLine) MockHWInput(v int) {
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
}

func (m *MockLine) Callback(msg *message.Msg) {
	str := msg.String()
	switch str {
	case "on":
		fallthrough

	case "1":
		m.MockHWInput(1)

	case "off":
		fallthrough

	case "0":
		m.MockHWInput(0)

	default:
		return
	}
	return
}
