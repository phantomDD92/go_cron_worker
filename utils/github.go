package utils

import (
	"encoding/json"
	"errors"
	"go_proxy_worker/db"
	"go_proxy_worker/models"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Github_RepoInfo struct {
	Url        string    `json:"url"`
	Name       string    `json:"name"`
	Summary    string    `json:"summary"`
	Owner      string    `json:"owner"`
	Language   string    `json:"language"`
	Libraries  string    `json:"libraires"`
	Stars      int       `json:"stars"`
	Forks      int       `json:"forks"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Maintained bool      `json:"maintained"`
	Readme     string    `json:"readme"`
}

type Github_RepoEntry struct {
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

type Github_SearchResult struct {
	Payload struct {
		Results   []Github_RepoEntry `json:"results"`
		PageCount int                `json:"page_count"`
	} `json:"payload"`
}

type Github_RepoDetail struct {
	Props struct {
		InitialPayload struct {
			Repo struct {
				CreatedAt string `json:"createdAt"`
			} `json:"repo"`
		} `json:"initialPayload"`
	} `json:"props"`
}

func Github_GetRepoResult(entry *Github_RepoEntry, info *Github_RepoInfo) error {
	var createdAt, updatedAt time.Time
	url := "https://github.com/" + entry.Repo.Repository.OwnerLogin + "/" + entry.Repo.Repository.Name
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
	if fork, err := ParseIntFromPattern(forkButton.Text(), `([\d\.]+)`); err == nil {
		info.Forks = fork
	}
	if scriptTag := doc.Find("react-partial[partial-name='repos-overview'] > script"); scriptTag.Length() > 0 {
		var detail Github_RepoDetail
		if err := json.Unmarshal([]byte(scriptTag.Text()), &detail); err != nil {
			return err
		}
		createdAt, err = time.Parse("2006-01-02T15:04:05.000Z", detail.Props.InitialPayload.Repo.CreatedAt)
		if err != nil {
			return err
		}
	}
	if articleTag := doc.Find("article").First(); articleTag.Length() > 0 {
		info.Readme = ExtractTextFromTag(articleTag)
	}
	updatedAt, err = time.Parse("2006-01-02T15:04:05.000Z", entry.Repo.Repository.UpdatedAt)
	if err != nil {
		return err
	}
	duration := time.Since(updatedAt)
	info.Url = url
	info.Name = entry.Repo.Repository.Name
	info.Summary = entry.Description
	info.Owner = entry.Repo.Repository.OwnerLogin
	info.CreatedAt = createdAt
	info.UpdatedAt = updatedAt
	info.Maintained = duration.Hours() < 24*30*60
	info.Stars = entry.Followers
	info.Language = entry.Language
	info.Libraries = strings.Join(entry.Topics, ",")
	return nil
}

func Github_GetSearchResult(query string, page int, info *Github_SearchResult) error {
	url := "https://github.com/search?q=" + strings.ReplaceAll(query, " ", "+") + "&type=repositories&p=" + strconv.Itoa(page)
	body, err := GetHtmlByProxy(url, false)
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

func Github_FindRecord(url string, record *models.GithubRepo) bool {
	dbInst := db.GetDB()
	result := dbInst.Where("repo_url = ?", url).First(record)
	return result.Error == nil
}

func Github_UpdateRecord(record *models.GithubRepo) error {
	dbInst := db.GetDB()
	result := dbInst.Save(&record)
	return result.Error
}

func Github_CreateRecord(record *models.GithubRepo) error {
	dbInst := db.GetDB()
	result := dbInst.Create(&record)
	return result.Error
}
