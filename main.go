package main

import (
	"fmt"

	"github.com/rickschubert/ebay-camel-camel-camel/crawler"
	"github.com/rickschubert/ebay-camel-camel-camel/database"
	"github.com/rickschubert/ebay-camel-camel-camel/time"
)

func main() {
	db := database.New()
	userEmail := db.GetUserEmail("1")
	fmt.Println(userEmail)
	tracking := db.GetTracking("749143c6-0c79-496b-9d71-d7063036c2e1")

	articles := crawler.GetAuctions(tracking.SearchTerm)

	for _, article := range articles {
		priceLowerThanDesiredMaximum := article.Price < tracking.Price
		AuctionEndsSoon := ((article.Finish - time.GetCurrentTime()) < tracking.MaxTime.ToMs())
		if priceLowerThanDesiredMaximum && AuctionEndsSoon {
			fmt.Println("-------------")
			fmt.Println(article.Link)
		}
	}
}
