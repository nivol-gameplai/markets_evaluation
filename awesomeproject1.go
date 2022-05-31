package main

import (
	"fmt"
	"io"
	"log"
	fi "marketsEvaluation/fast_markets_operations"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, "App Server Running\n")
		if err != nil {
			return
		}
	})
	http.HandleFunc("/markets/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fi.FetchAndEvaluate("28810038", "away_2pt", 111)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}
