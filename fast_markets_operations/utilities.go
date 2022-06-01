package fast_markets_operations

var ScoreFulfillments = newScoreFulfillmentRegistry()

func newScoreFulfillmentRegistry() *scoreFulFillRegistry {

	return &scoreFulFillRegistry{
		away_2pt: map[string]interface{}{
			"pointsToEval": 2,
			"teamorplayer": "team",
		},
		away_3pt: map[string]interface{}{
			"pointsToEval": 3,
			"teamorplayer": "team",
		},
		away_ft: map[string]interface{}{
			"pointsToEval": 1,
			"teamorplayer": "team",
		},
		home_2pt: map[string]interface{}{
			"pointsToEval": 2,
			"teamorplayer": "team",
		},
		home_3pt: map[string]interface{}{
			"pointsToEval": 3,
			"teamorplayer": "team",
		},
		home_ft: map[string]interface{}{
			"pointsToEval": 1,
			"teamorplayer": "team",
		},
	}
}

type scoreFulFillRegistry struct {
	away_2pt map[string]interface{}
	away_3pt map[string]interface{}
	away_ft  map[string]interface{}
	home_2pt map[string]interface{}
	home_3pt map[string]interface{}
	home_ft  map[string]interface{}
}

func marketSelection(market string) map[string]interface{} {
	if market == "away_2pt" {
		return ScoreFulfillments.away_2pt
	} else if market == "away_3pt" {
		return ScoreFulfillments.away_3pt
	} else if market == "away_ft" {
		return ScoreFulfillments.away_ft
	} else if market == "home_2pt" {
		return ScoreFulfillments.home_2pt
	} else if market == "home_3pt" {
		return ScoreFulfillments.home_3pt
	} else if market == "home_ft" {
		return ScoreFulfillments.home_ft
	} else {
		return map[string]interface{}{
			"pointsToEval": 0,
			"teamorplayer": "",
		}
	}
}

//func main() {
//	fmt.Println(ScoreFulfillments.away_2pt)
//}
