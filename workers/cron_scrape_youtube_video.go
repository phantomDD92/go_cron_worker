package workers

import (
	"fmt"
	"go_proxy_worker/models"
	"go_proxy_worker/slack"
	"go_proxy_worker/utils"
	"log"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

func CronScrapeYoutubeVideo() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Minutes().Do(RunScrapeYoutubeVideo)
	s.StartBlocking()
}

func RunScrapeYoutubeVideo() {
	YouTube_StartSearch("how to scrape amazon")
}

func YouTube_StartSearch(query string) {
	url := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=vid", strings.ReplaceAll(query, " ", "+"))
	page := 1
	createCount := 0
	updateCount := 0
	var newLinks []string
	for ; url != "" && page < 10; page += 1 {
		println("### Page :", page, " : ", url)
		links, nextUrl, err := utils.YouTube_ParseYoutubeLinks(url, page)
		if err != nil {
			log.Println(err.Error())
			break
		}
		for _, link := range links {
			info := utils.YouTubeVideo_Info{}
			answer := utils.ChatGPT_Answer{}
			record := models.MarketingWebsiteVideo{}
			if utils.YouTube_FindRecord(link, &record) {
				if err := utils.YouTube_ParseVideo(link, &info, true); err != nil {
					println("--- : ", link)
					continue
				}
				record.VideoLikes = info.Likes
				record.VideoViews = info.Views
				if err := utils.YouTube_UpdateRecord(&record); err != nil {
					println("--- : ", link)
					continue
				}
				updateCount += 1
				println("*** : ", link)
			} else {
				if err := utils.YouTube_ParseVideo(link, &info, false); err != nil {
					println("--- : ", link)
					continue
				}
				if err := utils.ChatGPT_GetAnswerForYoutube(info.Desc, info.Transcript, &answer); err != nil {
					println("... : ", link)
					continue
				}
				record.VideoChannel = info.Channel
				record.VideoChannelURL = info.ChannelURL
				record.VideoDate = info.Date
				record.VideoLength = info.Duration
				record.VideoLikes = info.Likes
				record.VideoName = info.Name
				record.VideoURL = link
				record.VideoViews = info.Views
				record.VideoPreviewImageUrl = info.Thumbnail
				record.CodeGithubRepo = info.Repo
				record.CodeLanguage = answer.Language
				record.CodeLevel = answer.CodeLevel
				record.CodeLibraries = answer.Libraries
				record.WebsitePageTypes = answer.PageType
				record.VideoSummary = answer.Summary
				record.Website = answer.WebSite
				if err := utils.YouTube_CreateRecord(&record); err != nil {
					println("--- : ", link)
					continue
				}
				newLinks = append(newLinks, link)
				println("+++ : ", link)
				createCount += 1
			}
		}
		url = nextUrl
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Number of searched : %d\n", createCount+updateCount))
	builder.WriteString(fmt.Sprintf("Number of videos updated : %d\n", updateCount))
	builder.WriteString(fmt.Sprintf("Number of new videos found : %d\n", createCount))
	builder.WriteString("New video found : \n")
	for _, repo := range newLinks {
		builder.WriteString(repo + "\n")
	}
	headline := fmt.Sprintf("* Search report for '%s' *", query)
	slack.SlackStatsAlert("#test-slack-notifications", headline, builder.String())
}
