package otto

import "fmt"

var (
	topics    map[string]int
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
	t := fmt.Sprintf(topicBase, "c", StationName, topic)
	topics[t]++
	return t
}
