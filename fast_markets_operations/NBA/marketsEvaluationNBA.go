package NBA

import (
	"cloud.google.com/go/firestore"
	"sync"
	"time"
)

// teamMarketsEvaluation veluates the markets based on the metric and fulfillmentMetric and current metric value
// i.e current score (current metric) = 125 if market = 2pts for score then if current metric (score) on odd doc
// is 123 and additive metric is 2 (as of 2pts evaluation) then the fulfillmentMetric is true 125 and the market is won
//TODO send messaged to client's for betslip evaluation
func teamMarketsEvaluation(odd map[string]interface{}, metric int64, additiveMetric int, metricType string) (response *string) {
	queryMapForMarket := DetermineQueryMetricParameter(metricType)
	currentMetricFieldString := queryMapForMarket["currentMetric"]
	fulfillmentMetricFieldString := queryMapForMarket["fulfillmentMetric"]
	fulfillmentMetric := odd[fulfillmentMetricFieldString].(int64)
	currentMetric := odd[currentMetricFieldString].(int64)
	if (currentMetric + int64(additiveMetric)) == fulfillmentMetric {
		if metric == fulfillmentMetric {
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

// EvaluateMarket receives an odd document as input, and the based on market id,
// it retrives fullfilment and market details before calling the MarketsEvaluation
// to determing outcome for the odd(market)
// finally it updates the odd doc
func EvaluateMarket(odd map[string]interface{}, playerOrTeamScore int64,
	tx *firestore.Transaction, ref firestore.DocumentRef, wg *sync.WaitGroup, metricType string) {

	defer wg.Done()
	var oddState *string
	oddState = nil
	marketID := odd["market_id"].(string)
	additiveScoreMap := MarketSelection(marketID)
	scoreToAdd := additiveScoreMap["metricToEval"].(int)
	isItTeamOrPlayer := additiveScoreMap["teamorplayer"].(string)
	if isItTeamOrPlayer == "team" {
		oddState = teamMarketsEvaluation(odd, playerOrTeamScore, scoreToAdd, metricType)
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
