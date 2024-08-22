package workers

import (
	"fmt"
	"go_proxy_worker/models"
	"go_proxy_worker/slack"
	"go_proxy_worker/utils"
	"log"
	"strings"

	"github.com/go-co-op/gocron"

	"time"
)

func Google_SearchArticles(query string) {
	url := fmt.Sprintf("https://www.google.co.uk/search?q=%s", strings.ReplaceAll(query, " ", "+"))
	page := 1
	createCount := 0
	updateCount := 0
	urlMap := make(map[string]int)
	var newLinks []string
	for ; url != "" && page < 5; page += 1 {
		println("### Page :", page, " : ", url)
		entries, nextUrl, err := utils.Article_ParseGoogleLinks(url, page)
		if err != nil {
			log.Println(err.Error())
			break
		}
		for _, entry := range entries {
			// check already parsed
			record := models.MarketingWebsiteArticle{}
			answer := utils.ChatGPT_Answer{}
			_, ok := urlMap[entry.Url]
			if ok {
				continue
			}
			if utils.Article_FindRecord(entry.Url, &record) {
				if entry.Date != "" && record.ArticleDate.Format("Jan 2, 2006") != entry.Date {
					println("$$$ : ", record.ArticleDate.Format("Jan 2, 2006"), " : ", entry.Date)
					if err := utils.Article_ParseEntry(&entry); err != nil {
						log.Println(err.Error())
						println("---4: ", entry.Url)
						continue
					}
					if err := utils.ChatGPT_GetAnswerByContent(entry.Content, &answer); err != nil {
						log.Println(err.Error())
						println("---5 : ", entry.Url)
						continue
					}
					utils.Article_FillRecordByEntry(&record, &entry, &answer)
					if err := utils.Article_UpdateRecord(&record); err != nil {
						println("---6 : ", entry.Url)
						log.Println(err.Error())
						continue
					}
					println("*** : ", entry.Url)
					updateCount += 1
				} else {
					println("... : ", entry.Url)
				}
			} else {
				if err := utils.Article_ParseEntry(&entry); err != nil {
					log.Println(err.Error())
					println("---1 : ", entry.Url)
					continue
				}
				if err := utils.ChatGPT_GetAnswerByContent(entry.Content, &answer); err != nil {
					log.Println(err.Error())
					println("---2 : ", entry.Url)
					continue
				}
				utils.Article_FillRecordByEntry(&record, &entry, &answer)
				if err := utils.Article_CreateRecord(&record); err != nil {
					log.Println(err.Error())
					println("---3 : ", entry.Url)
					continue
				}
				println("+++ : ", entry.Url)
				newLinks = append(newLinks, entry.Url)
				createCount += 1
			}
			urlMap[entry.Url] = 1
		}
		url = nextUrl
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Number of searched : %d\n", createCount+updateCount))
	builder.WriteString(fmt.Sprintf("Number of articles updated : %d\n", updateCount))
	builder.WriteString(fmt.Sprintf("Number of new articles found : %d\n", createCount))
	builder.WriteString("New articles found : \n")
	for _, repo := range newLinks {
		builder.WriteString(repo + "\n")
	}
	headline := fmt.Sprintf("Article Scraping Report for '%s' ", query)
	slack.SlackStatsAlert("#test-slack-notifications", headline, builder.String())
}

func CronScrapeGoogleArticle() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Minutes().Do(RunScrapeGoogleArticle)
	s.StartBlocking()
}

func RunScrapeGoogleArticle() {
	Google_SearchArticles("how to scrape amazon")
}
