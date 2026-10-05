package main



import (
	"fmt"
	"webcrawler/urlParser"
)



func main() {
	urls ,err := urlParser.Crawl("https://www.google.com/")

	if err!= nil {
		fmt.Printf("%s",err)
	} else {
		fmt.Println(urls)
	}
}