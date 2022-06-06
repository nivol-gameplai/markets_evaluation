package NBA

var ScoreFulfillments = newScoreFulfillmentRegistry()
var MetricQueryParameter = newMetricregistry()

// a way to build enum values in go
// we create a struct with rules for each market
// that is initiated in newScoreFulfillmentRegistry
type scoreFulFillRegistry struct {
	away_2pt map[string]interface{}
	away_3pt map[string]interface{}
	away_ft  map[string]interface{}
	home_2pt map[string]interface{}
	home_3pt map[string]interface{}
	home_ft  map[string]interface{}
}

// a way to build enum values in go
// we create a struct with rules for metrics to be tracked on markets
// these can be expanded as we add more markets
// that is initiated in newMetricregistry
type metricString struct {
	score      map[string]string
	assists    map[string]string
	threes     map[string]string
	twos       map[string]string
	freethrows map[string]string
	dunks      map[string]string
	blocks     map[string]string
}

// initialization of ths struct/enum
func newMetricregistry() *metricString {

	return &metricString{
		score: map[string]string{
			"currentMetric":     "current_team_score",
			"fulfillmentMetric": "fulfillment_score",
		},
		assists: map[string]string{
			"currentMetric":     "current_assists",
			"fulfillmentMetric": "fulfillment_assists",
		},
		threes: map[string]string{
			"currentMetric":     "current_threes",
			"fulfillmentMetric": "fulfillment_threes",
		},
		twos: map[string]string{
			"currentMetric":     "current_twos",
			"fulfillmentMetric": "fulfillment_twos",
		},
		freethrows: map[string]string{
			"currentMetric":     "current_freethrows",
			"fulfillmentMetric": "fulfillment_freethrows",
		},
		dunks: map[string]string{
			"currentMetric":     "current_dunks",
			"fulfillmentMetric": "fulfillment_dunks",
		},
		blocks: map[string]string{
			"currentMetric":     "current_blocks",
			"fulfillmentMetric": "fulfillment_blocks",
		},
	}
}

// initialization of ths struct/enum
func newScoreFulfillmentRegistry() *scoreFulFillRegistry {

	return &scoreFulFillRegistry{
		away_2pt: map[string]interface{}{
			"metricToEval": 2,
			"teamorplayer": "team",
		},
		away_3pt: map[string]interface{}{
			"metricToEval": 3,
			"teamorplayer": "team",
		},
		away_ft: map[string]interface{}{
			"metricToEval": 1,
			"teamorplayer": "team",
		},
		home_2pt: map[string]interface{}{
			"metricToEval": 2,
			"teamorplayer": "team",
		},
		home_3pt: map[string]interface{}{
			"metricToEval": 3,
			"teamorplayer": "team",
		},
		home_ft: map[string]interface{}{
			"metricToEval": 1,
			"teamorplayer": "team",
		},
	}
}

// function to return maket values based on enum
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
			"metricToEval": 0,
			"teamorplayer": "",
		}
	}
}

// function to return metric values based on enum
func DetermineQueryMetricParameter(metricType string) map[string]string {
	if metricType == "score" {
		return MetricQueryParameter.score
	} else if metricType ==
		"assists" {
		return MetricQueryParameter.assists
	} else if metricType ==
		"threes" {
		return MetricQueryParameter.threes
	} else if metricType ==
		"twos" {
		return MetricQueryParameter.twos
	} else if metricType ==
		"freethrows" {
		return MetricQueryParameter.freethrows
	} else if metricType ==
		"dunks" {
		return MetricQueryParameter.dunks
	} else if metricType == "blocks" {
		return MetricQueryParameter.blocks
	}
	return MetricQueryParameter.score
}
