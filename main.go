package main

import (
	"fmt"
	"net/http"

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

func main() {
	resp, err := http.Get("https://www.ebay.co.uk/sch/i.html?_from=R40&_sacat=0&LH_Auction=1&_nkw=spider-man+ps4&_sop=1")
	handleError(err, "Could not fetch response.")
	defer resp.Body.Close()
	fmt.Println("Status code: ", resp.StatusCode)
	body, err := goquery.NewDocumentFromReader(resp.Body)
	handleError(err, "Could not read response body.")

	var articles []article
	body.Find("#ListViewInner > li").Each(func(i int, s *goquery.Selection) {
		linkValue, _ := s.Find("a").Attr("href")
		priceValue := s.Find(".lvprice").Text()
		finishValue, _ := s.Find(".timeleft .timeMs").Attr("timems")
		articles = append(articles, article{link: linkValue, price: priceValue, finish: finishValue})
	})

	fmt.Println("These are the 50 most recent auctions:")
	articlesStringified, _ := fmt.Printf("%v", articles)
	fmt.Println(articlesStringified)
}
