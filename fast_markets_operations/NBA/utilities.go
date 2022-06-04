package NBA

var ScoreFulfillments = newScoreFulfillmentRegistry()
var MetricQueryParameter = newMetricregistry()

type scoreFulFillRegistry struct {
	away_2pt map[string]interface{}
	away_3pt map[string]interface{}
	away_ft  map[string]interface{}
	home_2pt map[string]interface{}
	home_3pt map[string]interface{}
	home_ft  map[string]interface{}
}

type metricString struct {
	score      string
	assists    string
	threes     string
	twos       string
	freethrows string
	dunks      string
	blocks     string
}

func newMetricregistry() *metricString {

	return &metricString{
		score:      "current_team_score",
		assists:    "current_assists",
		threes:     "current_threes",
		twos:       "current_twos",
		freethrows: "current_freethrows",
		dunks:      "current_dunks",
		blocks:     "current_blocks",
	}
}

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

func MarketSelection(market string) map[string]interface{} {
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

func DetermineQueryMetricParameter(metricType string) string {
	if metricType == "score" {
		return MetricQueryParameter.score
	} else if metricType == "assists" {
		return MetricQueryParameter.assists
	} else if metricType == "threes" {
		return MetricQueryParameter.threes
	} else if metricType == "twos" {
		return MetricQueryParameter.twos
	} else if metricType == "freethrows" {
		return MetricQueryParameter.freethrows
	} else if metricType == "dunks" {
		return MetricQueryParameter.dunks
	} else if metricType == "blocks" {
		return MetricQueryParameter.blocks
	}
	return MetricQueryParameter.score
}
