package mesh

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	// mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MeshNetwork represents a network of devices that have meshed up.
type MeshNetwork struct {
	ID         string `json:"id"`
	Pass       string `json:"passwd"`
	MeshRouter `json:"router"`

	RootId string               `json:"rootid"`
	Nodes  map[string]*MeshNode `json:"nodes"`
}

var (
	mn MeshNetwork
)

// GetNode will return the node associated with the given ID
func (m *MeshNetwork) GetNode(nid string) (mn *MeshNode) {
	var e bool

	if mn, e = m.Nodes[nid]; !e {
		mn = &MeshNode{Id: nid}
		m.Nodes[nid] = mn
	}
	return mn
}

// UpdateRoot will reroot the mesh network with the new root id
func (m *MeshNetwork) UpdateRoot(rootid string) {

	// TODO create a fully configured node and schedule network topology updates.
	// l.Printf("%s.%s %s[%.0f]: rootid: %s, self: %s, parent: %s\n",
	//	addr, typ, msgtype, layer, rootid, self, parent);
	if m.RootId != rootid {
		// we have a change of roots
		slog.Info("Root Node has changed", "from", m.RootId, "to", rootid)
		m.RootId = rootid
	}
}

// ServeHTTP provides a REST interface to the config structure
func (m MeshNetwork) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(m)

	case "POST", "PUT":
		// TODO
		http.Error(w, "Not Yet Supported", 401)
	}
}

// MeshRouter is the optional IP router for the mesh network
type MeshRouter struct {
	SSID string `json:"ssid"`
	Pass string `json:"passwd"`
	Host string `json:"host"`
}

// MeshNode represents a single node in the ESP-MESH network, this allows us
// to keep track of our inventory fleet.
type MeshNode struct {
	Id       string            `json:"id"`
	Parent   string            `json:"parend"`
	Layer    int               `json:"layer"`
	Children map[string]string `json:"children"`
	Updated  time.Time         `json:"time"`
}

// NewNode will create a new network node
func NewNode(d map[string]interface{}) *MeshNode {
	self := d["self"].(string)
	parent := d["parent"].(string)
	pnode := mn.GetNode(parent)

	mn := &MeshNode{
		Id:      self,
		Parent:  pnode.Id,
		Layer:   int(d["layer"].(float64)),
		Updated: time.Now(),
	}
	return mn
}

// UpdateParent will reestablish the the given parent node
func (n *MeshNode) UpdateParent(p *MeshNode) {
	if n.Parent != p.Id {
		slog.Info("n.Parent has changed", "from", n.Parent, "to", p.Id)
	}
	n.Parent = p.Id
}

// UpdateChild will re-establish the given child node
func (n *MeshNode) UpdateChild(c *MeshNode) {
	slog.Info("Update child", "Parent ", n.Id)
	if n.Children == nil {
		n.Children = make(map[string]string)
	}

	if _, e := n.Children[c.Id]; e {
		slog.Info("update existing child ")
	} else {
		slog.Info(" ADDING NEW child ")
	}
	n.Children[c.Id] = c.Id
	slog.Info(c.Id)
}

// String representation of this meshnode
func (n *MeshNode) String() string {
	str := fmt.Sprintf("NODE self - %s :=: parent - %s :=: layer - %d last update: %q\n",
		n.Id, n.Parent, n.Layer, n.Updated)
	if len(n.Children) < 1 {
		return str
	}
	str += "Chilren:\n"
	for _, mn := range n.Children {
		str += "\t" + mn + "\n"
	}
	return str
}

// MeshMessage is what is passed around amoung mesh nodes
type MeshMessage struct {
	Addr string `json:"addr"`
	Typ  string `json:"type"`
	Data []byte `json:"data"`
}

// MeshHeartBeat is a periodic message that advertises the liveness of the
// given mesh node
type MeshHeartBeat struct {
	Typ    string `json:"type"`   // heartbeat
	Id     string `json:"self"`   // macaddr of advertising node
	Parent string `json:"parent"` // macaddr of parent
	Layer  int    `json:"layer"`  // node layer
}

// MsgRecv reads a meshnode message and handles the payload
func (mn MeshNetwork) MsgRecv(topic string, payload []byte) {
	/*
		var m MeshMessage
		err := json.Unmarshal(payload, &m)
		if err != nil {
			slog.Error("Failed to unmarshal payload")
			return
		}

		// unravel the json message and verify our current node information
		paths := strings.Split(topic, "/")
		if len(paths) != 3 {
			slog.Error("Error unsupported path len", "pathlen", len(paths))
			return
		}

		rootid := paths[1]
		data := m.Data
		msgtype := m.data["type"]

		switch msgtype {
		case "heartbeat":

			self, _ := data["self"].(string)
			parent, _ := data["parent"].(string)
			layer, _ := data["layer"].(int)
			mn.Update(rootid, self, parent, layer)

		default:
			slog.Fatalln("Unknown message type: ", msgtype)
		}
	*/
	return
}

// Update a mesh network with the given information
func (mn MeshNetwork) Update(rootid, id, parent string, layer int) {

	slog.Debug("[MESH] Update [id/parent/rootid/layer]: ", id, parent, rootid, layer)
	if mn.RootId != rootid {
		slog.Info("[MESH] Root %s has changed to %s\n", mn.RootId, rootid)
		mn.RootId = rootid
	}

	node := mn.GetNode(id)
	if node == nil || node.Id == "" {
		node.Parent = parent
	}

	if node.Layer != layer {
		slog.Info("[MESH] Node %s layer has changed from %d to %d\n", node.Id, node.Layer, "layer", layer)
	}
}
