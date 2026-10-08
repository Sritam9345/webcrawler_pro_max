package urlParser



import (
	"fmt"
	"net/http"
	"net/url"
	"github.com/PuerkitoBio/goquery"
)




func Crawl(pageURL string) ([]string, error) {

	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}

	var urls []string

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {

		href, exists := s.Attr("href")
		if !exists {
			return
		}

		
		u, err := baseURL.Parse(href)
		if err != nil {
			return
		}

		
		if u.Scheme != "http" && u.Scheme != "https" {
			return
		}

		urls = append(urls, u.String())
	})

	return urls, nil
}


//all mem-bounded , working fine