package utils

import (
	"errors"
	"fmt"
	"go_proxy_worker/db"
	"go_proxy_worker/models"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Article_GoogleEntry struct {
	Url     string `json:"url"`
	Date    string `json:"date"`
	Icon    string `json:"icon"`
	Title   string `json:"title"`
	Domain  string `json:"domain"`
	Content string `json:"content"`
	WebPage string `json:"web_page"`
	Repo    string `json:"repo"`
}

func Article_ParseEntry(entry *Article_GoogleEntry) error {
	body, err := GetHtmlByProxy(entry.Url, true)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return err
	}
	repo1 := ""
	repo2 := ""
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if strings.Contains(href, "https://github.com") {
			repo2 = href
			if repo, err := ParseGithubRepo(href); err == nil {
				repo1 = repo
			}
		}
	})
	if repo1 != "" {
		entry.Repo = repo1
	} else {
		entry.Repo = repo2
	}
	if titleTag := doc.Find("title").First(); titleTag.Length() > 0 {
		entry.Title = titleTag.Text()
	}
	bodyTag := doc.Find("body").First()
	if bodyTag.Length() == 0 {
		return errors.New("cannot find body tag")
	}
	bodyTag.Find("script").Remove()
	bodyTag.Find("noscript").Remove()
	bodyTag.Find("style").Remove()
	bodyTag.Find("iframe").Remove()
	entry.Content = bodyTag.Text()
	return nil
}

func Article_ParseGoogleLinks(url string, page int) (entries []Article_GoogleEntry, nextUrl string, err error) {
	body, err := GetHtmlByProxy(url, false)
	if err != nil {
		return entries, nextUrl, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return entries, nextUrl, err
	}
	doc.Find("#search #rso").Find("div[jscontroller='SC7lYd']").Each(func(i int, s *goquery.Selection) {
		entry := Article_GoogleEntry{}
		if linkTag := s.Find("a").First(); linkTag.Length() > 0 {
			href := linkTag.AttrOr("href", "")
			entry.Url = href
			domain := strings.Split(strings.ReplaceAll(href, "https://", ""), "/")[0]
			if segs := strings.Split(domain, "."); len(segs) >= 2 {
				entry.Domain = segs[len(segs)-2] + "." + segs[len(segs)-1]
			}
			if imgTag := linkTag.Find("img").First(); imgTag.Length() > 0 {
				entry.Icon = imgTag.AttrOr("src", "")
			}
		}
		if titleTag := s.Find("h3").First(); titleTag.Length() > 0 {
			entry.Title = strings.TrimSpace(titleTag.Text())
		}
		if dateTag := s.Find(".Sqrs4e > span").First(); dateTag.Length() > 0 {
			dateStr := strings.TrimSpace(dateTag.Text())
			entry.Date = dateStr
		}
		entries = append(entries, entry)
	})
	nextLinkTag := doc.Find("#botstuff").Find("[role='navigation']").Find("a#pnnext").First()
	if nextLinkTag.Length() > 0 {
		nextUrl = fmt.Sprintf("https://www.google.com%s", nextLinkTag.AttrOr("href", ""))
	}
	return entries, nextUrl, nil
}

func Article_FindRecord(url string, record *models.MarketingWebsiteArticle) bool {
	dbInst := db.GetDB()
	result := dbInst.Where("article_url = ?", url).First(record)
	return result.Error == nil
}

func Article_UpdateRecord(record *models.MarketingWebsiteArticle) error {
	dbInst := db.GetDB()
	result := dbInst.Save(&record)
	return result.Error
}

func Article_CreateRecord(record *models.MarketingWebsiteArticle) error {
	dbInst := db.GetDB()
	result := dbInst.Create(&record)
	return result.Error
}

func Article_FillRecordByEntry(record *models.MarketingWebsiteArticle, entry *Article_GoogleEntry, answer *ChatGPT_Answer) {
	record.Website = entry.WebPage
	record.ArticleName = entry.Title
	record.ArticleURL = entry.Url
	if date, err := ParseDate(entry.Date); err == nil {
		record.ArticleDate = date
	}
	record.ArticleDomain = entry.Domain
	record.CodeGithubRepo = entry.Repo
	record.Website = answer.WebSite
	record.WebsitePageTypes = answer.PageType
	record.ArticleSummary = answer.Summary
	record.CodeLanguage = answer.Language
	record.CodeLibraries = answer.Libraries
	record.CodeLevel = answer.CodeLevel
	record.CodeLanguage = answer.Language
}
