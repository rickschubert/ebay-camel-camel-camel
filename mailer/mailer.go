package mailer

import (
	"fmt"
	"os"
	"strings"

	"github.com/rickschubert/ebay-camel-camel-camel/crawler"
	gomail "gopkg.in/gomail.v2"
)

func prepareMessage(emailAddress string, filteredArticles []crawler.Article) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", "update@ebaycamelcamelcamel.com")
	m.SetHeader("To", emailAddress)
	m.SetHeader("Subject", "Ebay camel camel camel found new articles you want")
	var mailBody strings.Builder
	mailBody.WriteString(fmt.Sprintf("<p>There are new articles you want!</p>"))
	for _, article := range filteredArticles {
		mailBody.WriteString(fmt.Sprintf("<p>%v<br>Current bid: %v<br>Time remaining for auction: %v</p>", article.Link, article.Price, article.Finish.Readable()))
	}
	fmt.Println(mailBody.String())
	m.SetBody("text/html", mailBody.String())
	return m
}

func NotifyUsersOfNewAuctions(emailAddress string, filteredArticles []crawler.Article) {
	m := prepareMessage(emailAddress, filteredArticles)
	fmt.Println(m)
	mailUser, mailUserEx := os.LookupEnv("MAIL_USER")
	if !mailUserEx {
		panic("The environment variable MAIL_USER is not defined.")
	}
	mailPassword, mailPasswordEx := os.LookupEnv("MAIL_PASSWORD")
	if !mailPasswordEx {
		panic("The environment variable MAIL_PASSWORD is not defined.")
	}
	d := gomail.NewDialer("smtp.gmail.com", 465, mailUser, mailPassword)
	// d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
