package main

import (
	"fmt"
	"net/http"
)

// Handles path: GET /admin/api/tracks_count
// Returns the current count of all tracks in the DB
func adminApiTracksCount(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	client := mongoConnect()

	fmt.Fprintln(w, "Current count of the tracks in DB is: ", countAllTracks(client))
}

// Handles path: DELETE /admin/api/track
// It only works with DELETE method, and this handler deletes all tracks in the DB
func adminApiTracks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	client := mongoConnect()

	// Notifying the admin first for the current count of the track
	fmt.Fprintln(w, "Count of the tracks removed from DB is: ", countAllTracks(client))

	// Deleting all the track in DB
	deleteAllTracks(client)

}
