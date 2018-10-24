package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)


/*

URLs for testing (IGC FILES):

	http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc
	http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc
	http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc
	http://skypolaris.org/wp-content/uploads/IGS%20Files/Medellin%20Guatemala.igc

*/


// Saving the time when the server started (i use this for later to calculate for how long the server is running)
var timeStarted = time.Now()



// ***DATA STRUCTURES USED IN THIS API*** //

// This structure holds Meta Information about the API
type MetaInformation struct {
	Uptime string		`json:"uptime"`
	Info string			`json:"info"`
	Version string 		`json:"version"`
}

// This structure holds the general informations about tracks posted in this API
type Track struct {
	Id string    				`json:"id"`
	Url string					`json:"url"`
	IgcTrack igc.Track 			`json:"igc_track"`
	TimeRecorded time.Time		`json:"time_recorded"`
}

// More specific info about tracks
type TrackInfo struct {
	HDate string			`json:"h_date"`
	Pilot string			`json:"pilot"`
	Glider string			`json:"glider"`
	GliderId string			`json:"glider_id"`
	TrackLength float64		`json:"track_length"`
	TrackSrcUrl string		`json:"track_src_url"`
}

type URLStruct struct {
	URL string `json:"url"`
}



//
//// Get the Port from the environment so we can run on Heroku
//func GetPort() string {
//	var port = os.Getenv("PORT")
//	// Set a default port if there is nothing in the environment
//	if port == "" {
//		port = "4747"
//		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
//	}
//	return ":" + port
//}



func main() {

	// I'm using Gorilla Mux router for routing different paths to assigned functions
	router := mux.NewRouter()

	// Paths for TracksAPI
	router.HandleFunc("/paragliding/", igcInfo)
	router.HandleFunc("/paragliding/api", getApi)
	router.HandleFunc("/paragliding/api/track", getApiIgc)
	router.HandleFunc("/paragliding/api/track/{id}", getApiIgcId)
	router.HandleFunc("/paragliding/api/track/{id}/{field}", getApiIgcField)

	// Paths for TickerAPI
	router.HandleFunc("/paragliding/api/ticker/latest", getApiTickerLatest)
	router.HandleFunc("/paragliding/api/ticker/", getApiTicker)
	router.HandleFunc("/paragliding/api/ticker/{timestamp}", getApiTickerTimestamp)

	// Paths for WebhookAPI
	router.HandleFunc("/paragliding/api/webhook/new_track/", webhookNewTrack)
	router.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}", webhookID)

	// Paths for AdminAPI
	router.HandleFunc("/paragliding/admin/api/tracks_count", adminApiTracksCount)
	router.HandleFunc("/paragliding/admin/api/tracks", adminApiTracks)

	//// Set http to listen and serve for different requests in the port found in the GetPort() function
	//err := http.ListenAndServe(GetPort(), router)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}


	log.Fatal(http.ListenAndServe(":8080", router))
}



// ***THE HANDLERS FOR THE CERTAIN PATHS*** //


// Handles path: GET /paragliding/
// If you request this path you get a 404 Not Found error for GET method and you get redirected to /paragliding/api path.
// For the other methods you get a status code of 501 Not Implemented
func igcInfo(w http.ResponseWriter, r *http.Request) {

	// For the methods that are not GET, it returns an http error 501 Not Implemented
	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	// Redirect to /paragliding/api
	http.Redirect(w, r, "/paragliding/api", http.StatusFound)

}


// This handler handles the path GET /igcinfo/api/.
// With the GET method on this path you get as a response Meta Information for the Api
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
	if len(urlVars) != 3 {
		http.Error(w, "400 - Bad Request, too many URLStruct arguments.", http.StatusBadRequest)
		return
	}

	// Creating an instance of MetaInformation stuct to save the data that is going to be Encoded with JSON
	metaInfo := &MetaInformation{}

	// The FormatTimeSince function gets the time that the server started as an argument and returns the formated time with ISO 8601 JSON format,
	// which represents the time for how long the server has been running.
	metaInfo.Uptime = FormatTimeSince(timeStarted)
	metaInfo.Info = "Service for IGC tracks"
	metaInfo.Version = "v2.2.0"

	// Encoding with JSON the meta information
	json.NewEncoder(w).Encode(metaInfo)
}


// Handle POST methods in which you enter an URL and you get back an ID for that inserted track on the database.
// Handle GET method where as a response you get all the IDs of the tracks saved in the database.
func getApiIgc(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

		w.Header().Set("Content-Type", "application/json")

		apiURL := &URLStruct{}

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
		igcFile.Id = strconv.Itoa(uniqueId)		// Converting int number to a string and saving it ti igcFile stuct
		igcFile.IgcTrack = track				// Saving all the other igcFile information
		igcFile.Url = apiURL.URL				// Saving the URL of that track, used later to check for duplicates before appending that file into igcFilesDB
		igcFile.TimeRecorded = time.Now()		// Saving the time that igcFile was recorded

		// Connecting to DB
		client := mongoConnect()

		// Specifying the specific collection which is going to be used
		collection := client.Database("igcFiles").Collection("track")

		// Checking for duplicates so that the user doesn't add into the database igc files with the same URL
		// If there is duplicates the function returns true, false otherwise
		duplicate := urlInMongo(igcFile.Url, collection)

		// If there are not duplicates Insert that track to the collection
		if !duplicate {

			res, err := collection.InsertOne(context.Background(), igcFile)
			if err != nil { log.Fatal(err) }
			id := res.InsertedID

			if id == nil{
				http.Error(w,"",300)
			}

			// Encoding the ID of the track that was just added to DB
			json.NewEncoder(w).Encode(igcFile.Id)

			// Counting the number of track at the moment
			trackCount := int32(countAllTracks(client))

			// Get tha latest webhook
			latestWebhook := getLatestWebhook(client)
			// Saving its minimal trigger value for later use
			minTriggerValue := latestWebhook.MinTriggerValue

			// If there is a webhook available or saved trigger that webhook, othewise just continue
			if latestWebhook.WebhookURL != "" {
				// Check according to minTriggerValue when to trigger the webhook
				if trackCount%minTriggerValue == 0 {
					// This function is called whenever a Track is registered in DB
					// The frequency of this function to be triggered depends on the minTriggerValue, which
					// indicates the frequency of updates - after how many tracks the webhook should be called
					triggerWhenTrackIsAdded(w, r, latestWebhook.WebhookURL)
				}
			}

		} else {

			// If there is another track in DB with that URL, get that track
			trackInDB := getTrack(client, igcFile.Url)

			// Notifying the user that the IGC File posted is already in our DB
			http.Error(w, "409 Conflict - The Igc File you entered is already in our database!", http.StatusConflict)

			// Printing that igcFile's ID
			fmt.Fprintln(w, "\nThe file you entered has the following ID: ", trackInDB.Id)

			return
		}




	case "GET":
		w.Header().Set("Content-Type", "application/json")

		urlVars := strings.Split(r.URL.Path, "/")
		if len(urlVars) != 4 {
			http.Error(w, "400 - Bad Request, too many URLStruct arguments.", http.StatusBadRequest)
			return
		}

		// Make a slice where are going to be saved all the IDs of the track in igcFiles database
		igcTrackIds := make([]string, 0, 0)

		client := mongoConnect()

		collection := client.Database("igcFiles").Collection("track")

		// Find all the documents in track collection
		cursor, err := collection.Find(context.Background(),nil,nil)
		if err != nil {
			log.Fatal(err)
		}

		// 'Close' the cursor
		defer cursor.Close(context.Background())

		track := Track{}

		// Point the cursor at whatever is found
		for cursor.Next(context.Background()) {
			// Decoding the findings
			err = cursor.Decode(&track)
			if err != nil {
				log.Fatal(err)
			}

			// Append all the track found ID's in igcTrackIDs slice
			igcTrackIds = append(igcTrackIds, track.Id)
		}

		// Encoding all IDs of the track in IgcFilesDB
		json.NewEncoder(w).Encode(igcTrackIds)


	default:
		// For other methods except GET and POST, requested in this handler you get this error
		http.Error(w, "Method not implemented yet", http.StatusNotImplemented)
		return

	}


}


// This handler handles the path: GET /igcinfo/api/igc/{id}/,
// which searches into igcFilesDB with the ID defined in that path
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


	client := mongoConnect()

	collection := client.Database("igcFiles").Collection("track")

	cursor, err := collection.Find(context.Background(),nil,nil)
	if err != nil {
		log.Fatal(err)
	}

	// 'Close' the cursor
	defer cursor.Close(context.Background())

	track := Track{}

	// Point the cursor at whatever is found
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&track)
		if err != nil {
			log.Fatal(err)
		}

		// Search into igcFiles found for the trackInfo with the requested ID and save the info from that trackInfo into the TrackInfo instance
		if track.Id == urlVars["id"] {
			trackInfo.HDate = track.IgcTrack.Date.String()
			trackInfo.Pilot = track.IgcTrack.Pilot
			trackInfo.Glider = track.IgcTrack.GliderType
			trackInfo.GliderId = track.IgcTrack.GliderID
			trackInfo.TrackLength = trackLength(track.IgcTrack)  // The trackLength function calculates the track length of a specific track, this function is defined at the end of this script
			trackInfo.TrackSrcUrl = track.Url

			json.NewEncoder(w).Encode(trackInfo)

			return
		}
	}

	// If the track if the requested ID doesn't exist in igcFilesDB, it returns an error
	http.Error(w, "404 - The trackInfo with that id doesn't exists in IGC Files", http.StatusNotFound)
	return

}


// This handler handles the path: GET /igcinfo/api/igc/{id}/{field}/
// You can search into igcFilesDB using the name of the field of the track you want to see, but first it has to be defined the ID of that particular track.
func getApiIgcField(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	urlVars := mux.Vars(r)

	regExId, _ := regexp.Compile("[0-9]+")					// Regular Expression for IDs
	regExField, _ := regexp.Compile("[a-zA-Z_]+")			// Regular Expression for Field

	if !regExId.MatchString(urlVars["id"]) {
		http.Error(w, "400 - Bad Request, you entered an invalid ID in URL.", http.StatusBadRequest)
		return
	}

	if !regExField.MatchString(urlVars["field"]) {
		http.Error(w, "400 - Bad Request, you entered an invalid Field in URL.", http.StatusBadRequest)
		return
	}


	client := mongoConnect()

	collection := client.Database("igcFiles").Collection("track")

	cursor, err := collection.Find(context.Background(),nil,nil)
	if err != nil {
		log.Fatal(err)
	}

	// 'Close' the cursor
	defer cursor.Close(context.Background())

	track := Track{}

	// Point the cursor at whatever is found
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&track)
		if err != nil {
			log.Fatal(err)
		}

		// Search into igcFilesDB for the track with the requested ID
		if track.Id == urlVars["id"] {

			// Mapping the track info into a Map
			fields := map[string]string {
				"pilot" :        track.IgcTrack.Pilot,
				"glider" :       track.IgcTrack.GliderType,
				"glider_id" :    track.IgcTrack.GliderID,
				"track_length" : FloatToString(trackLength(track.IgcTrack)),	// Calculate the track field for the specific track and convertin it to String
				"h_date" :       track.IgcTrack.Date.String(),
				"track_src_url": track.Url,
			}

			// Taking the field variable from the URL path and converting it to lower case to skip some potential errors
			field := urlVars["field"]
			field = strings.ToLower(field)

			// Searching into the map created above for the specific field that was requested
			if fieldData, ok := fields[field]; ok {
				// Encoding the data contained in the specific field saved in the map
				json.NewEncoder(w).Encode(fieldData)
				return
			} else {
				// If there is not a field like the one entered by the user. the user gets this error:
				http.Error(w, "400 - Bad Request, the field you entered is not on our database!", http.StatusBadRequest)
				return
			}

		}
	}

	// If the track with the requested ID isn't in our Database:
	http.Error(w, "404 - The track with that id doesn't exists in IGC Files", http.StatusNotFound)
	return
}




