package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ericchiang/css"
	"golang.org/x/net/html"
)

func handleError(err error, preExitMsg string) {
	if err != nil {
		fmt.Println(preExitMsg)
		fmt.Println(err)
		panic(err)
	}
}

func retrieveAuctionLink(article *html.Node) string {
	productLinkSelector := getSelector(".lvtitle > a")
	link := productLinkSelector.Select(article)
	for _, attribute := range link[0].Attr {
		if attribute.Key == "href" {
			return attribute.Val
		}
	}
	return ""
}

func retrievePrice(article *html.Node) string {
	priceSelector := getSelector("span.bold")
	priceSpan := priceSelector.Select(article)
	if len(priceSpan) == 0 {
		fmt.Println("What the fuck, this is empty?! QARS")
	} else {
		fmt.Println("QARS")
		fmt.Println(priceSpan[0].Data)
	}
	// for _, nodeText := range priceSpan[0].Data {
	// }
	return ""
}

func getSelector(selector string) *css.Selector {
	selectorElement, err := css.Compile(selector)
	handleError(err, "Selector creation failed.")
	return selectorElement
}

type article struct {
	link  string
	price string
}

func main() {
	fmt.Println("Program executed")
	resp, err := http.Get("https://www.ebay.co.uk/sch/i.html?_from=R40&_sacat=0&LH_Auction=1&_nkw=spider-man+ps4&_sop=1")
	defer resp.Body.Close()
	handleError(err, "Could not fetch response.")
	fmt.Println("Status code: ", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	handleError(err, "Could not read response body.")

	articleContainerSelector := getSelector("#ListViewInner > li")

	document, err := html.Parse(strings.NewReader(string(body)))
	handleError(err, "Unable to parse HTML document.")

	var articles []article
	for _, articleNode := range articleContainerSelector.Select(document) {
		fmt.Println("\n\nTHiS is an article \n")
		html.Render(os.Stdout, articleNode)
		individualArticle := article{link: retrieveAuctionLink(articleNode), price: retrievePrice(articleNode)}
		articles = append(articles, individualArticle)
	}

	fmt.Println("Below now the final result")
	s, _ := fmt.Printf("%v", articles)
	fmt.Println(s)
}
