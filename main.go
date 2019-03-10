package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func handleError(err error, preExitMsg string) {
	if err != nil {
		fmt.Println(preExitMsg)
		fmt.Println(err)
		panic(err)
	}
}

type Article struct {
	link   string
	price  float64
	finish int64
}

func getAuctions(searchTerm string) []Article {
	url := fmt.Sprintf("https://www.ebay.co.uk/sch/i.html?_from=R40&_sacat=0&LH_Auction=1&_nkw=%v&_sop=1", url.QueryEscape(searchTerm))
	fmt.Println(url)
	return crawl(url)
}

func crawl(url string) []Article {
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

	var articles []Article
	body.Find(selectors["articleContainer"]).Each(func(i int, s *goquery.Selection) {
		linkValue, _ := s.Find("a").Attr("href")

		priceValue := s.Find(selectors["price"]).Text()
		priceValueTrimmed := strings.TrimSpace(priceValue)
		priceValuePlain := strings.Replace(priceValueTrimmed, "£", "", 1)
		pricePlainNumber, _ := strconv.ParseFloat(priceValuePlain, 64)

		finishValue, _ := s.Find(selectors["finishTime"]).Attr("timems")
		finishValueInt, _ := strconv.ParseInt(finishValue, 0, 64)

		articles = append(articles, Article{link: linkValue, price: pricePlainNumber, finish: finishValueInt})
	})
	return articles
}

func getCurrentTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {
	const searchTerm = "Spider-Man PS4"
	const maxPrixe = 20
	const maxTimeLeft = 300 * 60000 // X minutes in milliseconds

	articles := getAuctions(searchTerm)

	for _, article := range articles {
		if article.price < maxPrixe && article.finish-getCurrentTime() < maxTimeLeft {
			fmt.Println("-------------")
			fmt.Println(article.link)
		}
	}
}
