package markets_suspension

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/api/iterator"
	"log"
	"marketsEvaluation/fast_markets_operations/configs"
	"sync"
)

var redisPool *redis.Pool

func DetermineSuspensionStrategy(strategy string, currentTeamPossession string, homeTeamId string, gameId string,
	awayTeamId string, score int64) {
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
			_, err = MarketSuspension(gameId, teamMarketsToSuspend, score)
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

func MarketSuspension(gameId string, teamId string, currentScore int64) ([]byte, error) {
	ctx := context.Background()
	client := configs.CreateClient(ctx)
	defer configs.CloseClient(client)
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
