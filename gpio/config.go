package gpio

type Configuration struct {
	Mock bool
}

var config Configuration

func init() {
	config.Mock = true
}
