package prioritizer

import (
	"encoding/csv"
	"io"
	"webcrawler/frontQueues"
	"fmt"
)

func Run(file io.Reader) ([]string,error) {


	reader := csv.NewReader(file)

	urls := make([]string,0,5)
	
	_, err := reader.Read()

	if err != nil {
		return nil,err
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			continue
		}

		url := record[1]

		urls = append(urls,url)
	}

	fmt.Println(urls)

	frontQueues.AddToFrontQueue(urls)

	return urls,nil
}