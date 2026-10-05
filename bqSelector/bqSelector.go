package bqSelector

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"webcrawler/schemas"
	"webcrawler/worker"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	Mu      sync.Mutex
	Workers map[string]any
}


func Run() error {

	fmt.Println("hi there")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	manager := &schemas.Manager{
		Workers: make(map[string]any),
	}

	

	for {

	array,_ := rdb.BLPop(ctx,0,
	"waiting_list",
	).Result()

	domainName := array[1]

	_,ok := manager.Workers[domainName]

	if ok == true {
		continue
	}

	limit,_ := rdb.Get(ctx, "queue_limit").Int()

	if limit >= 10{
		full_error := errors.New("Queue is full")
		fmt.Println(full_error)
		continue
	}

	ctx,cancel := context.WithCancel(context.Background())

	manager.Workers[domainName] = []any{ctx,cancel}

	fmt.Println("Running the worker")

	go worker.Run(domainName,manager)

}

	
}

