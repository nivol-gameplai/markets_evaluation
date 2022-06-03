package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"log"
	fi "marketsEvaluation/fast_markets_operations"
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
	Score                 string `json:"score"`
	Awayteamid            string `json:"awayteamid"`
	Currentteampossession string `json:"currentteampossession"`
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
	score, err := strconv.ParseInt(data.Score, 10, 64)
	if data.Service == "freeze" {
		log.Println("market suspension service")
		fi.DetermineSuspensionStrategy("next_bucket", data.Currentteampossession,
			data.Hometeamid, data.Gameid, data.Awayteamid)
	} else if data.Service == "evaluation" {
		log.Println("market evaluation service")
		_, err = fi.FetchAndEvaluate(data.Gameid, data.Marketid, score)
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
	json.Unmarshal([]byte(string(e)), &data)
	score, _ := strconv.ParseInt(data.Score, 10, 64)
	log.Println(score)
	fmt.Println(reflect.TypeOf(score))
	fi.MarketSuspension(data.Gameid, data.Hometeamid, score)
	w.WriteHeader(http.StatusOK)
}

func marketsSuspensionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fi.MarketSuspension("28810038", "3416", -1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func marketsEvaluationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fi.FetchAndEvaluate("28810038", "home_2pt", 125)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
