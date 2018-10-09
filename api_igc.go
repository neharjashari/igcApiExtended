package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"log"
	"os"
)


/*

URLs for testing:

	http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc
	http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc
	http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc
	http://skypolaris.org/wp-content/uploads/IGS%20Files/Medellin%20Guatemala.igc

 */


// Saving the time when the server started (i use this for later to calculate for how long the server is running)
var timeStarted = time.Now()


// Creating a slice of type Track to save all the track that are posted in our api
var igcFilesDB []Track




// ***DATA STRUCTURES USED IN THIS API*** //

type MetaInformation struct {
	Uptime string		`json:"uptime"`
	Info string			`json:"info"`
	Version string 		`json:"version"`
}

type Track struct {
	Id string	`json:"id"`
	igcTrack igc.Track `json:"igc_track"`
}

type TrackInfo struct {
	HDate string			`json:"h_date"`
	Pilot string			`json:"pilot"`
	Glider string			`json:"glider"`
	GliderId string			`json:"glider_id"`
	TrackLength float64		`json:"track_length"`
}

type url struct {
	URL string `json:"url"`
}






func main() {

	// I'm using Gorilla Mux router for routing different paths to assigned functions
	router := mux.NewRouter()

	router.HandleFunc("/igcinfo/", igcInfo)
	router.HandleFunc("/igcinfo/api/", getApi)
	router.HandleFunc("/igcinfo/api/igc/", getApiIgc)
	router.HandleFunc("/igcinfo/api/igc/{id}/", getApiIgcId)
	router.HandleFunc("/igcinfo/api/igc/{id}/{field}/", getApiIgcField)

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, router))
}



// ***THE HANDLERS FOR THE CERTAIN PATHS*** //


// The first handler where if you request this path you get a 404 Not Found error for GET method. For the other methods you get a status code of 501 Not Implemented
func igcInfo(w http.ResponseWriter, r *http.Request) {

	// For the methods that are not GET, it returns an http error 501 Not Implemented
	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	http.Error(w, "404 - Page not found!", http.StatusNotFound)
	return

}


// This handler handles the path /igcinfo/api/. With the GET method on this path you get as a response Meta Information for the Api
func getApi(w http.ResponseWriter, r *http.Request) {
	// Set response content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// For the methods that are not GET, it returns an http error 501 Not Implemented
	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	// Check for URL malformed
	urlVars := strings.Split(r.URL.Path, "/")
	if len(urlVars) != 4 {
		http.Error(w, "400 - Bad Request, too many url arguments.", http.StatusBadRequest)
		return
	}

	// Creating an instance of MetaInformation stuct to save the data that is going to be Encoded with JSON
	metaInfo := &MetaInformation{}

	// The FormatSince function gets the time that the server started as an argument and returns the formated time with ISO 8601 JSON format,
	// which represents the time for how long the server has been running.
	metaInfo.Uptime = FormatSince(timeStarted)
	metaInfo.Info = "Service for IGC tracks"
	metaInfo.Version = "v1.0.0"

	json.NewEncoder(w).Encode(metaInfo)
}


// Handle POST methods in which you enter an URL and you get back an ID for that inserted track on the database.
// Handle GET method where as a response you get all the IDs of the tracks saved in the database.
func getApiIgc(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":
		w.Header().Set("Content-Type", "application/json")

		apiURL := &url{}

		// Decoding the URL sent by POST method into the apiURL variable
		var error = json.NewDecoder(r.Body).Decode(apiURL)
		if error != nil {
			fmt.Fprintln(w, "Error made: ", error)
			return
		}

		// Parsing that URL into igc library where the results are saved in the track variable
		track, err := igc.ParseLocation(apiURL.URL)
		if err != nil {
			fmt.Fprintln(w, "Error made: ", err)
			return
		}

		// Creating an "unique ID" as a random number for the specific track
		uniqueId := rand.Intn(1000)

		// Creating an instance of Track struct where it's saved the ID of the new track and the other info about it
		igcFile := Track{}
		igcFile.Id = strconv.Itoa(uniqueId)
		igcFile.igcTrack = track

		// Appending the added track into our database defined at the beginning of the program for all the tracks
		igcFilesDB = append(igcFilesDB, igcFile)


		json.NewEncoder(w).Encode(igcFile.Id)


	case "GET":
		w.Header().Set("Content-Type", "application/json")

		urlVars := strings.Split(r.URL.Path, "/")
		if len(urlVars) != 5 {
			http.Error(w, "400 - Bad Request, too many url arguments.", http.StatusBadRequest)
			return
		}

		// Make a slice where are going to be saved all the IDs of the track in igcFilesDB database
		igcTrackIds := make([]string, 0, 0)

		for i := range igcFilesDB {
			// Appending all the igcFilesDB IDs into the slice created earlier
			igcTrackIds = append(igcTrackIds, igcFilesDB[i].Id)
		}

		json.NewEncoder(w).Encode(igcTrackIds)


	default:
		http.Error(w, "Method not implemented yet", http.StatusNotImplemented)
		return

	}


}
	

// This handler handles the path: /igcinfo/api/igc/{id}/, which searches into igcFilesDB with the ID defined in that path
func getApiIgcId(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")


	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	// Using mux router to save the ID variable defined in the requested path
	urlVars := mux.Vars(r)


	// Regular Expression to check for ID validity, the ID can only be a set of integer numbers
	regExId, _ := regexp.Compile("[0-9]+")

	if !regExId.MatchString(urlVars["id"]) {
		http.Error(w, "400 - Bad Request, you entered an invalid ID in URL.", http.StatusBadRequest)
		return
	}

	// Creating an instance of TrackInfo structure where are saved the main information about the trackInfo
	trackInfo := &TrackInfo{}

	// Search into igcFilesDB for the trackInfo with the requested ID and save the info from that trackInfo into the TrackInfo instance
	for i := range igcFilesDB {
		if igcFilesDB[i].Id == urlVars["id"] {
			trackInfo.HDate = igcFilesDB[i].igcTrack.Date.String()
			trackInfo.Pilot = igcFilesDB[i].igcTrack.Pilot
			trackInfo.Glider = igcFilesDB[i].igcTrack.GliderType
			trackInfo.GliderId = igcFilesDB[i].igcTrack.GliderID
			// The trackLength function calculates the track length of a specific track, this function is defined at the end of this script
			trackInfo.TrackLength = trackLength(igcFilesDB[i].igcTrack)

			json.NewEncoder(w).Encode(trackInfo)

			return
		}
	}

	// If the track if the requested ID doesn't exist in igcFilesDB, it returns an error
	http.Error(w, "404 - The trackInfo with that id doesn't exists in IGC Files", http.StatusNotFound)
	return

}


// This handler handles the path: /igcinfo/api/igc/{id}/{field}/
// You can search into igcFilesDB using the name of the field of the track you want to see, but first it has to be defined the ID of that particular track.
func getApiIgcField(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	urlVars := mux.Vars(r)


	regExId, _ := regexp.Compile("[0-9]+")
	regExField, _ := regexp.Compile("[a-zA-Z_]+")

	if !regExId.MatchString(urlVars["id"]) {
		http.Error(w, "400 - Bad Request, you entered an invalid ID in URL.", http.StatusBadRequest)
		return
	}

	if !regExField.MatchString(urlVars["field"]) {
		http.Error(w, "400 - Bad Request, you entered an invalid Field in URL.", http.StatusBadRequest)
		return
	}

	// Search into igcFilesDB for the track with the requested ID
	for i := range igcFilesDB {
		if igcFilesDB[i].Id == urlVars["id"] {

			// Mapping the track info into a Map
			mapping := map[string]string {
				"pilot" :        igcFilesDB[i].igcTrack.Pilot,
				"glider" :       igcFilesDB[i].igcTrack.GliderType,
				"glider_id" :    igcFilesDB[i].igcTrack.GliderID,
				"track_length" : FloatToString(trackLength(igcFilesDB[i].igcTrack)),
				"h_date" :       igcFilesDB[i].igcTrack.Date.String(),
			}

			// Taking the field variable from the URL path and converting it to lower case to skip some potential errors
			field := urlVars["field"]
			field = strings.ToLower(field)

			// Encoding the data contained in the specific field saved in the map
			if fieldData, ok := mapping[field]; ok {
				json.NewEncoder(w).Encode(fieldData)
				return
			} else {
				http.Error(w, "400 - Bad Request, the field you entered is not on our database!", http.StatusBadRequest)
				return
			}

		}
	}

	http.Error(w, "404 - The track with that id doesn't exists in IGC Files", http.StatusNotFound)
	return
}




// This function formats time with iso8601 json format
// i found this function online because i couldn't find a better way of formating with json format
func FormatSince(t time.Time) string {

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



// This function calculates the track length of each track based on the track Points.
func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}


func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 4, 64)
}