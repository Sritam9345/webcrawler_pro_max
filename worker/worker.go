package worker

import (
	"context"
	"fmt"
	"time"
	"webcrawler/bqSelector"
	"webcrawler/fileUploader"
	"webcrawler/urlParser"
	"github.com/redis/go-redis/v9"
)




func Run(domainName string,manager *bqSelector.Manager) {

	fmt.Printf("Statrting %s worker",domainName)

	if manager == nil {
		return
	}


	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",
		DB:       0,
	})

	ctx := context.Background()


	empty:= 0

	for {

		time.Sleep(5 * time.Second)

		if empty == 1 {

			result,err := rdb.BLPop(ctx,10*time.Second,domainName).Result()

			if err == redis.Nil {

				fmt.Printf("Closing %s worker , empty queue",domainName)

				
				
				manager.Mu.Lock()

				limit,_ := rdb.Get(ctx, "queue_limit").Int()

				limit-=1

				delete(manager.Workers,domainName)

				err_write := rdb.Set(ctx,"queue_limit",limit,0).Err()

				if err_write != nil {
				
					fmt.Println("Error Deleting Worker trying again")
					manager.Mu.Unlock()
					continue
		
				}

				manager.Mu.Unlock()

			} else {

				go func(url string){

					urls,err:= urlParser.Crawl(url)

					if err!= nil {
						fmt.Printf("url %s can't be crawled",url)
						return
					}

					fileUploader.UploadURLs(urls)

				}(result[1])

				continue
			}

			


			

		}


		res,err := rdb.LPop(ctx,domainName).Result()

			if err == redis.Nil {
				empty=1
				continue
			}


			go func(url string){

					urls,err:= urlParser.Crawl(url)

					if err!= nil {
						fmt.Printf("url %s can't be crawled",url)
						return
					}

					fileUploader.UploadURLs(urls)

				}(res)

				

	}
	
}




