package kubenet

import "time"

// duration is expriment time, in second
// delta is time interval, in milisecond
func FakeTrajectory(sat Satellite, curTime time.Time, duration int, speedupRatio float64, delta int) []Location {
	var locations []Location

	const Second = 1e9
	duration *= Second

	for ofs := 0; ofs*delta <= duration; ofs++ {
		timeOfs := int64(float64(ofs*delta) * speedupRatio)
		long, lat, alt := sat.LocationAtTime(curTime.Add(time.Duration(timeOfs) * time.Millisecond))
		locations = append(locations, Location{
			Longitude: long,
			Latitude:  lat,
			Altitude:  alt,
			Time:      int64(ofs * delta),
		})
	}

	return locations
}
