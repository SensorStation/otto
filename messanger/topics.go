package messanger

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sensorstation/otto/station"
)

type TopicList map[string]int

var (
	Topics    TopicList
	TopicBase string
)

func init() {
	TopicBase = "ss/%s/%s/%s"
	Topics = make(map[string]int)
}

func TopicControl(topic string) string {
	t := fmt.Sprintf(TopicBase, "c", station.StationName, topic)
	Topics[t]++
	return t
}

func TopicData(topic string) string {
	t := fmt.Sprintf(TopicBase, "d", station.StationName, topic)
	Topics[t]++
	return t
}

func (t TopicList) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jstr, err := json.Marshal(t)
	if err != nil {
		http.Error(w, "Not Yet Supported", 401)
	}
	fmt.Fprint(w, jstr)
}
