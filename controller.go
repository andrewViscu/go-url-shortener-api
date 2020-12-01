package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"math/rand"
	"encoding/json"
	"context"
	"os"
	"time"
	"log"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asaskevich/govalidator"
)

type URL struct {
	Original string `json:"original,omitempty"`
	Shorten  string `json:"shorten"`
}

var urlCollection = DBConnect().Database(os.Getenv("DB_DEFAULT_DATABASE")).Collection(os.Getenv("DB_COLLECTION"))

func Index(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results []primitive.M                      
	cur, err := urlCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	defer cur.Close(ctx) 
	for cur.Next(ctx) { 

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}


const (
	urlLength = 6
	lettersAndNumbersBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func isValidURL(url string) bool{
	return govalidator.IsURL(url)
}

func generateURL() string {
	b := make([]byte, urlLength)
    for i := range b {
        b[i] = lettersAndNumbersBytes[rand.Intn(len(lettersAndNumbersBytes))]
    }
    return string(b)
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var (
		url URL
		notFound error
		stringURL string
	)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		log.Print("JSON Decoder error: ", err)
	}
	if !isValidURL(url.Original) {
		http.Error(w, "Invalid url '" + url.Original + "'", 400)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	opts := options.FindOne().SetCollation(&options.Collation{})
	i := 0
	for notFound == nil{ //loop til you find not taken url 
		stringURL = generateURL()
		found := urlCollection.FindOne(ctx, bson.D{{Key: "shorten", Value: stringURL}}, opts)
		notFound = found.Err()
		i++
	}
	log.Println(stringURL)
	url.Shorten = stringURL
	urlCollection.InsertOne(ctx, url)
	json.NewEncoder(w).Encode(url)
}

func RedirectFromURL(w http.ResponseWriter, r *http.Request) {
	var result URL 
	stringURL := mux.Vars(r)["shorten_url"]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	opts := options.FindOne().SetCollation(&options.Collation{})
	res := urlCollection.FindOne(ctx, bson.D{{Key: "shorten",Value: stringURL}}, opts)
	err := res.Err()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"shortened url not found: ` + err.Error() + `", "code":404}`))
		return
	}
	res.Decode(&result)
	http.Redirect(w, r, result.Original, http.StatusFound)
}