package NBA

import (
	"cloud.google.com/go/firestore"
	"sync"
	"time"
)

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
		//	TODO SEND MESSAGE FOR BETS EVAL
	}
	return nil
}

func EvaluateMarket(odd map[string]interface{}, playerOrTeamScore int64,
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
	additiveScoreMap := MarketSelection(marketID)
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
