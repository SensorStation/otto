package blasters

import (
	"fmt"
	"testing"
	"time"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/mocks"
)

func TestBlasters(t *testing.T) {
	blasters := NewMQTTBlasters(5)
	if blasters.Count != 5 {
		t.Errorf("expected (5) blasters got (%d)", blasters.Count)
	}

	if blasters.Running {
		t.Error("expected running to be (false) got (true)")
	}

	if blasters.Wait != 2000 {
		t.Errorf("expected wait to be (2000) got (%d)", blasters.Wait)
	}

	for i := 0; i < 5; i++ {
		b := blasters.Blasters[i]
		stid := fmt.Sprintf("station-%d", i)
		if b.ID != stid {
			t.Errorf("expected station id (%s) got (%s)", stid, b.ID)
		}

		topic := fmt.Sprintf("ss/d/%s/temphum", stid)
		if b.Topic != topic {
			t.Errorf("expected topic (%s) got (%s)", topic, b.Topic)
		}
	}
}

func TestBlasting(t *testing.T) {
	c := mocks.GetMockClient()
	m := otto.GetMQTTClient(c)
	m.Connect()

	blasters := NewMQTTBlasters(5)
	for _, bl := range blasters.Blasters {
		m.Subscribe(bl.Topic, bl.Station.DataManager)
	}

	go blasters.Blast()
	time.Sleep(2 * time.Second)

	for _, bl := range blasters.Blasters {
		fmt.Printf("bl: %+v\n", bl)
		fmt.Printf("\tst: %+v\n", bl.Station.DataManager)
	}

	blasters.Stop()
}
