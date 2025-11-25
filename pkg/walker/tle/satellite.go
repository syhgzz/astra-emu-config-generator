package tle

import (
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/joshuaferrara/go-satellite"
	gosate "github.com/joshuaferrara/go-satellite"
)

type Satellite struct {
	Name      string
	Satellite gosate.Satellite
}

func NewStatellite(name, line1, line2 string) (*Satellite, error) {
	if name == "" {
		return nil, errors.New("the name of satellite is empty...")
	}
	if line1 == "" {
		return nil, errors.New("the first line of satellite's TLE is empty")
	}
	if line2 == "" {
		return nil, errors.New("the second line of satellite's TLE is empty")
	}
	return &Satellite{
		Name:      name,
		Satellite: gosate.TLEToSat(line1, line2, satellite.GravityWGS72),
	}, nil
}

func NewConstellation(filePath string) ([]Satellite, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", filePath)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	constellation := []Satellite{}
	for idx := 0; idx < len(lines); idx += 3 {
		name := strings.ReplaceAll(strings.Trim(lines[idx], " "), " ", "-")
		// name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "(", "")
		name = strings.ReplaceAll(name, ")", "")
		name = strings.ToLower(name)
		sat, err := NewStatellite(name, lines[idx+1], lines[idx+2])
		if err != nil {
			return nil, errors.Wrap(err, "failed to create constellation")
		}
		constellation = append(constellation, *sat)
	}

	return constellation, nil
}

// Return the location of satellite expressed by longitude/latitude/altitude at given time
func (sat *Satellite) LocationAtTime(time time.Time) (long, lat, alt float64) {
	year, month, day, hour, minute, second :=
		time.Year(),
		int(time.Month()),
		time.Day(),
		time.Hour(),
		time.Minute(),
		time.Second()
	// Get current position
	position, _ := gosate.Propagate(
		sat.Satellite, year, month,
		day, hour, minute, second,
	)
	// Convert the current time to Galileo system time (GST)
	gst := gosate.GSTimeFromDate(
		year, month, day,
		hour, minute, second,
	)
	// Get satellite's latitude & longtitude & altitude
	altitude, _, latlong := gosate.ECIToLLA(position, gst)
	latlong = gosate.LatLongDeg(latlong)
	// log.Printf("%s: long is %2f, lat is %2f, alt is %2f\n",
	// 	sat.Name, latlong.Longitude, latlong.Latitude, altitude)

	return latlong.Longitude, latlong.Latitude, altitude * 1000
}
