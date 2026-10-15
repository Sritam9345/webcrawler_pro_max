package fqSelector

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
	"webcrawler/backQueues"
	"github.com/redis/go-redis/v9"
)

type HashSet struct {
	mp map[string]struct{}
	mu sync.Mutex
}

func NewHashSet() *HashSet{
	a := &HashSet{
		mp: make(map[string]struct{}),
	}
	return a
}


func (set *HashSet) AddItem (item string) {

	set.mu.Lock()
	set.mp[item]=struct{}{}
	set.mu.Unlock()
}

func (set *HashSet) DeleteItem (item string) {
	
	set.mu.Lock()
	delete(set.mp,item)
	set.mu.Unlock()
}

func (set *HashSet) Contains(item string) bool {

	set.mu.Lock()
	_,exists := set.mp[item]
	set.mu.Unlock()

	return exists
}


func Run(){

	fmt.Println("Running the front queue selector")

	rdb := redis.NewClient(&redis.Options{ Addr: "localhost:6379", Password: "", DB: 0, })

	ctx := context.Background()

	high := 0
	low := 0
	med := 0

	totalCrawled := 0
	hitUrls := 0
	missUrls := 0
	discardedUrls := 0


	ticker := time.NewTicker(2 * time.Second)

go func() {
    for range ticker.C {
        metrics(&missUrls,&hitUrls,&totalCrawled,&discardedUrls)
    }
}()




	set := NewHashSet()
	
	for {


		if (high==1 && low==1 && med==1) {

			high = 0
			low = 0
			med = 0


			result, _ := rdb.BLPop(ctx, 0,
        "high_queue",
        "medium_queue",
        "low_queue",
    		).Result()

			

			exists := set.Contains(result[1])

			

			if exists {
				hitUrls++
				continue
			} else {
				missUrls++
   			    totalCrawled++
				
				isFull := handleNewUrl(set,result[1])

				if isFull == false {
				backQueues.Run(rdb,ctx,result[1])
			} else {
				discardedUrls++
			}
			}


			
			


		}

		

		
	
		num := rand.N(100)

		if (num > 40) {
			
			
			
			result , err := readHighQ(rdb,ctx)

			

			if err == redis.Nil {

				result, err = readMedQ(rdb,ctx)
				high = 1
			}

			if err == redis.Nil {
				result, err = readLowQ(rdb,ctx)
				med = 1
			}

			if err == redis.Nil {
				low = 1
			}

		
			if high==1 && low==1 && med==1 {
				continue
			}

			exists := set.Contains(result)

			

			if exists {
				hitUrls++
				continue
			} else {
				missUrls++
    			totalCrawled++
				isFull := handleNewUrl(set,result)

				if isFull == false {
				backQueues.Run(rdb,ctx,result)
			} else {
				discardedUrls++
			}
			}

			


		} else if (num >10) {

			


			result , err := readMedQ(rdb,ctx)

			

			if err == redis.Nil {

				result, err = readHighQ(rdb,ctx)
				med = 1
			}

			if err == redis.Nil {
				result, err = readLowQ(rdb,ctx)
				high = 1
			}

			if err == redis.Nil {
				low = 1
			}

		
			if high==1 && low==1 && med==1 {
				continue
			}

			exists := set.Contains(result)

			

			if exists {
				hitUrls++
				continue
			} else {
				missUrls++
   				 totalCrawled++
				isFull := handleNewUrl(set,result)

				if isFull == false {
				backQueues.Run(rdb,ctx,result)
			} else {
				discardedUrls++
			}
			}

			

			
			

		} else {

			

			result , err := readLowQ(rdb,ctx)


			if err == redis.Nil {

				result, err = readHighQ(rdb,ctx)
				low = 1
			}

			if err == redis.Nil {
				result, err = readMedQ(rdb,ctx)
				high = 1
			}

			if err == redis.Nil {
				med = 1
			}

		
			if high==1 && low==1 && med==1 {
				continue
			}

			exists := set.Contains(result)

			

			if exists {
				hitUrls++
				continue
			} else {
				missUrls++
   				 totalCrawled++
				
				isFull := handleNewUrl(set,result)

				if isFull == false {
				backQueues.Run(rdb,ctx,result)
			} else {
				discardedUrls++
			}
			}
			

		}

		
	}

}


func readHighQ(rdb *redis.Client, ctx context.Context) (string,error){

	result, err := rdb.LPop(ctx, "high_queue").Result()

	if err == redis.Nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return result, nil

}


func readMedQ(rdb *redis.Client, ctx context.Context) (string, error) {

	result, err := rdb.LPop(ctx, "medium_queue").Result()

	if err == redis.Nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return result, nil
}


func readLowQ(rdb *redis.Client, ctx context.Context) (string,error) {

	result, err := rdb.LPop(ctx,"low_queue").Result()

	if err == redis.Nil {
		return "",err
	}

	if err != nil {
		return "",err
	}

	return result,nil
}

func handleNewUrl(set *HashSet,url string) bool {
	
	if len(set.mp) > 1000 {
		return true
	}
				

	set.AddItem(url)

	time.AfterFunc(10*time.Second,func(){
		set.DeleteItem(url)
	})

	return false

}


func metrics(missUrls *int , hitUrls *int,totalCrawled *int,discardedUrls *int){
	hitMissRatio := 0.0
	discardedRate := 0.0

		if *missUrls > 0 {
    hitMissRatio = float64(*hitUrls) / float64(*missUrls)
					}

		if *totalCrawled >0{
			discardedRate = float64(*discardedUrls) / float64(*totalCrawled)
		}

	fmt.Printf("\rTotal Crawled: %d | Hits: %d | Misses: %d | Hit/Miss Ratio: %.2f | Discarded Rate: %.2f ",
    *totalCrawled, *hitUrls, *missUrls, hitMissRatio,discardedRate)
}


//all mem-bounded , working fine