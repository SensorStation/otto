package mocks

import (
	"strings"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

type node struct {
	nodes    map[string]*node
	handlers []gomqtt.MessageHandler
}

var (
	root node
)

func init() {
	root = node{
		nodes: make(map[string]*node),
	}
}

func newNode(index string) *node {
	n := &node{
		nodes: make(map[string]*node),
	}
	return n
}

func (n *node) insert(topic string, mh gomqtt.MessageHandler) {
	indexes := strings.Split(topic, "/")
	pn := n
	for _, idx := range indexes {
		if nn, ex := pn.nodes[idx]; !ex {
			nn = newNode(idx)
			pn.nodes[idx] = nn
			pn = nn
		}
	}
	// The last node push the callback on the callback list
	pn.handlers = append(pn.handlers, mh)
}

func (n *node) lookup(topic string) *node {
	indexes := strings.Split(topic, "/")
	pn := n
	for _, idx := range indexes {
		nn, ex := pn.nodes["#"]
		if ex {
			// found the wildcard. XXX this will be a bug if it is not
			// the last node on the  topic
			return nn
		}

		nn, ex = pn.nodes["+"]
		if ex {
			// we will accept this path portion of the wildcard, but
			// must continue on
			pn = nn
			continue
		}

		nn, ex = pn.nodes[idx]
		if !ex {
			return nil
		}
		pn = nn
	}
	return pn
}

func (n *node) pub(c gomqtt.Client, m gomqtt.Message) {
	for _, h := range n.handlers {
		h(c, m)
	}
}
