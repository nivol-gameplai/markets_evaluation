package fast_markets_operations

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"log"
	"marketsEvaluation/fast_markets_operations/NBA"
	"marketsEvaluation/fast_markets_operations/configs"
	"sync"
)

func FetchAndEvaluate(gameId string, marketId string, metric int64, metricType string, sport string) ([]byte, error) {
	// Get a Firestore client.
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
		}
	}(client)

	var query firestore.Query

	var wg sync.WaitGroup
	var queryParameter string

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		col := client.Collection("Odds")
		if sport == "NBA" {
			queryParameter = NBA.DetermineQueryMetricParameter(metricType)
		} else if sport == "soccer" {
			log.Println("TBA soccer")
		}

		query = col.Where("game_id", "==", gameId).Where("status", "==",
			"open").Where("market_id", "==", marketId).Where(queryParameter, "<=", metric)

		docs := tx.Documents(query)
		for {
			doc, err := docs.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
			}
			wg.Add(1)
			if sport == "NBA" {
				go NBA.EvaluateMarket(doc.Data(), metric, tx, *doc.Ref, &wg)
			} else if sport == "soccer" {
				log.Println("TBA soccer")
			}

			if err != nil {
				return err
			}
		}
		wg.Wait()
		return nil
	})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil, nil
}
