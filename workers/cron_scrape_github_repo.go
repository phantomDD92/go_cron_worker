package workers

import (
	"fmt"
	"go_proxy_worker/logger"
	"go_proxy_worker/models"
	"go_proxy_worker/slack"
	"go_proxy_worker/utils"
	"log"
	"strings"

	"github.com/go-co-op/gocron"

	"time"
)

func Github_SearchRepos(query string) {
	fileName := "cron_scrape_github_repo.go"
	emptyErrMap := make(map[string]interface{})
	var searchResult utils.Github_SearchResult
	var updateCount, createCount int
	var newRepos []string
	pageMax := 100
	page := 1
	for page = 1; page <= pageMax; page += 1 {
		log.Println("### : Page ", page)
		if err := utils.Github_GetSearchResult(query, page, &searchResult); err != nil {
			logger.LogError("ERROR", fileName, err, "Error scrape github search", emptyErrMap)
			break
		}
		pageMax = searchResult.Payload.PageCount
		for _, el := range searchResult.Payload.Results {
			record := models.GithubRepo{}
			info := utils.Github_RepoInfo{}
			if err := utils.Github_GetRepoResult(&el, &info); err != nil {
				log.Println("--- : ", info.Url)
				continue
			}
			if utils.Github_FindRecord(info.Url, &record) {
				record.RepoStars = info.Stars
				record.RepoForks = info.Forks
				record.RepoLastUpdated = info.UpdatedAt
				record.RepoMaintained = info.Maintained
				if err := utils.Github_UpdateRecord(&record); err != nil {
					println("--- : ", info.Url)
					continue
				}
				updateCount += 1
				log.Println("*** : ", info.Url)
			} else {
				if err := utils.Github_ParseByChatGpt(&info); err != nil {
					println("... : ", info.Url)
					continue
				}
				record.Website = info.Website
				record.WebsitePageTypes = info.PageType
				record.RepoName = info.Name
				record.RepoURL = info.Url
				record.RepoSummary = info.Summary
				record.RepoOwner = info.Owner
				record.RepoCreatedDate = info.CreatedAt
				record.RepoLastUpdated = info.UpdatedAt
				record.RepoMaintained = info.Maintained
				record.RepoStars = info.Stars
				record.RepoForks = info.Forks
				record.CodeLanguage = info.Language
				record.CodeLibrary = info.Libraries
				record.CodeLevel = info.CodeLevel
				record.RepoImageURL = fmt.Sprintf("https://github.com/%s.png", info.Owner)
				if err := utils.Github_CreateRecord(&record); err != nil {
					println("--- : ", info.Url)
					continue
				}
				createCount += 1
				newRepos = append(newRepos, info.Url)
				log.Println("+++ : ", info.Url)
			}
		}
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Number of keywords searched : %d\n", createCount+updateCount))
	builder.WriteString(fmt.Sprintf("Number of repos updated : %d\n", updateCount))
	builder.WriteString(fmt.Sprintf("Number of new repos found : %d\n", createCount))
	builder.WriteString("New repos found : \n")
	for _, repo := range newRepos {
		builder.WriteString(repo + "\n")
	}
	headline := fmt.Sprintf("Github Scraping Report for '%s' ", query)
	slack.SlackStatsAlert("#test-slack-notifications", headline, builder.String())
}

func CronScrapeGithubRepo() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Minutes().Do(RunScrapeGithubRepo)
	s.StartBlocking()
}

func RunScrapeGithubRepo() {
	Github_SearchRepos("amazon scraper")
}
