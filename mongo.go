package main

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"strings"
	"time"
)

// *** DB METHODS *** //

// This function connects the API with Mongo Database and returns that connection
func mongoConnect() *mongo.Client {
	// Connect to MongoDB
	conn, err := mongo.Connect(context.Background(), "mongodb://localhost:27017", nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return conn
}

// Check if the track already exists in the database
func urlInMongo(url string, trackColl *mongo.Collection) bool {

	// Read the documents where the trackurl field is equal to url parameter
	cursor, err := trackColl.Find(context.Background(),
		bson.NewDocument(bson.EC.String("url", url)))
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
	}

	if track.URL == "" { // If there is an empty field, in this case, `url`, it means the track is not on the database
		return false
	}
	return true
}

// Returns the track with the specified URL as function parameters
func getTrack(client *mongo.Client, url string) Track {
	db := client.Database("igcFiles")    // igcFiles Database
	collection := db.Collection("track") // track Collection

	// Query collection to find the specific track with that URL
	cursor, err := collection.Find(context.Background(),
		bson.NewDocument(bson.EC.String("url", url)))

	if err != nil {
		log.Fatal(err)
	}

	resTrack := Track{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resTrack)
		if err != nil {
			log.Fatal(err)
		}
	}

	return resTrack

}

// Delete all tracks
func deleteAllTracks(client *mongo.Client) {
	db := client.Database("igcFiles")
	collection := db.Collection("track")

	// Delete the tracks
	collection.DeleteMany(context.Background(), bson.NewDocument())
}

// Count all tracks
func countAllTracks(client *mongo.Client) int64 {
	db := client.Database("igcFiles")
	collection := db.Collection("track")

	// Count the tracks
	count, _ := collection.Count(context.Background(), nil, nil)

	return count
}

// Return track names
// And also t_stop track
func returnTracks(n int) (string, time.Time) {
	var response string
	var tStop time.Time

	conn := mongoConnect()

	resultTracks := getAllTracks(conn)

	for key, val := range resultTracks { // Go through the slice
		response += `"` + val.ID + `",`
		if key == n-1 || key == len(resultTracks)-1 {
			tStop = val.TimeRecorded
			break
		}
	}

	// Get rid of that last `,` of JSON will freak out
	response = strings.TrimRight(response, ",")

	return response, tStop
}


// Get all tracks
func getAllTracks(client *mongo.Client) []Track {
	db := client.Database("igcFiles")
	collection := db.Collection("track")

	var cursor mongo.Cursor
	var err error

	cursor, err = collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	resTracks := []Track{}
	resTrack := Track{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resTrack)
		if err != nil {
			log.Fatal(err)
		}
		resTracks = append(resTracks, resTrack) // Append each resTrack to resTracks slice
	}

	return resTracks
}
