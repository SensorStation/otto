package iote

import (
	"fmt"
	"log"
	"strconv"
)

type Msg struct {
	Station string		`json:"station"`
	Sensor	string		`json:"sensor"`
	Time	int64		`json:"time"`
	Data	[]byte		`json:"-"`
}

type MsgFloat64 struct {
	Msg
	Value	float64		`json:"value"`
}

func (m Msg) ToMsgFloat64() (m64 MsgFloat64) {
	m64 = MsgFloat64{
		Msg: m,
	}
	m64.Value = m.Float64()
	return m64
}

func (m *Msg) Float64() float64 {
	val, err := strconv.ParseFloat(string(m.Data), 64)
	if err != nil {
	 	log.Println(err)
		return 0.0
	}
	return val
}

func (m *Msg) String() (dstr string) {
	str := fmt.Sprintf("%q - %s - %x - %s",
		dstr, m.Station, m.Sensor, m.Time, string(m.Data))
	return str
}
