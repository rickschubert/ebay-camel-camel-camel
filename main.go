package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func handleError(err error, preExitMsg string) {
	if err != nil {
		fmt.Println(preExitMsg)
		fmt.Println(err)
		panic(err)
	}
}

type article struct {
	link   string
	price  string
	finish string
}

const searchTerm = "Spider-Man PS4"

func getAuctions(searchTerm string) []article {
	url := fmt.Sprintf("https://www.ebay.co.uk/sch/i.html?_from=R40&_sacat=0&LH_Auction=1&_nkw=%v&_sop=1", url.QueryEscape(searchTerm))
	fmt.Println(url)
	return crawl(url)
}

func crawl(url string) []article {
	resp, err := http.Get(url)
	handleError(err, "Could not fetch response.")
	defer resp.Body.Close()
	fmt.Println("Status code: ", resp.StatusCode)
	body, err := goquery.NewDocumentFromReader(resp.Body)
	handleError(err, "Could not read response body.")

	selectors := map[string]string{
		"articleContainer": "#ListViewInner > li",
		"price":            ".lvprice",
		"finishTime":       ".timeleft .timeMs",
	}

	var articles []article
	body.Find(selectors["articleContainer"]).Each(func(i int, s *goquery.Selection) {
		linkValue, _ := s.Find("a").Attr("href")
		priceValue := s.Find(selectors["price"]).Text()
		finishValue, _ := s.Find(selectors["finishTime"]).Attr("timems")
		articles = append(articles, article{link: linkValue, price: priceValue, finish: finishValue})
	})
	return articles
}

func main() {
	articles := getAuctions(searchTerm)
	fmt.Println("These are the 50 most recent auctions:")
	articlesStringified, _ := fmt.Printf("%v", articles)
	fmt.Println(articlesStringified)
}
