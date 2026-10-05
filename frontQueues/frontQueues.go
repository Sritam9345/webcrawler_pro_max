package frontQueues


import (
	"context"
	"github.com/redis/go-redis/v9"
	"math/rand/v2"
	"fmt"
)

func AddToFrontQueue(urls []string) {
	
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	for _,url := range urls {

		redisKey := fmt.Sprintf("%s_queue",classifier())

		

		err:= rdb.RPush(ctx,redisKey,url).Err()

		if err != nil {
			panic(err)
		}

	}


	

}




func classifier() string {
	
	randomNum := rand.N(100)

	if randomNum >= 75{
		return "high"

	} else if randomNum >=50 {
		return "medium"
	
	} else {
		return "low"
	}

}