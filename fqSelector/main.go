package main


import (
	"context"
	"github.com/redis/go-redis/v9"
	"math/rand/v2"
	"webcrawler/backQueues"
	"fmt"
	"time"
)

type HashSet map[string]struct{}

func NewHashSet() *HashSet{
	a := make(HashSet)
	return &a
}


func (set HashSet) AddItem (item string){
	set[item] = struct{}{}
}

func (set HashSet) DeleteItem (item string){
	delete(set,item)
}

func (set HashSet) Contains(item string) bool{
	_,exists := set[item]
	return exists
}


func fQselector(){

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

			
			fmt.Println(result[1])

			exists := set.Contains(result[1])

			if exists {
				continue
			} else {
				
				go handleNewUrl(set,result[1])
			}


			backQueues.Run(rdb,ctx,result[1])
			
			


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

			backQueues.Run(rdb,ctx,result)

			


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

			backQueues.Run(rdb,ctx,result)

			
			

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

			backQueues.Run(rdb,ctx,result)

			

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

func handleNewUrl(set *HashSet,url string) {

	set.AddItem(url)

	time.AfterFunc(2*time.Second,func(){
		set.DeleteItem(url)
	})

}


func main(){
	fQselector()
}