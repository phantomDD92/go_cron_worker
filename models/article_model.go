package models

import (
	"time"
)

// MarketingWebsiteArticle represents an article on a marketing website, managed by GORM
type MarketingWebsiteArticle struct {
	ID                     int64     `gorm:"primaryKey;autoIncrement"`
	Website                string    `gorm:"size:255"`
	WebsitePageTypes       string    `gorm:"size:255"`
	ArticleName            string    `gorm:"size:255"`
	ArticlePreviewImageURL string    `gorm:"size:255"`
	ArticleURL             string    `gorm:"size:255"`
	ArticleDate            time.Time `gorm:"type:timestamp"`
	ArticleDomain          string    `gorm:"size:255"`
	ArticleSummary         string
	ArticleRating          string `gorm:"size:255"`
	ArticleLength          int
	CodeLanguage           string `gorm:"size:255"`
	CodeLibraries          string `gorm:"size:255"`
	CodeLevel              string `gorm:"size:255"`
	CodeWorks              string `gorm:"size:255"`
	CodeGithubRepo         string `gorm:"size:255"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
