package worker

import (
	"context"
	"fmt"
	"time"
	"webcrawler/fileUploader"
	"webcrawler/urlParser"
	"github.com/redis/go-redis/v9"
	"webcrawler/schemas"
	"encoding/json"
)




func Run(rdb *redis.Client,domainName string,manager *schemas.Manager) {

	startTime := 0
	count := 0
	isWaiting := 0
	unCrawlable := 0

	ctx := context.Background()

	fmt.Printf("\nStatrting %s worker\n",domainName)

	ticker := time.NewTicker(2 * time.Second)

go func() {
    for range ticker.C {
        metricsUpdate(ctx,rdb,&isWaiting,&startTime,&count,domainName,&unCrawlable)
    }
}()

	if manager == nil {
		
		return
	}


	

	empty:= 0


	defer func(){	
			manager.Mu.Lock()

				limit,_ := rdb.Get(ctx, "queue_limit").Int()

				rdb.SRem(ctx,"registered_queues",domainName)

				limit-=1

				if limit == 0 {
					fmt.Println("All workers are closed ... none left")
				}

				delete(manager.Workers,domainName)

				err_write := rdb.Set(ctx,"queue_limit",limit,0).Err()

				

				if err_write != nil {

					fmt.Println("error for write")
					fmt.Println(err_write)
				
					fmt.Println("Error Deleting Worker/Queue trying again stage3")
		
				}

				manager.Mu.Unlock()
			}()

		defer ticker.Stop()

	for {

		time.Sleep(1 * time.Second)

		

		if empty == 1 {

			isWaiting = 0

			result,err := rdb.BLPop(ctx,10*time.Second,domainName).Result()

			if err == redis.Nil {

				fmt.Printf("Closing %s worker , empty queue\n",domainName)

				return

			} else {

				if startTime == 0 {
				startTime = int(time.Now().Unix())
			}

			count+=1

				go func(url string){

					urls,err:= urlParser.Crawl(url)


					if err!= nil {
						unCrawlable+=1
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

			isWaiting = 0

			if startTime == 0 {
				startTime = int(time.Now().Unix())
			}

			count+=1

			func(url string){


					urls,err:= urlParser.Crawl(url)

					if err!= nil {
						unCrawlable+=1
						return
					}

					fileUploader.UploadURLs(urls)

				}(res)

	}
	
}



//all mem-bounded , working fine

func metricsUpdate(
    ctx context.Context,
    rdb *redis.Client,
    isWaiting *int,
    startTime *int,
    count *int,
    domainName string,
    unCrawlable *int,
) error {
    if *isWaiting == 1 {
        return nil
    }

    currentTime := time.Now().Unix()
    totalCrawlTime := currentTime - int64(*startTime)

    var speed float64
    if totalCrawlTime > 0 {
        speed = float64(*count) / float64(totalCrawlTime)
    }

    metrics := schemas.WorkerMetrics{
        Count:          *count,
        Speed:          speed,
        Uncrawable:     *unCrawlable,
        TotalCrawlTime: totalCrawlTime,
    }

    jsonData, err := json.Marshal(metrics)
    if err != nil {
        return fmt.Errorf("marshal metrics for %s: %w", domainName, err)
    }

    
    if err := rdb.HSet(
        ctx,
        "worker_metrics",
        domainName,
        string(jsonData),
    ).Err(); err != nil {
        return fmt.Errorf("store metrics for %s: %w", domainName, err)
    }


    return nil
}