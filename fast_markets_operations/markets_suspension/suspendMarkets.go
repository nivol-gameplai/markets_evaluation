package markets_suspension

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/api/iterator"
	"log"
	"marketsEvaluation/fast_markets_operations/configs"
	"sync"
	"time"
)

var redisPool *redis.Pool

// DetermineSuspensionStrategy this function accepts tries to enforce the rules of the markets suspension strategy
// calls MarketSuspension to be enforced based on the strategy
// INPUTS:
//   strategy string i.e. next_bucket, team that has the ball possession (soccer, nba) string i.e. home/away string i.e.
//   and home/away team ids and gameid as string. Finally, current metric as an int
//  (could be score, assists and other metrics related to markets.)
// everything Reads and Writes from Firestore are performed under a transaction to make sure clients are not reading stale
// data
func DetermineSuspensionStrategy(strategy string, currentTeamPossession string, homeTeamId string, gameId string,
	awayTeamId string, metric int64) {
	if strategy == "next_bucket" {
		var previousTeamPossession string
		previousTeamPossession = ""
		var teamMarketsToSuspend string
		if redisPool == nil {
			// Pre-declare err to avoid shadowing redisPool
			var err error
			redisPool, err = configs.InitializeRedis()
			if err != nil {
				return
			}
		}
		conn := redisPool.Get()
		defer func(conn redis.Conn) {
			err := conn.Close()
			if err != nil {
				return
			}
		}(conn)
		exists, err := redis.Int(conn.Do("EXISTS", gameId+"_possession"))
		if err != nil {
			return
		} else if exists != 0 {
			previousTeamPossession, err = redis.String(conn.Do("GET", gameId+"_possession"))
			if err != nil {
				return
			}
		}
		if currentTeamPossession != previousTeamPossession {
			if currentTeamPossession == homeTeamId {
				teamMarketsToSuspend = awayTeamId
			} else {
				teamMarketsToSuspend = homeTeamId
			}
			_, err = MarketSuspension(gameId, teamMarketsToSuspend, metric)
			if err != nil {
				log.Fatal(err)
				return
			}
			_, err = conn.Do("SET", gameId+"_possession", currentTeamPossession, "EX", "86400")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// MarketSuspension executes odds queries based on parameters to suspend
// specific markets based on metric (will determine the markets), the teamid (will only need to suspend markets of specific
// team or if teamId ="NA" it will suspend all active markets
// calls updateSuspensionOnDocs as a go routine (threads) to update all market docs on parallel
func MarketSuspension(gameId string, teamId string, metric int64) ([]byte, error) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer configs.CloseClient(client)
	var query firestore.Query
	col := client.Collection("Odds")
	var wg sync.WaitGroup

	if metric < 0 && teamId != "NA" {

		query = col.Where("game_id", "==", gameId).Where("is_active", "==",
			true).Where("team_id", "==", teamId)
	} else if metric >= 0 && teamId != "NA" {
		query = col.Where("game_id", "==", gameId).Where("is_active", "==",
			true).Where("team_id", "==", teamId).Where("current_team_score", "<=", metric)
	} else {
		query = col.Where("game_id", "==", gameId).Where("is_active", "==",
			true)
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

			go updateSuspensionOnDocs(doc.Data(), metric, tx, *doc.Ref, &wg)

		}
		wg.Wait()
		return nil
	})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil, nil
}

//updateSuspensionOnDocs peferms the actual suspension (freeze) on the markets and saves them back to the
// irrelevant of the amount of the docs (could be 10s of thousdands) the transactions are completed
// in parallel in milliseconds
func updateSuspensionOnDocs(odd map[string]interface{}, metric int64,
	tx *firestore.Transaction, ref firestore.DocumentRef, wg *sync.WaitGroup) {
	defer wg.Done()
	if odd["freeze"] == true {

		odd["freeze"] = false
	} else {
		if metric < 0 {
			odd["freeze"] = true
		}
	}
	if metric >= 0 {
		odd["is_active"] = false
		odd["freeze"] = false
	}
	odd["modified"] = time.Now().UnixNano() / int64(time.Millisecond)
	err := tx.Set(&ref, odd)
	if err != nil {
		return
	}

}
