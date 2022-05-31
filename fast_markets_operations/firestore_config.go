package fast_markets_operations

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"log"
)

// [END admin_import]

func CreateClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "braided-triode-232313"
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func evaluateMarket(odd map[string]interface{}, homeOrAwayScore int64) map[string]interface{} {
	//TODO CHECK IF THE CURRENT SCORE ON THE ODD (THE ONE OF DURING THE ODD CREATION)
	//TODO + THE 2 POINTS ETC == FULFILLMENT SCORE AND CURRENT SCORE FROM MESSAGE AND IF MISSED HAVE A REPLAY FUNCTION
	//TODO RUN THE API AGAIN AND CHECK ALL ACTIVES AND RE-CLOSE BY REPLAYING ALL EVENTS FROM SCRATCH AND PUT STATUS TO DEFFERED
	fulfillmentScore := odd["fulfillment_score"].(int64)
	if homeOrAwayScore == fulfillmentScore {
		odd["result"] = "WON"
	} else {
		odd["result"] = "LOST"
	}
	odd["status"] = "CLOSED"
	odd["is_active"] = false

	return odd
}

func FetchAndEvaluate(gameId string, marketId string, homeOrAwayScore int64) ([]byte, error) {
	// Get a Firestore client.
	ctx := context.Background()
	client := CreateClient(ctx)
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
		}
	}(client)
	var odd map[string]interface{}
	var query firestore.Query

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		col := client.Collection("Odds")
		query = col.Where("game_id", "==", gameId).Where("is_active", "==",
			true).Where("market_id", "==", marketId).Limit(1)

		docs := tx.Documents(query)
		for {
			doc, err := docs.Next()
			if err != nil {
				if err == iterator.Done {

					break
				}
			}
			odd = evaluateMarket(doc.Data(), homeOrAwayScore)
			err = tx.Set(doc.Ref, odd, firestore.MergeAll)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil, nil
}
