package backQueues

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"github.com/redis/go-redis/v9"
)


func Run(rdb *redis.Client,ctx context.Context,url string) error {


	match := extractDomainName(url)



	if match == ""{
		full_error := errors.New("Invalid url")
		return full_error
	}


	exists,err := rdb.SIsMember(ctx, "registered_queues", match).Result()
	
	limit,_ := rdb.Get(ctx, "queue_limit").Int()

	if limit >= 10 && exists == false {
		full_error := errors.New("All workers are busy ")
		return full_error
	}

	if exists == false {
		
	    limit+=1
		rdb.SAdd(ctx,"registered_queues",match)
	}
	
	err_write := rdb.Set(ctx,"queue_limit",limit,0).Err()

	if err_write != nil {
		return err_write
	}


	if err !=nil {
		return err
	}

	len, err := rdb.LLen(ctx, match).Result()

	

	if int(len) < 10 {
		
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

var domainRegex = regexp.MustCompile(
	`(?i)(?:^|\.)(google|amazon|reddit|wikipedia|bbc|twitter|linkedin|youtube|github|` +
		`pinterest|medium|stackoverflow|oracle|microsoft|apple|nvidia|intel|samsung|` +
		`sony|tesla|spacex|facebook|instagram|whatsapp|telegram|discord|slack|` +
		`zoom|dropbox|canva|figma|openai|anthropic|perplexity|huggingface)(?:\.|$)`,
)


func extractDomainName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	host := strings.ToLower(u.Hostname())

	
	match := domainRegex.FindStringSubmatch(host)
	if len(match) > 1 {
		return match[1]
	}

	
	if host == "localhost" {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) > 0 {
			match := domainRegex.FindStringSubmatch(
				strings.ToLower(parts[0]),
			)
			if len(match) > 1 {
				return match[1]
			}
		}
	}

	return ""
}


//all mem-bounded , working fine