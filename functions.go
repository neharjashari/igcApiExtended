package main

import (
	"fmt"
	"github.com/marni/goigc"
	"strconv"
	"time"
)


// FormatTimeSince function formats time with iso8601 json format
// i found this function online because i couldn't find a better way of formating with json format
func FormatTimeSince(t time.Time) string {

	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
	)

	timeSince := time.Since(t)

	sign := time.Duration(1)

	if timeSince < 0 {
		sign = -1
		timeSince = -timeSince
	}

	timeSince += +Decisecond / 2

	days := sign * (timeSince / Day)
	timeSince = timeSince % Day

	hours := timeSince / time.Hour
	timeSince = timeSince % time.Hour

	minutes := timeSince / time.Minute
	timeSince = timeSince % time.Minute

	seconds := timeSince / time.Second
	timeSince = timeSince % time.Second

	f := timeSince / Decisecond

	years := days / 365

	return fmt.Sprintf("P%dY%dD%dH%dM%d.%dS", years, days, hours, minutes, seconds, f)
}

// trackLength function calculates the track length of each track based on the track Points.
func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}

// FloatToString function returns a float number formated as a string
func FloatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 4, 64)
}
