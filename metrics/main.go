
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"webcrawler/schemas"
)

func runMetricsDisplay(
	ctx context.Context,
	rdb *redis.Client,
) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if err := renderMetrics(ctx, rdb); err != nil {
			fmt.Printf("\033[HMetrics error: %v\033[J", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func renderMetrics(
	ctx context.Context,
	rdb *redis.Client,
) error {
	// Get the currently registered workers.
	domains, err := rdb.SMembers(
		ctx, "registered_queues",
	).Result()
	if err != nil {
		return err
	}

	// Fetch all worker metrics in one Redis command.
	metricsMap, err := rdb.HGetAll(
		ctx, "worker_metrics",
	).Result()
	if err != nil {
		return err
	}

	// Keep display order consistent across refreshes.
	sort.Strings(domains)

	var output strings.Builder

	// Move cursor to the top of the terminal.
	output.WriteString("\033[H")

	if len(domains) == 0 {
		output.WriteString("No active workers\n")
	} else {
		for _, domainName := range domains {
			jsonData, exists := metricsMap[domainName]

			if !exists {
				fmt.Fprintf(
					&output,
					"Worker: %-12s | Metrics pending...\033[K\n",
					domainName,
				)
				continue
			}

			var metrics schemas.WorkerMetrics

			if err := json.Unmarshal(
				[]byte(jsonData), &metrics,
			); err != nil {
				fmt.Fprintf(
					&output,
					"Worker: %-12s | Invalid metrics JSON\033[K\n",
					domainName,
				)
				continue
			}

			fmt.Fprintf(
				&output,
				"Worker: %-12s | URLs: %-6d | Time: %-6ds | Speed: %.2f URLs/sec | Uncrawled: %-6d\033[K\n",
				domainName,
				metrics.Count,
				metrics.TotalCrawlTime,
				metrics.Speed,
				metrics.Uncrawable,
			)
		}
	}

	// Clear leftover lines from the previous display.
	output.WriteString("\033[J")

	fmt.Print(output.String())
	return nil
}


func main(){

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	runMetricsDisplay(ctx,rdb)
}