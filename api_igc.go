package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"log"
	"math/rand"
	"net/http"
	"net/url"
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


// Creating a slice of type Track to save all the track that are posted in our api
var igcFilesDB []Track




// ***DATA STRUCTURES USED IN THIS API*** //

type MetaInformation struct {
	Uptime string		`json:"uptime"`
	Info string			`json:"info"`
	Version string 		`json:"version"`
}

type Track struct {
	Id       string    	`json:"id"`
	IgcTrack igc.Track 	`json:"igc_track"`
	Url      string		`json:"URLStruct"`
	Timestamp string 	`json:"timestamp"`
}

type TrackInfo struct {
	HDate string			`json:"h_date"`
	Pilot string			`json:"pilot"`
	Glider string			`json:"glider"`
	GliderId string			`json:"glider_id"`
	TrackLength float64		`json:"track_length"`
	TrackSrcUrl string		`json:"track_src_url"`
}

type URLStruct struct {
	URL string `json:"URLStruct"`
}



type Webhook struct {
	WebhookURL string		`json:"webhook_url"`
	MinTriggerValue int		`json:"min_trigger_value"`
	WebhookID string		`json:"webhook_id"`
}

type WebhookInfo struct {
	TLatest string		`json:"t_latest"`
	Tracks string		`json:"tracks"`
	Processing string	`json:"processing"`
}

type Discord struct {
	Username string		`json:"username"`
	Content string		`json:"content"`
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

	router.HandleFunc("/paragliding/", igcInfo)
	router.HandleFunc("/paragliding/api", getApi)
	router.HandleFunc("/paragliding/api/track", getApiIgc)
	router.HandleFunc("/paragliding/api/track/{id}", getApiIgcId)
	router.HandleFunc("/paragliding/api/track/{id}/{field}", getApiIgcField)
	router.HandleFunc("/paragliding/api/ticker/latest", getApiTickerLatest)
	router.HandleFunc("/paragliding/api/ticker/", getApiIgc)
	router.HandleFunc("/paragliding/api/ticker/{timestamp}", getApiIgc)

	router.HandleFunc("/paragliding/api/webhook/new_track/", webhookNewTrack)
	router.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}", webhookID)

	//// Set http to listen and serve for different requests in the port found in the GetPort() function
	//err := http.ListenAndServe(GetPort(), router)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}


	log.Fatal(http.ListenAndServe(":8080", router))
}



// ***THE HANDLERS FOR THE CERTAIN PATHS*** //

var webhooksDB []Webhook

func webhookNewTrack(w http.ResponseWriter, r *http.Request) {

	//timeStartedRequest := time.Now()

	w.Header().Set("Content-Type", "application/json")

	webhook := Webhook{}
	webhookInfo := &WebhookInfo{}

	// Decoding the URL sent by POST method into the apiURL variable
	var error = json.NewDecoder(r.Body).Decode(&webhook)
	if error != nil {
		fmt.Fprintln(w, "Error made: ", error)
		return
	}

	uniqueId := rand.Intn(1000)
	webhook.WebhookID = strconv.Itoa(uniqueId)

	//webhookInfo.TLatest = igcFilesDB[len(igcFilesDB) - 1].Timestamp
	webhookInfo.TLatest = "sgddsgg"

	WebhookInfoTrackIDs := make([]string, 0, 0)
	tracksString := "["
	for i := range igcFilesDB {
		// Appending all the igcFilesDB IDs into the slice created earlier
		WebhookInfoTrackIDs = append(WebhookInfoTrackIDs, igcFilesDB[i].Id)
		tracksString += WebhookInfoTrackIDs[i] + ","
	}
	tracksString += "]"
	webhookInfo.Tracks = tracksString

	//timeElapsed := timeStartedRequest - time.Now()
	webhookInfo.Processing = time.Now().String()


	webhookURL := webhook.WebhookURL

	content := "```css"
	content += "\n{ \n\t\"t_latest\" : \"" + webhookInfo.TLatest + "\","
	content += " \n\t\"tracks\" : " + webhookInfo.Tracks + ","
	content += " \n\t\"processing\" : \"" + webhookInfo.Processing + "\" \n}\n"
	content += "```"

	//resource := "/user/"
	data := url.Values{}
	data.Set("username", "tracks")
	data.Add("content", content)

	u, _ := url.ParseRequestURI(webhookURL)
	//u.Path = resource
	urlStr := u.String() // 'https://api.com/user/'

	client := &http.Client{}
	r, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		fmt.Fprintln(w,"Error constructing the POST request, ", err)
	}
	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Fprintln(w,"Error executing the POST request, ", err)
	}

	//fmt.Println(resp.Status)

	defer resp.Body.Close()

	// Checking for duplicates so that the user doesn't add into the database igc files with the same URL
	for i := range webhooksDB {
		if webhooksDB[i].WebhookID == webhook.WebhookID {
			// If there is another file in igcFilesDB with that URL return and tell the user that that IGC FILE is already in the database
			http.Error(w, "409 Conflict - The Igc File you entered is already in our database!", http.StatusConflict)
			fmt.Fprintln(w, "\nThe file you entered has the following ID: ", webhooksDB[i].WebhookID)
			return
		}
	}

	// Appending the added track into our database defined at the beginning of the program for all the tracks
	webhooksDB = append(webhooksDB, webhook)

	json.NewEncoder(w).Encode(webhook.WebhookID)

}


func webhookID(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		w.Header().Set("Content-Type", "application/json")

		urlVars := mux.Vars(r)
		currentWebhook := Webhook{}

		for i := range webhooksDB {
			if webhooksDB[i].WebhookID == urlVars["webhook_id"] {
				currentWebhook = webhooksDB[i]
			}
		}

		// Encoding all IDs of the track in IgcFilesDB
		json.NewEncoder(w).Encode(currentWebhook)

	case "DELETE":
		w.Header().Set("Content-Type", "application/json")

		urlVars := mux.Vars(r)
		currentWebhook := Webhook{}

		for i := range webhooksDB {
			if webhooksDB[i].WebhookID == urlVars["webhook_id"] {
				currentWebhook = webhooksDB[i]
				webhooksDB = append(webhooksDB[:i], webhooksDB[i+1:]...)
			}
		}

		json.NewEncoder(w).Encode(currentWebhook)

	default:
		// For other methods except GET and DELETE, requested in this handler you get this error
		http.Error(w, "Method not implemented yet", http.StatusNotImplemented)
		return

	}

}



// The first handler where if you request this path you get a 404 Not Found error for GET method. For the other methods you get a status code of 501 Not Implemented
func igcInfo(w http.ResponseWriter, r *http.Request) {

	// For the methods that are not GET, it returns an http error 501 Not Implemented
	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	http.Redirect(w, r, "/paragliding/api", http.StatusFound)

	//// As described	in the assignment when you require the GET method in the root path you get a 404 Not Found error
	//http.Error(w, "404 - Page not found!", http.StatusNotFound)
	//return

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
	metaInfo.Version = "v1.0.0"

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
		igcFile.IgcTrack = track
		igcFile.Url = apiURL.URL				// Saving the URL of that track, used later to check for duplicates before appending that file into igcFilesDB

		/* TODO : Each track, ie. each IGC file uploaded to the system, will have a timestamp
		 represented as a LONG number, that must be unique and monotonic. This can be achieved
		 by storing a milisecond of the upload time to the server, although, you have to plan how
		 to make it thread safe and scalable. This is relevant for the ticker API. Hint - you could
		 use mongoDB IDs that are monotonic, but then you would have some security considerations -
		 think which ones. */
		igcFile.Timestamp = time.Now().String()



		// Checking for duplicates so that the user doesn't add into the database igc files with the same URL
		for i := range igcFilesDB {
			if igcFilesDB[i].Url == igcFile.Url {
				// If there is another file in igcFilesDB with that URL return and tell the user that that IGC FILE is already in the database
				http.Error(w, "409 Conflict - The Igc File you entered is already in our database!", http.StatusConflict)
				fmt.Fprintln(w, "\nThe file you entered has the following ID: ", igcFilesDB[i].Id)
				return
			}
		}


		// Appending the added track into our database defined at the beginning of the program for all the tracks
		igcFilesDB = append(igcFilesDB, igcFile)


		// Encoding the ID of the track that was just added to DB
		json.NewEncoder(w).Encode(igcFile.Id)


	case "GET":
		w.Header().Set("Content-Type", "application/json")

		urlVars := strings.Split(r.URL.Path, "/")
		if len(urlVars) != 4 {
			http.Error(w, "400 - Bad Request, too many URLStruct arguments.", http.StatusBadRequest)
			return
		}

		// Make a slice where are going to be saved all the IDs of the track in igcFilesDB database
		igcTrackIds := make([]string, 0, 0)

		for i := range igcFilesDB {
			// Appending all the igcFilesDB IDs into the slice created earlier
			igcTrackIds = append(igcTrackIds, igcFilesDB[i].Id)
		}

		// Encoding all IDs of the track in IgcFilesDB
		json.NewEncoder(w).Encode(igcTrackIds)


	default:
		// For other methods except GET and POST, requested in this handler you get this error
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
			trackInfo.HDate = igcFilesDB[i].IgcTrack.Date.String()
			trackInfo.Pilot = igcFilesDB[i].IgcTrack.Pilot
			trackInfo.Glider = igcFilesDB[i].IgcTrack.GliderType
			trackInfo.GliderId = igcFilesDB[i].IgcTrack.GliderID
			// The trackLength function calculates the track length of a specific track, this function is defined at the end of this script
			trackInfo.TrackLength = trackLength(igcFilesDB[i].IgcTrack)
			trackInfo.TrackSrcUrl = igcFilesDB[i].Url

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

	// Search into igcFilesDB for the track with the requested ID
	for i := range igcFilesDB {
		if igcFilesDB[i].Id == urlVars["id"] {

			// Mapping the track info into a Map
			fields := map[string]string {
				"pilot" :        igcFilesDB[i].IgcTrack.Pilot,
				"glider" :       igcFilesDB[i].IgcTrack.GliderType,
				"glider_id" :    igcFilesDB[i].IgcTrack.GliderID,
				"track_length" : FloatToString(trackLength(igcFilesDB[i].IgcTrack)),	// Calculate the track field for the specific track and convertin it to String
				"h_date" :       igcFilesDB[i].IgcTrack.Date.String(),
				"track_src_url": igcFilesDB[i].Url,
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



func getApiTickerLatest(w http.ResponseWriter, r *http.Request) {

	// For the methods that are not GET, it returns an http error 501 Not Implemented
	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	//// Check for URL malformed
	//urlVars := strings.Split(r.URL.Path, "/")
	//if len(urlVars) != 3 {
	//	http.Error(w, "400 - Bad Request, too many URLStruct arguments.", http.StatusBadRequest)
	//	return
	//}


	// Creating a type that holds the struct Track, which will later save the latest added track
	// Saving in latestTrack the last item on the igcFilesDB
	igcFilesLen := len(igcFilesDB)
	latestTrack := igcFilesDB[igcFilesLen - 1]

	timestamp := latestTrack.Timestamp

	json.NewEncoder(w).Encode(timestamp)
}






// This function formats time with iso8601 json format
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