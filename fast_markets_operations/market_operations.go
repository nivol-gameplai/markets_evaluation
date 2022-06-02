package fast_markets_operations

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"log"
	"sync"
	"time"
)

func MarketSuspension(gameId string, teamId string, currentScore int64) ([]byte, error) {
	ctx := context.Background()
	client := CreateClient(ctx)
	defer closeClient(client)
	var query firestore.Query
	col := client.Collection("Odds")
	var wg sync.WaitGroup

	if currentScore < 0 {

		query = col.Where("game_id", "==", gameId).Where("is_active", "==",
			true).Where("team_id", "==", teamId)
	} else {
		query = col.Where("game_id", "==", gameId).Where("is_active", "==",
			true).Where("team_id", "==", teamId).Where("current_team_score", "<=", currentScore)
	}

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docs := tx.Documents(query)
		for {
			doc, err := docs.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
			}
			wg.Add(1)

			go updateSuspensionOnDocs(doc.Data(), currentScore, tx, *doc.Ref, &wg)

		}
		wg.Wait()
		return nil
	})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil, nil
}

func updateSuspensionOnDocs(odd map[string]interface{}, currentScore int64,
	tx *firestore.Transaction, ref firestore.DocumentRef, wg *sync.WaitGroup) {
	defer wg.Done()
	if odd["freeze"] == true {

		odd["freeze"] = false
	} else {
		if currentScore < 0 {
			odd["freeze"] = true
		}
	}
	if currentScore >= 0 {
		odd["is_active"] = false
	}
	err := tx.Set(&ref, odd)
	if err != nil {
		return
	}

}

func teamMarketsEvaluation(odd map[string]interface{}, teamScore int64, additiveScore int) (response *string) {
	fulfillmentScore := odd["fulfillment_score"].(int64)
	currentTeamScore := odd["current_team_score"].(int64)
	if (currentTeamScore + int64(additiveScore)) == fulfillmentScore {
		if teamScore == fulfillmentScore {
			ret := "WON"
			return &ret
		} else {
			ret := "LOST"
			return &ret
		}
	}
	return nil
}

func evaluateMarket(odd map[string]interface{}, playerOrTeamScore int64,
	tx *firestore.Transaction, ref firestore.DocumentRef, wg *sync.WaitGroup) {
	//TODO CHECK IF THE CURRENT SCORE ON THE ODD (THE ONE OF DURING THE ODD CREATION)
	//TODO + THE 2 POINTS ETC == FULFILLMENT SCORE AND CURRENT SCORE FROM MESSAGE AND IF MISSED HAVE A REPLAY FUNCTION
	//TODO RUN THE API AGAIN AND CHECK ALL ACTIVES AND RE-CLOSE BY REPLAYING ALL EVENTS FROM SCRATCH AND PUT STATUS TO DEFFERED
	//TODO THIS SHOULD BE THE LIVE EVENTS SERVICE RUN WITHOUT THE REDIS CHECK OF EVENT ID BEING PROCESSED, THIS CAN BE
	//TODO DONE CONTINIOUSLY THROUGHOUT THE GAME , TO MAKE SURE ALL MARKETS ARE PROCESSED WHILE ACTIVE

	//TODO ALOSSSSSOOOO WHEN MARKETS CLOSED AS WON OR LOST SEND MESSAGE TO EVAL BETS AS NOW WE ARE SURE MARKETS ARE CLOSED :)
	defer wg.Done()
	var oddState *string
	oddState = nil
	marketID := odd["market_id"].(string)
	additiveScoreMap := marketSelection(marketID)
	scoreToAdd := additiveScoreMap["pointsToEval"].(int)
	isItTeamOrPlayer := additiveScoreMap["teamorplayer"].(string)
	if isItTeamOrPlayer == "team" {
		oddState = teamMarketsEvaluation(odd, playerOrTeamScore, scoreToAdd)
	}
	if oddState != nil {
		odd["status"] = "CLOSED"
		odd["is_active"] = false
		odd["result"] = oddState
		odd["modified"] = time.Now().UnixNano() / int64(time.Millisecond)
	}
	err := tx.Set(&ref, odd)
	if err != nil {
		return
	}
	//log.Println(<-returnOdd)
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

	var query firestore.Query

	var wg sync.WaitGroup

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		col := client.Collection("Odds")
		query = col.Where("game_id", "==", gameId).Where("status", "==",
			"open").Where("market_id", "==", marketId).Where("current_team_score", "<=", homeOrAwayScore)

		docs := tx.Documents(query)
		for {
			doc, err := docs.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
			}
			wg.Add(1)
			go evaluateMarket(doc.Data(), homeOrAwayScore, tx, *doc.Ref, &wg)

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
