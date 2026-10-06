package worker

import (
	"context"
	"fmt"
	"time"
	"webcrawler/fileUploader"
	"webcrawler/urlParser"
	"github.com/redis/go-redis/v9"
	"webcrawler/schemas"
)




func Run(domainName string,manager *schemas.Manager) {

	

	fmt.Printf("Statrting %s worker\n",domainName)

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


	defer func(){	manager.Mu.Lock()

				limit,_ := rdb.Get(ctx, "queue_limit").Int()

				limit-=1

				if limit == 0 {
					fmt.Println("All workers are closed ... none left")
				}

				delete(manager.Workers,domainName)

				err_write := rdb.Set(ctx,"queue_limit",limit,0).Err()

				

				if err_write != nil {

					fmt.Println("error for write")
					fmt.Println(err_write)
				
					fmt.Println("Error Deleting Worker/Queue trying again")
		
				}

				manager.Mu.Unlock()
			}()

	for {

		time.Sleep(5 * time.Second)

		if empty == 1 {

			result,err := rdb.BLPop(ctx,10*time.Second,domainName).Result()

			if err == redis.Nil {

				fmt.Printf("Closing %s worker , empty queue\n",domainName)

				return

			} else {

				go func(url string){

					urls,err:= urlParser.Crawl(url)


					if err!= nil {
						fmt.Printf("url %s can't be crawled\n",url)
						return
					} else {
						fmt.Printf("successfully crawled %s\n",url)
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
						fmt.Printf("url %s can't be crawled\n",url)
						return
					} else {
						fmt.Printf("successfully crawled %s\n",url)
					}

					fileUploader.UploadURLs(urls)

				}(res)

	}
	
}




