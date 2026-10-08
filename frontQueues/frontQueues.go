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

	for _, url := range urls {

	redisKey := fmt.Sprintf("%s_queue", classifier())

	script := `
		local len = redis.call("LLEN", KEYS[1])

		if len >= 100 then
			return 0
		end

		redis.call("RPUSH", KEYS[1], ARGV[1])
		return 1
	`

	result, err := rdb.Eval(
		ctx,
		script,
		[]string{redisKey},
		url,
	).Int()

	if err != nil {
		panic(err)
	}

	if result == 0 {
		fmt.Println("Queue is full")
		continue
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


//all mem-bounded , working fine