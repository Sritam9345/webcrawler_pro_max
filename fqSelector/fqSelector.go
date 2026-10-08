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
	fmt.Printf("added %s to map\n",item)

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

	set := NewHashSet()
	
	for {

		if (high==1 && low==1 && med==1) {

			high = 0
			low = 0
			med = 0


			fmt.Println("all frontQueues are empty... waiting...")

			result, _ := rdb.BLPop(ctx, 0,
        "high_queue",
        "medium_queue",
        "low_queue",
    		).Result()

			

			exists := set.Contains(result[1])

			if exists {
				fmt.Printf("already crawled..%s\n",result[1])
				continue
			} else {
				
				go handleNewUrl(rdb,ctx,set,result[1])
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
				fmt.Printf("already crawled..%s\n",result)
				continue
			} else {
				
				go handleNewUrl(rdb,ctx,set,result)
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
				fmt.Printf("already crawled..%s\n",result)
				continue
			} else {
				
				go handleNewUrl(rdb,ctx,set,result)
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
				fmt.Printf("already crawled..%s\n",result)
				continue
			} else {
				
				go handleNewUrl(rdb,ctx,set,result)
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

func handleNewUrl(rdb *redis.Client,ctx context.Context,set *HashSet,url string) {
	
	if len(set.mp) > 100 {
		fmt.Println("hash-map is full can't add more")
		return
	}


	set.AddItem(url)

	time.AfterFunc(3600*time.Second,func(){
		set.DeleteItem(url)
	})

	backQueues.Run(rdb,ctx,url)

}




//all mem-bounded , working fine