package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"log"
	"marketsEvaluation/fast_markets_operations"
	"marketsEvaluation/fast_markets_operations/markets_suspension"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type Data struct {
	Gameid                string `json:"gameid"`
	Hometeamid            string `json:"hometeamid"`
	Marketid              string `json:"marketid"`
	Service               string `json:"service"`
	Metric                string `json:"metric"`     // the actual metric like score, assists etc
	Metrictype            string `json:"metrictype"` //the string representation of the metric like "score"
	Awayteamid            string `json:"awayteamid"`
	Currentteampossession string `json:"currentteampossession"`
	Sport                 string `json:"sport"`
	Playerid              string `json:"playerid"`
}

func listenToTriggers(w http.ResponseWriter, r *http.Request) {
	var e PubSubMessage
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Bad HTTP Request", http.StatusBadRequest)
		log.Printf("Bad HTTP Request: %v", http.StatusBadRequest)
		return
	}
	data := Data{}
	unmarshalErr := json.Unmarshal([]byte(string(e.Message.Data)), &data)
	if unmarshalErr != nil {
		// Error...
		log.Println("error unmarshall")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, err := strconv.ParseInt(data.Metric, 10, 64)
	if data.Service == "freeze" {
		log.Println("market suspension service")
		markets_suspension.DetermineSuspensionStrategy("next_bucket", data.Currentteampossession,
			data.Hometeamid, data.Gameid, data.Awayteamid, metric)
	} else if data.Service == "evaluation" {
		log.Println("market evaluation service")
		_, err = fast_markets_operations.FetchAndEvaluate(data.Gameid, data.Marketid, metric, data.Metrictype, data.Sport)
	}
	s := fmt.Sprintf("Hello, %s! ID: %s", "", string(r.Header.Get("Ce-Id")))
	_, err = fmt.Fprintln(w, s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", listenToTriggers)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/markets/", marketsEvaluationHandler)
	http.HandleFunc("/freeze/", marketsSuspensionHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := Data{}
	e := `{"gameid": "28810038", "teamid": "3416", "marketid": "NA","service": "freeze", "score":"-1"}`
	err := json.Unmarshal([]byte(string(e)), &data)
	if err != nil {
		return
	}
	score, _ := strconv.ParseInt(data.Metric, 10, 64)
	log.Println(score)
	fmt.Println(reflect.TypeOf(score))
	_, err = markets_suspension.MarketSuspension(data.Gameid, data.Hometeamid, score)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func marketsSuspensionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := markets_suspension.MarketSuspension("28810038", "3416", -1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func marketsEvaluationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fast_markets_operations.FetchAndEvaluate("28810038", "home_2pt", 125,
		"score", "NBA")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
