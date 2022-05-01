package iote

type Configuration struct {
	Broker		string		`json:"mqtt-broker"`

	DebugMQTT	bool		`json:"mqtt-debug"`
	Verbose		bool		`json:"verbose"`
}

var (
	config Configuration
)

func init() {
	config.Broker =  "tcp://localhost:1883"
}

func GetConfig() (Configuration) {
	return config
}
