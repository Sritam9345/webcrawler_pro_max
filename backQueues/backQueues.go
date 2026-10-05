package backQueues

import (
	"context"
	"github.com/redis/go-redis/v9"
	"errors"
	"strings"
	"net/url"
	"fmt"

)


func Run(rdb *redis.Client,ctx context.Context,url string) error {


	match := extractDomainName(url)


	exists,err:= rdb.LLen(ctx,match).Result()

	

	limit,_ := rdb.Get(ctx, "queue_limit").Int()
	fmt.Println(limit)

	if limit >= 10 {
		full_error := errors.New("All workers are busy ")
		return full_error
	}

	limit+=1
	


	err_write := rdb.Set(ctx,"queue_limit",limit,0).Err()

	if err_write != nil {
		return err_write
	}


	if err !=nil {
		return err
	}

	if exists < 10 {
		fmt.Println("reached here")
		err_push := rdb.LPush(ctx,match,url).Err()

		rdb.LPush(ctx,"waiting_list",match)
		
		if err_push!= nil{
			return err_push
		}

	} else{

		err = errors.New("Queue is full")
		return err

	}

	return nil

}



func extractDomainName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	host := u.Hostname()

	
	host = strings.TrimPrefix(host, "www.")

	
	parts := strings.Split(host, ".")

	if len(parts) == 0 {
		return ""
	}

	return parts[0]
}