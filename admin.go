package main

import (
	"fmt"
	"net/http"
)

func adminApiTracksCount(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	client := mongoConnect()

	fmt.Fprintln(w, "Current count of the tracks in DB is: ", countAllTracks(client))
}


func adminApiTracks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	client := mongoConnect()

	fmt.Fprintln(w, "Count of the tracks removed from DB is: ", countAllTracks(client))

	deleteAllTracks(client)

}
