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




// *** WEBHOOK HANDLERS *** //


func webhookNewTrack(w http.ResponseWriter, r *http.Request) {

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
	db := conn.Database("igcFiles")
	coll := db.Collection("webhooks")

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
	uniqueId := rand.Intn(1000)
	webhook.WebhookID = strconv.Itoa(uniqueId)

	// Insert webhook
	_, err = coll.InsertOne(context.Background(), webhook)

	if err != nil {
		log.Fatal(err)
	}


	// Encoding the ID of the track that was just added to DB
	json.NewEncoder(w).Encode(webhook.WebhookID)

}


func webhookID(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		w.Header().Set("Content-Type", "application/json")

		urlVars := mux.Vars(r)

		client := mongoConnect()

		collection := client.Database("igcFiles").Collection("webhooks")

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

		// If the track if the requested ID doesn't exist in igcFilesDB, it returns an error
		http.Error(w, "404 - The webhook with that ID doesn't exists in our Database", http.StatusNotFound)
		return


	case "DELETE":
		w.Header().Set("Content-Type", "application/json")

		urlVars := mux.Vars(r)

		client := mongoConnect()

		collection := client.Database("igcFiles").Collection("webhooks")

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

			deleteWebhook(client, webhook.WebhookID)

			return

		}

		// If the track if the requested ID doesn't exist in igcFilesDB, it returns an error
		http.Error(w, "404 - The webhook with that ID doesn't exists in our Database", http.StatusNotFound)
		return


	default:
		// For other methods except GET and DELETE, requested in this handler you get this error
		http.Error(w, "Method not implemented yet", http.StatusNotImplemented)
		return

	}

}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


func triggerWhenTrackIsAdded(w http.ResponseWriter, r *http.Request, webhookURL string) {

	webhookInfo := &WebhookContent{}

	processStart := time.Now() // Track when the process started


	timestamps := tickerTimestamps("")

	webhookInfo.TLatest = timestamps.latestTimestamp.String()


	WebhookInfoTrackIDs := make([]string, 0, 0)

	clientDB := mongoConnect()

	collection := clientDB.Database("igcFiles").Collection("track")

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
		WebhookInfoTrackIDs = append(WebhookInfoTrackIDs, track.Id)

	}

	webhookInfo.Tracks = WebhookInfoTrackIDs

	webhookInfo.Processing = strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 2, 64) + " ms"


	content := "```css"
	content += "\n{ \n\t\"t_latest\" : \"" + webhookInfo.TLatest + "\" ,"
	content += " \n\t\"tracks\" : [ " + strings.Join(webhookInfo.Tracks, ", ") + " ] ,"
	content += " \n\t\"processing\" : \"" + webhookInfo.Processing + "\" \n}\n"
	content += "```"

	data := url.Values{}
	data.Set("username", "tracks")
	data.Add("content", content)

	u, _ := url.ParseRequestURI(webhookURL)
	urlStr := u.String()

	client := &http.Client{}
	r, err = http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
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

	defer resp.Body.Close()

}


func getLatestWebhook(client *mongo.Client) Webhook {
	db := client.Database("igcFiles")
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


// Delete webhook
func deleteWebhook(client *mongo.Client, webhookID string) {
	db := client.Database("igcFiles")
	collection := db.Collection("webhooks")

	// Delete the webhook
	collection.DeleteOne(
		context.Background(), bson.NewDocument(
			bson.EC.String("webhookid", webhookID),
		),
	)
}
