package fast_markets_operations

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"log"
)

func teamMarketsEvaluation(odd map[string]interface{}, teamScore int64, additiveScore int64) (response *string) {
	fulfillmentScore := odd["fulfillment_score"].(int64)
	currentTeamScore := odd["current_team_score"].(int64)
	if (currentTeamScore + additiveScore) == fulfillmentScore {
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

func evaluateMarket(odd map[string]interface{}, playerOrTeamScore int64, returnOdd chan map[string]interface{}) {
	//TODO CHECK IF THE CURRENT SCORE ON THE ODD (THE ONE OF DURING THE ODD CREATION)
	//TODO + THE 2 POINTS ETC == FULFILLMENT SCORE AND CURRENT SCORE FROM MESSAGE AND IF MISSED HAVE A REPLAY FUNCTION
	//TODO RUN THE API AGAIN AND CHECK ALL ACTIVES AND RE-CLOSE BY REPLAYING ALL EVENTS FROM SCRATCH AND PUT STATUS TO DEFFERED
	//TODO THIS SHOULD BE THE LIVE EVENTS SERVICE RUN WITHOUT THE REDIS CHECK OF EVENT ID BEING PROCESSED, THIS CAN BE
	//TODO DONE CONTINIOUSLY THROUGHOUT THE GAME , TO MAKE SURE ALL MARKETS ARE PROCESSED WHILE ACTIVE

	//TODO ALOSSSSSOOOO WHEN MARKETS CLOSED AS WON OR LOST SEND MESSAGE TO EVAL BETS AS NOW WE ARE SURE MARKETS ARE CLOSED :)
	var oddState *string
	oddState = nil
	marketID := odd["market_id"].(string)
	additiveScoreMap := marketSelection(marketID)
	scoreToAdd := additiveScoreMap["pointsToEval"].(int64)
	isItTeamOrPlayer := additiveScoreMap["teamorplayer"].(string)
	if isItTeamOrPlayer == "team" {
		oddState = teamMarketsEvaluation(odd, playerOrTeamScore, scoreToAdd)
	}
	if oddState != nil {
		odd["status"] = "CLOSED"
		odd["is_active"] = false
		odd["result"] = oddState
	}
	returnOdd <- odd
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
	returnOdd := make(chan map[string]interface{})

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		col := client.Collection("Odds")
		query = col.Where("game_id", "==", gameId).Where("is_active", "==",
			true).Where("market_id", "==", marketId)

		docs := tx.Documents(query)
		for {
			doc, err := docs.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
			}

			go evaluateMarket(doc.Data(), homeOrAwayScore, returnOdd)
			err = tx.Set(doc.Ref, <-returnOdd, firestore.MergeAll)
			if err != nil {
				return err
			}
		}
		close(returnOdd)
		return nil
	})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil, nil
}
