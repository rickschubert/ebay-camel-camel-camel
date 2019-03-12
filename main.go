package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

const awsDatabaseRegion = "eu-west-2"

var dynamoClient *dynamodb.DynamoDB

func connectToDynamoDB() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsDatabaseRegion)},
	)
	if err != nil {
		panic("Could not initiate new session with Dynamo DB.")
	}
	dynamoClient = dynamodb.New(sess)
	fmt.Println("Established dynamodb session")
}

func setEnvironmentVariablesForDevelopmentPurposes() {
	os.Setenv("SEARCH_TERM", "Spider-Man PS4")
	os.Setenv("PRICE", "20")
	os.Setenv("MAX_TIME", "18000000")
	os.Setenv("USER_ID", "1")
}

func exitProgramOnUndefinedEnv(env string) {
	_, err := fmt.Printf("The environment variable %v needs to be defined.\n", env)
	if err != nil {
		fmt.Println("Had trouble formatting the error message.")
	}
	os.Exit(1)
}

func retrieveEnv(env string) string {
	envRetr, envex := os.LookupEnv(env)
	if !envex {
		exitProgramOnUndefinedEnv(env)
	}
	return envRetr
}

func _handleNumberParsingError(err error, env string, expectedType string) {
	msg := fmt.Sprintf("ERROR: Could not parse environment variable %v into %v\n", env, expectedType)
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func retrieveEnvParsedAsInt(env string) int64 {
	envRetr := retrieveEnv(env)
	intified, err := strconv.ParseInt(envRetr, 10, 64)
	_handleNumberParsingError(err, env, "int")
	return intified
}

func retrieveEnvParsedAsFloat(env string) float64 {
	envRetr := retrieveEnv(env)
	floatified, err := strconv.ParseFloat(envRetr, 64)
	_handleNumberParsingError(err, env, "float")
	return floatified
}

func main() {
	setEnvironmentVariablesForDevelopmentPurposes()

	connectToDynamoDB()
	searchTerm := retrieveEnv("SEARCH_TERM")
	userId := retrieveEnv("USER_ID")
	fmt.Println(userId)
	maxPrice := retrieveEnvParsedAsFloat("PRICE")
	maxTimeLeft := retrieveEnvParsedAsInt("MAX_TIME")

	articles := getAuctions(searchTerm)

	for _, article := range articles {
		if article.price < maxPrice && article.finish-getCurrentTime() < maxTimeLeft {
			fmt.Println("-------------")
			fmt.Println(article.link)
		}
	}

}
