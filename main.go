package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
	finish milliSeconds
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

		articles = append(articles, Article{
			link:   linkValue,
			price:  pricePlainNumber,
			finish: milliSeconds(finishValueInt),
		})
	})
	return articles
}

func getCurrentTime() milliSeconds {
	return milliSeconds(time.Now().UnixNano() / int64(time.Millisecond))
}

func connectToDynamoDB() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsDatabaseRegion)},
	))
	dynamoClient = dynamodb.New(sess)
	fmt.Println("Established dynamodb session")
}

type user struct {
	UUID  string `json:"UUID"`
	Email string `json:"email"`
}

func getUserEmailFromDatabase(userId string) string {
	result, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"UUID": {
				S: aws.String(userId),
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Could not retrieve the user with UUID %v", userId))
	}

	var userRetrieved user
	errMarsh := dynamodbattribute.UnmarshalMap(result.Item, &userRetrieved)
	if errMarsh != nil {
		panic(fmt.Sprintf("Failed to unmarshal record %v", err))
	}

	return userRetrieved.Email
}

type trackingInformation struct {
	SearchTerm string  `json:"searchTerm"`
	UserId     string  `json:"userId"`
	Price      float64 `json:"price"`
	MaxTime    minutes `json:"maxTime"`
}

func getTrackingFromDatabase(trackingId string) trackingInformation {
	result, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("trackings"),
		Key: map[string]*dynamodb.AttributeValue{
			"UUID": {
				S: aws.String(trackingId),
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Could not retrieve the tracking with UUID %v", trackingId))
	}

	var trackingRetrieved trackingInformation
	err = dynamodbattribute.UnmarshalMap(result.Item, &trackingRetrieved)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal record %v", err))
	}

	return trackingRetrieved
}

var dynamoClient *dynamodb.DynamoDB

const awsDatabaseRegion = "eu-west-2"

type minutes int64

type milliSeconds int64

func (m minutes) toMs() milliSeconds {
	return milliSeconds(m * 60000)
}

func main() {
	connectToDynamoDB()
	userEmail := getUserEmailFromDatabase("1")
	fmt.Println(userEmail)
	tracking := getTrackingFromDatabase("895ba9b7-df51-4f58-a73d-0f6c57294ad6")

	articles := getAuctions(tracking.SearchTerm)

	for _, article := range articles {
		priceLowerThanDesiredMaximum := article.price < tracking.Price
		timeInDesiredNotificationTimeFrame := ((article.finish - getCurrentTime()) < tracking.MaxTime.toMs())
		if priceLowerThanDesiredMaximum && timeInDesiredNotificationTimeFrame {
			fmt.Println("-------------")
			fmt.Println(article.link)
		}
	}

}
