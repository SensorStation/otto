package mocks

import (
	"testing"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestMockRoutes(t *testing.T) {
	root := newNode("")
	if root.nodes == nil {
		t.Error("root node is nil")
	}
	if len(root.nodes) != 0 {
		t.Errorf("root node expected (0) got (%d)", len(root.nodes))
	}
	if len(root.handlers) != 0 {
		t.Errorf("root handlers expected (0) got (%d)", len(root.handlers))
	}

	gotit := false
	topic := "test/path/full"
	c := MockClient{}
	m := MockMessage{topic: topic}
	root.insert(topic, func(c gomqtt.Client, m gomqtt.Message) {
		gotit = true
	})

	n := root.lookup(topic)
	if n == nil {
		t.Errorf("expect node for %s got nil", topic)
	}

	if n.handlers == nil {
		t.Errorf("expected node handlers (1) got (%d)", len(n.handlers))
	}

	n.pub(c, m)
	if !gotit {
		t.Errorf("message failed to get published")
	}
}
