package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	// mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MeshNetwork represents a network of devices that have meshed up.
type MeshNetwork struct {
	ID   string
	Pass string
	MeshRouter

	RootId string // Id of the root node
	Nodes  map[string]*MeshNode
}

var mn MeshNetwork

func (m *MeshNetwork) GetNode(nid string) (mn *MeshNode) {
	var e bool

	if mn, e = m.Nodes[nid]; !e {
		mn = &MeshNode{Id: nid}
		m.Nodes[nid] = mn
	}
	return mn
}

func (m *MeshNetwork) UpdateRoot(rootid string) {

	// TODO create a fully configured node and schedule network topology updates.
	// log.Printf("%s.%s %s[%.0f]: rootid: %s, self: %s, parent: %s\n",
	//	addr, typ, msgtype, layer, rootid, self, parent);
	if m.RootId != rootid {
		// we have a change of roots
		log.Println("Root Node has changed from ", m.RootId, " to ", rootid)
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
	SSID string
	Pass string
	Host string
}

// MeshNode represents a single node in the ESP-MESH network, this allows us
// to keep track of our inventory fleet.
type MeshNode struct {
	Id       string
	Parent   string
	Layer    int
	Children map[string]string
	Updated  time.Time
}

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

func (n *MeshNode) UpdateParent(p *MeshNode) {
	if n.Parent != p.Id {
		log.Printf("n.Parent has changed from %s -> %s\n", n.Parent, p.Id)
	}
	n.Parent = p.Id
}

func (n *MeshNode) UpdateChild(c *MeshNode) {
	log.Print("Parent ", n.Id)
	if n.Children == nil {
		n.Children = make(map[string]string)
	}

	if _, e := n.Children[c.Id]; e {
		log.Println(" update existing child ")
	} else {
		log.Println(" ADDING NEW child ")
	}
	n.Children[c.Id] = c.Id
	log.Println(c.Id)
}

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

type MeshMessage struct {
	Addr string `json:"addr"`
	Typ  string `json:"type"`
	Data []byte `json:"data"`
}

type MeshHeartBeat struct {
	Typ    string `json:"type"`   // heartbeat
	Id     string `json:"self"`   // macaddr of advertising node
	Parent string `json:"parent"` // macaddr of parent
	Layer  int    `json:"layer"`  // node layer
}

func (mn MeshNetwork) MsgRecv(topic string, payload []byte) {

	// var m MeshMessage
	// err := json.Unmarshal(payload, &m)
	// if err != nil {
	// 	log.Fatal("Failed to unmarshal payload")
	// }

	// // unravel the json message and verify our current node information
	// paths := strings.Split(topic, "/")
	// if len(paths) != 3 {
	// 	log.Fatal("Error unsupported path")
	// }

	// rootid := paths[1]
	// data := m.Data
	// msgtype := data["type"]

	// switch msgtype {
	// case "heartbeat":

	// 	self, _ := data["self"].(string)
	// 	parent, _ := data["parent"].(string)
	// 	layer, _ := data["layer"].(int)
	// 	mn.Update(rootid, self, parent, layer)

	// default:
	// 	log.Fatalln("Unknown message type: ", msgtype)
	// }
	return
}

func (mn MeshNetwork) Update(rootid, id, parent string, layer int) {

	var debug bool
	if debug {
		log.Println("[MESH] Update [id/parent/rootid/layer]: ", id, parent, rootid, layer)
	}

	if mn.RootId != rootid {
		log.Printf("[MESH] Root %s has changed to %s\n", mn.RootId, rootid)
		mn.RootId = rootid
	}

	node := mn.GetNode(id)
	if node == nil || node.Id == "" {
		node.Parent = parent
	}

	if node.Layer != layer {
		log.Printf("[MESH] Node %s layer has changed from %d to %d\n", node.Id, node.Layer, layer)
	}
}
