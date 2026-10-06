package prioritizer

import (
	"encoding/csv"
	"fmt"
	"io"
	"webcrawler/frontQueues"
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
			fmt.Println(err)
			continue
		}

		url := record[1]

		

		urls = append(urls,url)
	}

	frontQueues.AddToFrontQueue(urls)

	return urls,nil
}