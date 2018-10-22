package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// *** API TICKER ***


// Gometalinter
const (
	gmlOB  = `{`
	gmlCB  = `}`
	gmlCPC = `",`
)


func getApiTickerLatest(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet { // The request has to be of GET type

		timestamps := tickerTimestamps("")
		latestTimestamp := timestamps.latestTimestamp

		if latestTimestamp.IsZero() { // If you dont assign a time to a time.Time variable, it's value is 0 date. We can check with IsZero() function
			fmt.Fprintln(w, "There are no track records")
		} else { //If it's not zero, we can format and display it to the user
			fmt.Fprintln(w, latestTimestamp.Format("02.01.2006 15:04:05.000"))
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}

}


func getApiTicker(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet { // The request has to be of GET type

		processStart := time.Now() // Track when the process started

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		timestamps := tickerTimestamps("")

		oldestTS := timestamps.oldestTimestamp
		latestTS := timestamps.latestTimestamp

		// timestamps := returnTimestamps(5)

		response := gmlOB
		response += `"t_latest": "`
		if latestTS.IsZero() {
			response += gmlCPC
		} else {
			response += latestTS.Format("02.01.2006 15:04:05.000") + `",`
		}

		response += `"t_start": "`
		if oldestTS.IsZero() {
			response += gmlCPC
		} else {
			response += oldestTS.Format("02.01.2006 15:04:05.000") + `",`
		}

		// returnTracks returns the last element and the n number of tracks
		trackArray, tStop := returnTracks(5)

		// t_stop SHOULD BE ADDED HERE
		response += `"t_stop": "` + tStop.Format("02.01.2006 15:04:05.000") + `",`

		response += `"tracks":` + `[`

		// THAT 5 SHOULD BE ABLE TO CHANGE DYNAMICALLY
		response += trackArray // Maximum of 5 tracks

		response += `],`
		response += `"processing":` + `"` + strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 2, 64) + `ms"`
		response += gmlCB
		fmt.Fprintln(w, response)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}


func getApiTickerTimestamp(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet { // The request has to be of GET type

		processStart := time.Now() // Track when the process started

		pathArray := strings.Split(r.URL.Path, "/") // split the URL Path into chunks, whenever there's a "/"
		timestamp := pathArray[len(pathArray)-1]    // The part after the last "/", is the timestamp

		_, err := time.Parse("02.01.2006 15:04:05.000", timestamp) // Check if the timestamp provided is a valid time

		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // If there is an error, then return a bad request error
			return
		}

		timestamps := tickerTimestamps(timestamp)

		olderTS := timestamps.oldestNewerTimestamp
		latestTS := timestamps.latestTimestamp

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		response := gmlOB
		response += `"t_latest": "`
		if latestTS.IsZero() {
			response += gmlCPC
		} else {
			response += latestTS.Format("02.01.2006 15:04:05.000") + `",`
		}

		response += `"t_start": "`
		if olderTS.IsZero() {
			response += gmlCPC
		} else {
			response += olderTS.Format("02.01.2006 15:04:05.000") + `",`
		}

		// returnTracks returns the last element and the n number of tracks
		trackArray, tStop := returnTracks(5)

		// t_stop SHOULD BE ADDED HERE
		response += `"t_stop": "` + tStop.Format("02.01.2006 15:04:05.000") + `",`

		response += `"tracks":` + `[`

		// THAT 5 SHOULD BE ABLE TO CHANGE DYNAMICALLY
		response += trackArray // Maximum of 5 tracks

		response += `],`

		response += `"processing":` + `"` + strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 2, 64) + `ms"`
		response += gmlCB

		fmt.Fprintln(w, response)

	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}


// Timestamps for ticker API struct
type Timestamps struct {
	latestTimestamp      time.Time
	oldestTimestamp      time.Time
	oldestNewerTimestamp time.Time
}



// Return the latest timestamp
func latestTimestamp(resultTracks []Track) time.Time {
	var latestTimestamp time.Time // Create a variable to store the most recent track added

	for _, val := range resultTracks { // Iterate every track to find the most recent track added
		if val.TimeRecorded.After(latestTimestamp) { // If current track timestamp is after the current latestTimestamp...
			latestTimestamp = val.TimeRecorded // Set that one as the latestTimestamp
		}
	}

	return latestTimestamp
}

// Return the oldest timestamp
func oldestTimestamp(resultTracks []Track) time.Time {

	// Just the first time, add the first found timestamp
	// After that, check that one against the other timestamps in the slice
	// If there is none, JSON response will be an empty string ""
	// If there is one timestamp, that one is the oldest timestamp as well

	var oldestTimestamp time.Time // Create a variable to store the oldest track added

	for key, val := range resultTracks { // Iterate every track to find the oldest track added

		// Assign to oldestTimestamp a value, but just once
		// Then we check it against other timestamps of other tracks in the slice
		if key == 0 {
			oldestTimestamp = val.TimeRecorded
		}

		if val.TimeRecorded.Before(oldestTimestamp) { // If current track timestamp is before the current latestTimestamp...
			oldestTimestamp = val.TimeRecorded // Set that one as the latestTimestamp
		}
	}

	return oldestTimestamp
}

// Return the oldest timestamp which is newer than input timestamp
func oldestNewerTimestamp(inputTS string, resultTracks []Track) time.Time {

	ts := time.Now()
	testTs := ts

	parsedTime, _ := time.Parse("02.01.2006 15:04:05.000", inputTS) // Parse the string into time

	for _, val := range resultTracks { // Iterate every track to find the most recent track added
		if val.TimeRecorded.After(parsedTime) && val.TimeRecorded.Before(ts) { // If current track timestamp is after the current latestTimestamp...
			ts = val.TimeRecorded // Set that one as the latestTimestamp
		}
	}

	if testTs.Equal(ts) {
		return time.Time{}
	}

	return ts
}

func tickerTimestamps(inputTS string) Timestamps {
	conn := mongoConnect()
	resultTracks := getAllTracks(conn, true)

	timestamps := Timestamps{}

	timestamps.latestTimestamp = latestTimestamp(resultTracks)
	timestamps.oldestTimestamp = oldestTimestamp(resultTracks)
	timestamps.oldestNewerTimestamp = oldestNewerTimestamp(inputTS, resultTracks)

	return timestamps
}



