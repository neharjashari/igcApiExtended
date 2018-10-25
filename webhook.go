package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// *** WEBHOOK API *** //

// Webhook structure holds the info needed to register a webhook for later use
type Webhook struct {
	WebhookURL      string `json:"webhookURL"`
	MinTriggerValue int32  `json:"minTriggerValue"`
	WebhookID       string `json:"webhook_id"`
}

// WebhookContent keeps the webhook content to be send to Discord
type WebhookContent struct {
	TLatest    string   `json:"t_latest"`
	Tracks     []string `json:"tracks"`
	Processing string   `json:"processing"`
}

// Handles path: POST /api/webhook/new_track/
// Registration of new webhook for notifications about tracks being added to the system.
// Returns the details about the registration. The response contains the ID of the created resource
// The webhookURL is required parameter of the request.
// MinTriggerValue indicates the frequency of updates - after how many new tracks the webhook should be called.
func webhookNewTrack(w http.ResponseWriter, r *http.Request) {

	// It only works with POST requests
	if r.Method != "POST" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	webhook := Webhook{}

	// Decoding the URL sent by POST method into the apiURL variable
	var error = json.NewDecoder(r.Body).Decode(&webhook)
	if error != nil {
		fmt.Fprintln(w, "Error made: ", error)
		return
	}

	conn := mongoConnect()
	db := conn.Database("igcfiles")   // igcFiles Database
	coll := db.Collection("webhooks") // webhooks Collection

	// Check if Webhook exists
	cursor, err := coll.Find(context.Background(),
		bson.NewDocument(bson.EC.String("webhookurl", webhook.WebhookURL)))
	if err != nil {
		log.Fatal(err)
	}

	// 'Close' the cursor
	defer cursor.Close(context.Background())

	webhookInDB := Webhook{}

	// Point the cursor at whatever is found
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&webhookInDB)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintln(w, "The webhook you entered has been updated and has this ID: ", webhookInDB.WebhookID)

		// If the webhook is already in the DB, then update the minTriggerValue because that one can be changed even after
		// the webhook has been registered. But the ID doesn't change
		_, err := coll.UpdateOne(context.Background(),
			bson.NewDocument(
				bson.EC.String("webhookurl", webhook.WebhookURL),
			),
			bson.NewDocument(
				bson.EC.SubDocumentFromElements("$set",
					bson.EC.Int32("mintriggervalue", webhook.MinTriggerValue),
				),
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// Create an ID for the new webhook
	uniqueID := rand.Intn(1000)
	webhook.WebhookID = strconv.Itoa(uniqueID)

	// Insert the webhook if this one isn't in the Database
	_, err = coll.InsertOne(context.Background(), webhook)

	if err != nil {
		log.Fatal(err)
	}

	// Encoding the ID of the track that was just added to DB
	json.NewEncoder(w).Encode(webhook.WebhookID)

}

// Handles path: /api/webhook/new_track/<webhook_id>
func webhookID(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	// If the request is of GET type, then return the webhook registered with that ID
	case "GET":
		w.Header().Set("Content-Type", "application/json")

		urlVars := mux.Vars(r)

		client := mongoConnect()

		collection := client.Database("igcfiles").Collection("webhooks")

		cursor, err := collection.Find(context.Background(),
			bson.NewDocument(bson.EC.String("webhookid", urlVars["webhook_id"])))
		if err != nil {
			log.Fatal(err)
		}

		// 'Close' the cursor
		defer cursor.Close(context.Background())

		webhook := Webhook{}

		// Point the cursor at whatever is found
		for cursor.Next(context.Background()) {
			err = cursor.Decode(&webhook)
			if err != nil {
				log.Fatal(err)
			}

			json.NewEncoder(w).Encode(webhook)

			return

		}

		// If the webhook with the requested ID doesn't exist in the collection, return an error
		http.Error(w, "404 - The webhook with that ID doesn't exists in our Database", http.StatusNotFound)
		return

		// If the request is of DELETE type, then delete the webhook with the specified ID
	case "DELETE":
		w.Header().Set("Content-Type", "application/json")

		urlVars := mux.Vars(r)

		client := mongoConnect()

		collection := client.Database("igcfiles").Collection("webhooks")

		cursor, err := collection.Find(context.Background(),
			bson.NewDocument(bson.EC.String("webhookid", urlVars["webhook_id"])))
		if err != nil {
			log.Fatal(err)
		}

		// 'Close' the cursor
		defer cursor.Close(context.Background())

		webhook := Webhook{}

		// Point the cursor at whatever is found
		for cursor.Next(context.Background()) {
			err = cursor.Decode(&webhook)
			if err != nil {
				log.Fatal(err)
			}

			json.NewEncoder(w).Encode(webhook)

			// Delete the webhook that was found
			deleteWebhook(client, webhook.WebhookID)

			return

		}

		// If the webhook with the requested ID doesn't exist in the collection, return an error
		http.Error(w, "404 - The webhook with that ID doesn't exists in our Database", http.StatusNotFound)
		return

	default:
		// For other methods except GET and DELETE, requested in this handler you get this error
		http.Error(w, "Method not implemented yet", http.StatusNotImplemented)
		return

	}

}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// This function is called whenever a Track is registered in DB
// The frequency of this function to be triggered depends on the minTriggerValue, which
// indicates the frequency of updates - after how many tracks the webhook should be called
func triggerWhenTrackIsAdded(w http.ResponseWriter, r *http.Request) {

	// Connect to DB
	clientDB := mongoConnect()

	resultWebhooks := getAllWebhooks(clientDB)

	for _, val := range resultWebhooks {

		// Counting the number of track at the moment
		trackCount := int32(countAllTracks(clientDB))

		// Saving its minimal trigger value for later use
		minTriggerValue := val.MinTriggerValue

		// Check according to minTriggerValue when to trigger the webhook
		if trackCount%minTriggerValue == 0 {

			// Creating an instance of WebhookContent stuct
			webhookInfo := &WebhookContent{}

			processStart := time.Now() // Track when the process started

			timestamps := tickerTimestamps("")

			// Saving the latest added timestamp of the entire collection
			webhookInfo.TLatest = timestamps.latestTimestamp.String()

			// Creating a slice where all the IDs of track in DB are going to be saved
			WebhookInfoTrackIDs := make([]string, 0, 0)

			collection := clientDB.Database("igcfiles").Collection("track")

			// Find all the documents(tracks) in that collection
			cursor, err := collection.Find(context.Background(), nil, nil)
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

				// Append all the IDs of tracks in our DB to this slice
				WebhookInfoTrackIDs = append(WebhookInfoTrackIDs, track.ID)

			}

			webhookInfo.Tracks = WebhookInfoTrackIDs

			// Formating the processing time, time in ms of how long it took to process the request
			webhookInfo.Processing = strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 2, 64) + " ms"

			// Formating the content to be send to Discord as JSON
			content := "```css"
			content += "\n{ \n\t\"t_latest\" : \"" + webhookInfo.TLatest + "\" ,"
			content += " \n\t\"tracks\" : [ " + strings.Join(webhookInfo.Tracks, ", ") + " ] ,"
			content += " \n\t\"processing\" : \"" + webhookInfo.Processing + "\" \n}\n"
			content += "```"

			// Adding the values to URL
			data := url.Values{}
			data.Set("username", "TrackAdded")
			data.Add("content", content)

			u, _ := url.ParseRequestURI(val.WebhookURL)
			urlStr := u.String()

			client := &http.Client{}

			// Creating a new POST request to the webhook URL and sending the specified data to be printed in Discord
			r, err = http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
			if err != nil {
				fmt.Fprintln(w, "Error constructing the POST request, ", err)
			}

			// Specifying the request header parameters to send the data as JSON
			r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

			resp, err := client.Do(r)
			if err != nil {
				fmt.Fprintln(w, "Error executing the POST request, ", err)
			}

			defer resp.Body.Close()

		}

	}

}

// Get the latest webhook in DB
func getLatestWebhook(client *mongo.Client) Webhook {
	db := client.Database("igcfiles")
	collection := db.Collection("webhooks")

	cursor, err := collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	resWebhook := Webhook{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resWebhook)
		if err != nil {
			log.Fatal(err)
		}
	}

	return resWebhook

}

// Delete webhook with the ID specified in function parameters
func deleteWebhook(client *mongo.Client, webhookID string) {
	db := client.Database("igcfiles")
	collection := db.Collection("webhooks")

	// Delete the webhook
	collection.DeleteOne(
		context.Background(), bson.NewDocument(
			bson.EC.String("webhookid", webhookID),
		),
	)
}

// Get all webhooks
func getAllWebhooks(client *mongo.Client) []Webhook {
	db := client.Database("igcfiles")
	collection := db.Collection("webhooks")

	var cursor mongo.Cursor
	var err error
	cursor, err = collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	resWebhooks := []Webhook{}
	resWebhook := Webhook{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resWebhook)
		if err != nil {
			log.Fatal(err)
		}
		resWebhooks = append(resWebhooks, resWebhook) // Append each resTrack to resTracks slice
	}

	return resWebhooks
}

var latestTrackCounter = 1

//
func clockTrigger(w http.ResponseWriter, r *http.Request) {

	clientDB := mongoConnect()

	currentTrackCount := int(countAllTracks(clientDB))

	if latestTrackCounter != currentTrackCount {

		resultWebhooks := getAllWebhooks(clientDB)
		allTracks := getAllTracks(clientDB)

		allTracks = allTracks[latestTrackCounter-1:]

		for _, val := range resultWebhooks {

			webhookInfo := &WebhookContent{}

			processStart := time.Now() // Track when the process started

			timestamps := tickerTimestamps("")

			webhookInfo.TLatest = timestamps.latestTimestamp.String()

			WebhookInfoTrackIDs := make([]string, 0, 0)

			collection := clientDB.Database("igcfiles").Collection("track")

			cursor, err := collection.Find(context.Background(), nil, nil)
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
				WebhookInfoTrackIDs = append(WebhookInfoTrackIDs, track.ID)

			}

			webhookInfo.Tracks = WebhookInfoTrackIDs

			webhookInfo.Processing = strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 2, 64) + " ms"

			content := "\n"
			content += " Latest added track at: " + webhookInfo.TLatest + " ."
			content += " \nNew tracks are: [ " + strings.Join(webhookInfo.Tracks, ", ") + " ] ."
			content += " \nThe request took: " + webhookInfo.Processing + " time to process .\n"

			data := url.Values{}
			data.Set("username", "tracks")
			data.Add("content", content)

			u, _ := url.ParseRequestURI(val.WebhookURL)
			urlStr := u.String()

			client := &http.Client{}
			r, err = http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
			if err != nil {
				fmt.Fprintln(w, "Error constructing the POST request, ", err)
			}

			// Specifying the request header parameters to send the data as JSON
			r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

			resp, err := client.Do(r)
			if err != nil {
				fmt.Fprintln(w, "Error executing the POST request, ", err)
			}

			defer resp.Body.Close()

		}

		latestTrackCounter = currentTrackCount

	}

}
