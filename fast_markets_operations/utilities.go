package fast_markets_operations

var ScoreFulfillments = newScoreFulfillmentRegistry()

func newScoreFulfillmentRegistry() *scoreFulFillRegistry {
	return &scoreFulFillRegistry{
		away_2pt: 2,
		away_3pt: 3,
		away_ft:  1,
		home_2pt: 2,
		home_3pt: 3,
		home_ft:  1,
	}
}

type scoreFulFillRegistry struct {
	away_2pt int64
	away_3pt int64
	away_ft  int64
	home_2pt int64
	home_3pt int64
	home_ft  int64
}

func marketSelection(market string) int64 {
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
		return 0
	}
}

//func main() {
//	fmt.Println(ScoreFulfillments.away_2pt)
//}
