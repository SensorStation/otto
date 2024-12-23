package otto

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sensorstation/otto/message"
)

func getMsg() (*message.Msg, time.Time) {
	now := time.Now()
	path := "ss/d/%s/test"

	b := fmt.Sprintf("%d", 4)
	m := message.NewMsg(path, []byte(b), "test")
	m.Source = "be:ef:ca:fe:01"
	m.Time = now

	return m, now
}

func TestStationMsg(t *testing.T) {
	topic := "ss/d/be:ef:ca:fe:01/station"
	omsg, _ := getMsg()
	omsg.Source = "be:ef:ca:fe:01"

	j, err := json.Marshal(omsg)
	if err != nil {
		t.Errorf("json marshal failed %+v", err)
		return
	}

	msg := message.NewMsg(topic, j, "test")
	if msg == nil {
		t.Error("msg topic expected but is nil")
	}

	if msg.Topic != topic {
		t.Errorf("msg topic expected (%s) got (%s)", topic, msg.Topic)
	}

	path := strings.Split(topic, "/")
	if len(path) != len(msg.Path) {
		t.Errorf("msg path len expected (%d) got (%d)", len(path), len(msg.Path))
	}

	for i := 0; i < len(path); i++ {
		if path[i] != msg.Path[i] {
			t.Errorf("msg path[%d] expected (%s) got (%s)", i, path[i], msg.Path[i])
		}
	}

	if msg.Source != "test" {
		t.Errorf("msg source expected (test) got (%s)", msg.Source)
	}

	for i := 0; i < len(j); i++ {
		if msg.Data[i] != j[i] {
			t.Errorf("msg data[%d] expected (% x) got (% x)", i, j[i], msg.Data[i])
		}
	}
}

func TestJSON(t *testing.T) {
	m, _ := getMsg()

	_, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Failed to marshal message %+v", m)
	}
}
