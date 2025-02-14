package gtu7

import (
	"errors"
	"strconv"
	"strings"
)

type GPS struct {
	*GPGGA
	*GPGSA
	*GPGSV
	*GPGLL
	*GPRMC
	*GPVTG

	complete bool
}

type GPGGA struct {
	Timestamp         float64
	Latitude          float64
	Latdir            string
	Longitude         float64
	Longdir           string
	Quality           int64
	Satellites        int64
	HDOP              float64
	Altitude          float64
	Altunit           string
	Geodial           float64 // subtract from Altitude to get Height Above Ellipsoid (HAE)
	Geounit           string
	AgeCorrection     float64
	CorrectionStation int64 // can be empty
	BaseStationID     int64 // does not exist on gt-u7
	Checksum          string
}

type GPGSA struct {
	sentence string
}

type GPGSV struct {
	sentences []string
}

type GPGLL struct {
	sentence string
}

type GPRMC struct {
	sentence string
}

type GPVTG struct {
	sentence string
}

func (gps *GPS) Parse(input string) (*GPS, error) {
	data := strings.Split(input, ",")

	// Todo add errors
	switch data[0] {
	case "$GPGGA":
		g := &GPGGA{}
		g.Timestamp, _ = strconv.ParseFloat(data[1], 64)
		g.Latitude, _ = strconv.ParseFloat(data[2], 64)
		g.Latdir = data[3]
		g.Longitude, _ = strconv.ParseFloat(data[4], 64)
		g.Longdir = data[5]
		g.Quality, _ = strconv.ParseInt(data[6], 10, 32)
		g.Satellites, _ = strconv.ParseInt(data[7], 10, 64)
		g.HDOP, _ = strconv.ParseFloat(data[8], 64)
		g.Altitude, _ = strconv.ParseFloat(data[9], 64)
		g.Altunit = data[10]
		g.Geodial, _ = strconv.ParseFloat(data[11], 64)
		g.AgeCorrection, _ = strconv.ParseFloat(data[12], 64)
		g.CorrectionStation, _ = strconv.ParseInt(data[13], 10, 64)
		g.Checksum = data[14]
		gps.GPGGA = g

	case "$GPGSA":
		gps.GPGSA = &GPGSA{
			sentence: input,
		}

	case "$GPGSV":
		if gps.GPGSV == nil {
			gps.GPGSV = &GPGSV{}
		}
		gps.GPGSV.sentences = append(gps.GPGSV.sentences, input)

	case "$GPGLL":
		gps.GPGLL = &GPGLL{
			sentence: input,
		}

	case "$GPRMC":
		gps.GPRMC = &GPRMC{
			sentence: input,
		}

	case "$GPVTG":
		gps.GPVTG = &GPVTG{
			sentence: input,
		}

	default:
		return nil, errors.New("Unknown command: " + data[0])
	}

	return gps, nil
}

func (g *GPS) IsComplete() bool {
	if !g.complete {
		if g.GPGGA != nil &&
			g.GPGSA != nil &&
			g.GPGSV != nil &&
			g.GPGLL != nil &&
			g.GPRMC != nil &&
			g.GPVTG != nil {
			g.complete = true
		}
	}
	return g.complete
}
