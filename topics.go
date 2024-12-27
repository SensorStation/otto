package otto

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Topics map[string]int

var (
	topics    Topics
	topicBase string
)

func init() {
	topicBase = "ss/%s/%s/%s"
	topics = make(map[string]int)
}

func TopicControl(topic string) string {
	t := fmt.Sprintf(topicBase, "c", StationName, topic)
	topics[t]++
	return t
}

func TopicData(topic string) string {
	t := fmt.Sprintf(topicBase, "d", StationName, topic)
	topics[t]++
	return t
}

func (t Topics) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jstr, err := json.Marshal(t)
	if err != nil {
		http.Error(w, "Not Yet Supported", 401)
	}
	fmt.Fprint(w, jstr)
}
