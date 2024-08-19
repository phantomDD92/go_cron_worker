package workers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_proxy_worker/db"
	"go_proxy_worker/logger"
	"go_proxy_worker/models"
	"go_proxy_worker/slack"
	"go_proxy_worker/utils"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-co-op/gocron"

	"time"
)

type GithubRepoInfo struct {
	Id          string `json:"id"`
	Followers   int    `json:"followers"`
	Language    string `json:"language"`
	Description string `json:"hl_trunc_description"`
	Repo        struct {
		Repository struct {
			Id         uint64 `json:"id"`
			OwnerId    int    `json:"owner_id"`
			Name       string `json:"name"`
			OwnerLogin string `json:"owner_login"`
			UpdatedAt  string `json:"updated_at"`
		} `json:"repository"`
	} `json:"repo"`
	Topics []string `json:"topics"`
}

type GithubSearchResult struct {
	Payload struct {
		Results   []GithubRepoInfo `json:"results"`
		PageCount int              `json:"page_count"`
	} `json:"payload"`
}

type GithubRepoDetail struct {
	Props struct {
		InitialPayload struct {
			Repo struct {
				CreatedAt string `json:"createdAt"`
			} `json:"repo"`
		} `json:"initialPayload"`
	} `json:"props"`
}

func Github_GetRepoResult(info *GithubRepoInfo, repo *models.GithubRepo) error {
	var createdAt, updatedAt time.Time
	url := "https://github.com/" + info.Repo.Repository.OwnerLogin + "/" + info.Repo.Repository.Name
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	forkButton := doc.Find("a#fork-button").First()
	if fork, err := utils.ParseIntFromPattern(forkButton.Text(), `([\d\.]+)`); err == nil {
		repo.RepoForks = fork
	}
	if scriptTag := doc.Find("react-partial[partial-name='repos-overview'] > script"); scriptTag.Length() > 0 {
		var detail GithubRepoDetail
		if err := json.Unmarshal([]byte(scriptTag.Text()), &detail); err != nil {
			return err
		}
		createdAt, err = time.Parse("2006-01-02T15:04:05.000Z", detail.Props.InitialPayload.Repo.CreatedAt)
		if err != nil {
			return err
		}
	}
	updatedAt, err = time.Parse("2006-01-02T15:04:05.000Z", info.Repo.Repository.UpdatedAt)
	if err != nil {
		return err
	}
	duration := time.Since(updatedAt)
	repo.Website = "https://github.com"
	repo.RepoURL = url
	repo.RepoName = info.Repo.Repository.Name
	repo.RepoSummary = info.Description
	repo.RepoOwner = info.Repo.Repository.OwnerLogin
	repo.RepoCreatedDate = createdAt
	repo.RepoLastUpdated = updatedAt
	repo.RepoMaintained = duration.Hours() < 24*30*60
	repo.RepoStars = info.Followers
	repo.CodeLanguage = info.Language
	repo.CodeLibrary = strings.Join(info.Topics, ",")
	return nil
}

func Github_GetSearchResult(query string, page int, info *GithubSearchResult) error {
	url := "https://github.com/search?q=" + strings.ReplaceAll(query, " ", "+") + "&type=repositories&p=" + strconv.Itoa(page)
	body, err := utils.GetHtmlByProxy(url)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return err
	}
	scriptTag := doc.Find("script[data-target='react-app.embeddedData']").First()
	if scriptTag.Length() == 0 {
		return errors.New("data not found")
	}
	err = json.Unmarshal([]byte(scriptTag.Text()), info)
	return err
}

func Github_UpdateDB(repo *models.GithubRepo) (updated bool, err error) {
	var record models.GithubRepo
	db := db.GetDB()
	result := db.Where("repo_url = ?", repo.RepoURL).First(&record)
	if result.Error != nil {
		result := db.Create(&repo)
		return true, result.Error
	}
	record.RepoStars = repo.RepoStars
	record.RepoForks = repo.RepoForks
	record.RepoLastUpdated = repo.RepoLastUpdated
	record.RepoMaintained = repo.RepoMaintained
	result = db.Save(&record)
	return false, result.Error
}

func Github_SearchRepos(query string) {
	fileName := "cron_scrape_github_repo.go"
	emptyErrMap := make(map[string]interface{})
	var searchResult GithubSearchResult
	var searchCount, createCount, errorCount int
	var newRepos []string
	pageMax := 100
	page := 1
	for page = 1; page <= pageMax; page += 1 {
		// log.Println("### : Page ", page)
		if err := Github_GetSearchResult(query, page, &searchResult); err != nil {
			logger.LogError("ERROR", fileName, err, "Error scrape github search", emptyErrMap)
			break
		}
		pageMax = searchResult.Payload.PageCount
		for _, el := range searchResult.Payload.Results {
			searchCount += 1
			gitRepo := models.GithubRepo{}
			if err := Github_GetRepoResult(&el, &gitRepo); err != nil {
				errorCount += 1
				logger.LogError("ERROR", fileName, err, "Error scrape github repository", emptyErrMap)
				continue
			}
			created, err := Github_UpdateDB(&gitRepo)
			if err != nil {
				logger.LogError("ERROR", fileName, err, "Error updating database", emptyErrMap)
				errorCount += 1
				continue
			}
			if created {
				// log.Println("+++ : ", gitRepo.RepoURL)
				createCount += 1
				newRepos = append(newRepos, gitRepo.RepoURL)
			} else {
				// log.Println("--- : ", gitRepo.RepoURL)
			}
		}
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Number of keywords searched : %d\n", searchCount))
	builder.WriteString(fmt.Sprintf("Number of repos updated : %d\n", searchCount-createCount-errorCount))
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
