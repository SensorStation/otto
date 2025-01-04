package utils

var (
	stationName string
)

func init() {
	stationName = "station"
}

func SetStation(name string) {
	stationName = name
}

func Station() string {
	return stationName
}
